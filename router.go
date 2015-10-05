package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

type Router struct {
	handlers      []Handler
	errorHandlers []ErrorHandler
}

func NewRouter() *Router {
	return &Router{}
}

func (self *Router) Serve(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
	if err := self.OnUpdate(ctx, bot, update); err != nil {
		ctrl.Error(err)
	}

	ctrl.Next()
}

func (self *Router) Use(handlers ...Handler) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) UseFunc(handlers ...func(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, *Control)) {
	self.Use(mapHandlerFunc(handlers)...)
}

func (self *Router) UseOn(pattern string, handlers ...Handler) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, handlers)...)
}

func (self *Router) UseFuncOn(pattern string, handlers ...func(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, *Control)) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, mapHandlerFunc(handlers))...)
}

func (self *Router) UseErr(handlers ...ErrorHandler) {
	self.errorHandlers = append(self.errorHandlers, handlers...)
}

func (self *Router) UseErrFunc(handlers ...func(context.Context, *tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)) {
	self.UseErr(mapErrorHandlerFunc(handlers)...)
}

func (self *Router) OnUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) (err error) {
	err = iterate(ctx, self.handlers, bot, update)
	if err != nil {
		err = iterateError(ctx, self.errorHandlers, bot, update, err)
	}

	return
}

func iterate(ctx context.Context, handlers []Handler, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	group := WaitGroup{}
	group.Init()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, handler := range handlers {
		ctrl := NewControl(&group)
		group.Add(1)

		// start handling
		go handler.Serve(ctx, bot, update, ctrl)

		// race with ctx
		select {
		case <-ctx.Done():
			err := ctx.Err()
			// close channel
			ctrl.kill(err)
			group.Resolve(err)
			return err
		case out := <-ctrl.Out:
			switch out.Signal {
			case NEXT:
				continue
			case ERROR:
				group.Resolve(out.Error)
				return out.Error
			case STOP:
				group.Resolve(nil)
				return nil
			}
		}
	}

	// release whole loop
	group.Resolve(nil)
	return nil
}

func iterateError(ctx context.Context, handlers []ErrorHandler, bot *tgbotapi.BotAPI, update *tgbotapi.Update, initial error) error {
	group := WaitGroup{}
	group.Init()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := initial
	for _, handler := range handlers {
		ctrl := NewControl(&group)
		group.Add(1)

		go handler.ServeError(ctx, bot, update, err, ctrl)

		select {
		case <-ctx.Done():
			err := ctx.Err()
			ctrl.kill(err)
			group.Resolve(err)
			return err
		case out := <-ctrl.Out:
			switch out.Signal {
			case ERROR:
				err = out.Error
			case STOP, NEXT:
				group.Resolve(nil)
				return nil
			}
		}
	}

	// release whole loop
	group.Resolve(err)
	return err
}
