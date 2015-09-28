package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
)

type Router struct {
	handlers      []Handler
	errorHandlers []ErrorHandler
}

func NewRouter() *Router {
	return &Router{}
}

func (self *Router) Serve(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
	if err := self.OnUpdate(bot, update); err != nil {
		ctrl.Error(err)
	}

	ctrl.Next()
}

func (self *Router) Use(handlers ...Handler) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) UseFunc(handlers ...func(*tgbotapi.BotAPI, *tgbotapi.Update, *Control)) {
	self.Use(mapHandlerFunc(handlers)...)
}

func (self *Router) UseOn(pattern string, handlers ...Handler) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, handlers)...)
}

func (self *Router) UseFuncOn(pattern string, handlers ...func(*tgbotapi.BotAPI, *tgbotapi.Update, *Control)) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, mapHandlerFunc(handlers))...)
}

func (self *Router) UseErr(handlers ...ErrorHandler) {
	self.errorHandlers = append(self.errorHandlers, handlers...)
}

func (self *Router) UseErrFunc(handlers ...func(*tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)) {
	self.UseErr(mapErrorHandlerFunc(handlers)...)
}

func (self *Router) OnUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (err error) {
	err = iterate(self.handlers, bot, update)
	if err != nil {
		err = iterateError(self.errorHandlers, bot, update, err)
	}

	return
}

func iterate(handlers []Handler, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	group := WaitGroup{}
	group.Init()

	for _, handler := range handlers {
		ctrl := NewControl(&group)
		group.Add(1)

		// start timeout
		//		timeout := make(chan int)
		//		go func() {
		//			time.Sleep(timeout * time.Millisecond)
		//			timeout <- 1
		//			close(timeout)
		//		}()

		// start handling
		go handler.Serve(bot, update, ctrl)

		// race
		//		select {
		//		case out := <-ctrl.Out:
		out := <-ctrl.Out
		switch out.Signal {
		case NEXT:
			//
		case ERROR:
			group.Resolve(out.Error)
			return out.Error
		case STOP:
			group.Resolve(nil)
			return nil
		}
		//		case <-timeout:
		//			ctrl.kill()
		//		}
	}

	// release whole loop
	group.Resolve(nil)
	return nil
}

func iterateError(handlers []ErrorHandler, bot *tgbotapi.BotAPI, update *tgbotapi.Update, initial error) error {
	group := WaitGroup{}
	group.Init()

	err := initial
	for _, handler := range handlers {
		ctrl := NewControl(&group)
		group.Add(1)

		go handler.ServeError(bot, update, err, ctrl)

		out := <-ctrl.Out
		switch out.Signal {
		case ERROR:
			err = out.Error
		case STOP, NEXT:
			group.Resolve(nil)
			return nil
		}
	}

	// release whole loop
	group.Resolve(err)
	return err
}
