package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
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
	lock     sync.Mutex
	resolved bool
	count    int
	error    error
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
	self.error = err
	self.lock.Unlock()
}

func (self *WaitGroup) Wait() error {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.error
}

type Control struct {
	Out      chan Out
	lock     sync.Mutex
	group    *WaitGroup
	isKilled error
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

	if self.isKilled != nil {
		return fmt.Errorf("could not take control it is already closed: %s", self.isKilled)
	}

	self.Out <- out
	close(self.Out)
	self.isCalled = true

	return nil
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

func (self *Control) kill(err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	close(self.Out)
	self.isKilled = err
}

// main handler
type Handler interface {
	Serve(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, *Control)
}

type HandlerFunc func(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, *Control)

// main error handler
type ErrorHandler interface {
	ServeError(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)
}

type ErrorHandlerFunc func(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)

func (self ErrorHandlerFunc) ServeError(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
	self(ctx, bot, update, err, ctrl)
}

func (self HandlerFunc) Serve(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
	self(ctx, bot, update, ctrl)
}
