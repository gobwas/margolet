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

func (self *Router) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if err := self.HandleUpdate(ctrl.Context(), bot, update); err != nil {
		ctrl.Throw(err)
	}

	ctrl.Next()
}

func (self *Router) Use(handlers ...Handler) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) UseFunc(handlers ...func(*Control, *tgbotapi.BotAPI, *tgbotapi.Update)) {
	self.Use(mapHandlerFunc(handlers)...)
}

func (self *Router) UseOn(pattern string, handlers ...Handler) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, handlers)...)
}

func (self *Router) UseFuncOn(pattern string, handlers ...func(*Control, *tgbotapi.BotAPI, *tgbotapi.Update)) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, mapHandlerFunc(handlers))...)
}

func (self *Router) UseErr(handlers ...ErrorHandler) {
	self.errorHandlers = append(self.errorHandlers, handlers...)
}

func (self *Router) UseErrFunc(handlers ...func(*Control, *tgbotapi.BotAPI, *tgbotapi.Update, error)) {
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
		go handler.Serve(ctrl, bot, update)

		// race with ctx
		select {
		case <-ctx.Done():
			group.Done()

			err := ctx.Err()
			ctrl.Throw(err)

			return err

		case signal := <-ctrl.done:
			group.Done()

			switch signal {
			case s_NEXT:
				ctx = ctrl.Context()
				continue

			case s_ERROR:
				err := ctrl.error()
				cancel()
				return err

			case s_STOP:
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

		go handler.ServeError(ctrl, bot, update, err)

		select {
		case <-ctx.Done():
			group.Done()

			err := ctx.Err()
			ctrl.Throw(err)

			return err

		case signal := <-ctrl.done:
			group.Done()

			switch signal {
			case s_ERROR:
				err = ctrl.error()
				continue
			case s_NEXT:
				continue
			case s_STOP:
				cancel()
				return nil
			}
		}
	}

	return err
}
