package telegram

import (
	"errors"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var ErrorStopped = errors.New("control stopped")
var ErrorSignalSent = errors.New("telegram: could not modify control after signal already sent")

type Signal int

const (
	s_NEXT Signal = iota
	s_ERROR
	s_STOP
)

type Wait func()

// Control is a structure, that brings control over handling phase.
type Control struct {
	mu sync.Mutex

	// ctx is context of current handling phase
	ctx context.Context

	// nextCtx specifies context for the next handling phase
	nextCtx context.Context

	// signal is a channel of Signal of control
	signal chan Signal

	// complete specifies the channel that will be closed
	// when all underneath handlers will send some signal or will be timed out
	complete <-chan struct{}

	// err specifies error that was happened during handling phase
	err error

	// signalSent specifies the flag that shows,
	// that control has already sent the signal to the upper layer
	signalSent bool
}

// NewControl initializes a new Control with given context.
func NewControl(complete <-chan struct{}, ctx context.Context) *Control {
	return &Control{
		ctx:      ctx,
		nextCtx:  ctx,
		complete: complete,
		signal:   make(chan Signal),
	}
}

func (self *Control) setNextContext(nextCtx context.Context) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.signalSent {
		return ErrorSignalSent
	}

	self.nextCtx = nextCtx
	return nil
}

func (self *Control) setState(signal Signal, err error) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.signalSent {
		return ErrorSignalSent
	}
	self.signalSent = true

	self.err = err
	self.signal <- signal

	return nil
}

func (self *Control) Error() error {
	return self.err
}

func (self *Control) Context() context.Context {
	return self.ctx
}

func (self *Control) NextContext() context.Context {
	return self.nextCtx
}

func (self *Control) NextWithCancel() (cancel context.CancelFunc, err error) {
	nextCtx, cancel := context.WithCancel(self.nextCtx)
	err = self.setNextContext(nextCtx)
	if err != nil {
		return
	}

	return
}

func (self *Control) NextWithTimeout(duration time.Duration) (cancel context.CancelFunc, err error) {
	nextCtx, cancel := context.WithTimeout(self.nextCtx, duration)
	err = self.setNextContext(nextCtx)
	if err != nil {
		return
	}

	return
}

func (self *Control) NextWithDeadline(deadline time.Time) (cancel context.CancelFunc, err error) {
	nextCtx, cancel := context.WithDeadline(self.nextCtx, deadline)
	err = self.setNextContext(nextCtx)
	if err != nil {
		return
	}

	return
}

func (self *Control) NextWithValue(key interface{}, val interface{}) error {
	return self.setNextContext(context.WithValue(self.nextCtx, key, val))
}

func (self *Control) Next() error {
	if err := self.setState(s_NEXT, nil); err != nil {
		return err
	}

	// lock until all underneath chain will complete
	// this allows to do some work after all stuff is done
	<-self.complete

	return nil
}

func (self *Control) Throw(e error) error {
	if err := self.setState(s_ERROR, e); err != nil {
		return err
	}

	self.err = e

	// lock until all underneath chain will complete
	// this allows to do some work after all stuff is done
	<-self.complete

	return nil
}

func (self *Control) Stop() error {
	if err := self.setState(s_STOP, ErrorStopped); err != nil {
		return err
	}

	// lock until all underneath chain will complete
	// this allows to do some work after all stuff is done
	<-self.complete

	return nil
}
