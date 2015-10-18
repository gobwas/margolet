package telegram

import (
	"errors"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var Stopped = errors.New("control stopped")

type Signal int

const (
	s_NEXT Signal = iota
	s_ERROR
	s_STOP
)

type Wait func()

type Control struct {
	ctx      context.Context
	done     chan Signal
	wait     Wait
	mu       sync.Mutex
	err      error
	isCalled bool
}

func NewControl(ctx context.Context, wait Wait) *Control {
	return &Control{
		ctx:  ctx,
		wait: wait,
		done: make(chan Signal),
	}
}

func (self *Control) setContext(ctx context.Context) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.isCalled {
		panic("could not be called after control was taken")
	}

	self.ctx = ctx
}

func (self *Control) setState(signal Signal, err error) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.isCalled {
		panic("could not be called twice")
	}
	self.isCalled = true

	self.err = err
	self.done <- signal
	close(self.done)

	return nil
}

// returns error that wath throwed with `ctrl.Throw()`
func (self Control) error() error {
	return self.err
}

// returns context that was passed to the `ctrl.NextWithContext()` method
func (self Control) Context() context.Context {
	return self.ctx
}

func (self *Control) WithCancel() context.CancelFunc {
	ctx, cancel := context.WithCancel(self.ctx)
	self.setContext(ctx)
	return cancel
}

func (self *Control) WithTimeout(duration time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(self.ctx, duration)
	self.setContext(ctx)
	return cancel
}

func (self *Control) WithDeadline(deadline time.Time) context.CancelFunc {
	ctx, cancel := context.WithDeadline(self.ctx, deadline)
	self.setContext(ctx)
	return cancel
}

func (self *Control) WithValue(key interface{}, val interface{}) {
	self.setContext(context.WithValue(self.ctx, key, val))
}

func (self *Control) Next() error {
	if err := self.setState(s_NEXT, nil); err != nil {
		return err
	}

	self.wait()
	return nil
}

func (self *Control) Throw(e error) error {
	if err := self.setState(s_ERROR, e); err != nil {
		return err
	}

	self.err = e

	self.wait()
	return nil
}

func (self *Control) Stop() error {
	if err := self.setState(s_STOP, Stopped); err != nil {
		return err
	}

	self.wait()
	return nil
}
