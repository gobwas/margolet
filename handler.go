package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"sync"
)

type Signal int

const (
	NEXT Signal = iota
	ERROR
	STOP
)

type Out struct {
	Signal Signal
	Error  error
}

type WaitGroup struct {
	resolved bool
	value    error
	lock     sync.Mutex
	count    int
}

func (self *WaitGroup) Init() {
	if self.resolved {
		panic("could not add to resolved handlersGroup")
	}
	self.lock.Lock()
}

func (self *WaitGroup) Add(delta int) {
	if self.resolved {
		panic("could not add to resolved handlersGroup")
	}
	self.count += delta
}

func (self *WaitGroup) Resolve(err error) {
	if self.resolved {
		panic("could not add to resolved handlersGroup")
	}
	self.resolved = true
	self.value = err
	self.lock.Unlock()
}

func (self *WaitGroup) Wait() error {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.value
}

type Control struct {
	Out      chan Out
	lock     sync.Mutex
	group    *WaitGroup
	isClosed bool
	isCalled bool
}

func NewControl(group *WaitGroup) *Control {
	return &Control{
		group: group,
		Out:   make(chan Out),
	}
}

func (self *Control) call(out Out) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.isCalled {
		panic("could not be called twice")
	}

	if self.isClosed {
		return fmt.Errorf("could not take control - it is already closed")
	}

	self.Out <- out
	close(self.Out)

	self.isCalled = true
	self.isClosed = true

	return nil
}

func (self *Control) close() {
	self.lock.Lock()
	defer self.lock.Unlock()

	close(self.Out)
	self.isClosed = true
}

func (self *Control) Next() error {
	if err := self.call(Out{NEXT, nil}); err != nil {
		return err
	}

	_ = self.group.Wait()
	return nil
}

func (self *Control) Error(e error) error {
	if err := self.call(Out{ERROR, e}); err != nil {
		return err
	}

	_ = self.group.Wait()
	return nil
}

func (self *Control) Stop() error {
	if err := self.call(Out{STOP, nil}); err != nil {
		return err
	}

	_ = self.group.Wait()
	return nil
}

func (self *Control) kill() {
	self.close()
}

// main handler
type Handler interface {
	Serve(*tgbotapi.BotAPI, *tgbotapi.Update, *Control)
}

type HandlerFunc func(*tgbotapi.BotAPI, *tgbotapi.Update, *Control)

// main error handler
type ErrorHandler interface {
	ServeError(*tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)
}

type ErrorHandlerFunc func(*tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)

func (self ErrorHandlerFunc) ServeError(bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
	self(bot, update, err, ctrl)
}

func (self HandlerFunc) Serve(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
	self(bot, update, ctrl)
}
