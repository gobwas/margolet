package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
	"sync"
)

type Router struct {
	handlers      []Handler
	errorHandlers []ErrorHandler
}

func NewRouter() *Router {
	return &Router{}
}

func (self *Router) Serve(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if err := self.HandleUpdate(ctx, bot, update); err != nil {
		ctrl.Throw(err)
	}

	ctrl.Next()
}

func (self *Router) Use(handlers ...Handler) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) UseFunc(handlers ...func(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update)) {
	self.Use(mapHandlerFunc(handlers)...)
}

func (self *Router) UseOn(pattern string, handlers ...Handler) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, handlers)...)
}

func (self *Router) UseFuncOn(pattern string, handlers ...func(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update)) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, mapHandlerFunc(handlers))...)
}

func (self *Router) UseErr(handlers ...ErrorHandler) {
	self.errorHandlers = append(self.errorHandlers, handlers...)
}

func (self *Router) UseErrFunc(handlers ...func(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update, error)) {
	self.UseErr(mapErrorHandlerFunc(handlers)...)
}

func (self *Router) HandleUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) (err error) {
	err = self.traverseUpdate(ctx, bot, update)
	if err != nil {
		err = self.traverseError(ctx, bot, update, err)
	}

	return
}

func (self Router) traverseUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var group sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prevent release on first Done() in loop
	group.Add(1)
	defer group.Done()

	for _, handler := range self.handlers {
		ctrl := NewControl(ctx, Wait(group.Wait))
		group.Add(1)

		// start handling
		go handler.Serve(ctx, ctrl, bot, update)

		// race with ctx
		select {
		case <-ctx.Done():
			group.Done()

			err := ctx.Err()
			ctrl.Throw(err)

			return err

		case signal := <-ctrl.Done:
			group.Done()

			switch signal {
			case NEXT:
				ctx = ctrl.Context()
				continue

			case ERROR:
				err := ctrl.Error()
				cancel()
				return err

			case STOP:
				cancel()
				return nil
			}
		}
	}

	return nil
}

func (self Router) traverseError(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error) error {
	var group sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prevent release on first Done() in loop
	group.Add(1)
	defer group.Done()

	for _, handler := range self.errorHandlers {
		ctrl := NewControl(ctx, Wait(group.Wait))
		group.Add(1)

		go handler.ServeError(ctx, ctrl, bot, update, err)

		select {
		case <-ctx.Done():
			group.Done()

			err := ctx.Err()
			ctrl.Throw(err)

			return err

		case signal := <-ctrl.Done:
			group.Done()

			switch signal {
			case ERROR:
				err = ctrl.Error()
				continue
			case NEXT:
				return nil
			case STOP:
				cancel()
				return nil
			}
		}
	}

	return err
}
