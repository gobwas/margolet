package telegram

import (
	"errors"
	"golang.org/x/net/context"
	"sync"
)

var Stopped = errors.New("control stopped")

type Signal int

const (
	NEXT Signal = iota
	ERROR
	STOP
)

type Wait func()

type Control struct {
	Done chan Signal

	wait     Wait
	lock     sync.Mutex
	err      error
	ctx      context.Context
	isCalled bool
}

func NewControl(ctx context.Context, wait Wait) *Control {
	return &Control{
		wait: wait,
		ctx:  ctx,
		Done: make(chan Signal),
	}
}

func (self *Control) setState(signal Signal, err error) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.isCalled {
		panic("could not be called twice")
	}
	self.isCalled = true

	self.err = err
	self.Done <- signal
	close(self.Done)

	return nil
}

func (self Control) Error() error {
	return self.err
}

func (self Control) Context() context.Context {
	return self.ctx
}

func (self *Control) SetContext(ctx context.Context) {
	self.ctx = ctx
}

func (self *Control) Next() error {
	if err := self.setState(NEXT, nil); err != nil {
		return err
	}

	self.wait()
	return nil
}

func (self *Control) Throw(e error) error {
	if err := self.setState(ERROR, e); err != nil {
		return err
	}

	self.err = e

	self.wait()
	return nil
}

func (self *Control) Stop() error {
	if err := self.setState(STOP, Stopped); err != nil {
		return err
	}

	self.wait()
	return nil
}
