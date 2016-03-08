package telegram

import (
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v2"
)

type Router struct {
	handlers      []Handler
	errorHandlers []ErrorHandler
}

func NewRouter() *Router {
	return &Router{}
}

func (self *Router) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if err := self.HandleUpdate(ctrl.Context(), bot, update); err != nil {
		ctrl.Throw(err)
	}

	ctrl.Next()
}

func (self *Router) Use(handlers ...Handler) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) UseFunc(handlers ...func(*Control, *tgbotapi.BotAPI, tgbotapi.Update)) {
	self.Use(mapHandlerFunc(handlers)...)
}

func (self *Router) UseOn(pattern string, handlers ...Handler) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, handlers)...)
}

func (self *Router) UseFuncOn(pattern string, handlers ...func(*Control, *tgbotapi.BotAPI, tgbotapi.Update)) {
	self.handlers = append(self.handlers, mapRouteHandler(pattern, mapHandlerFunc(handlers))...)
}

func (self *Router) UseErr(handlers ...ErrorHandler) {
	self.errorHandlers = append(self.errorHandlers, handlers...)
}

func (self *Router) UseErrFunc(handlers ...func(*Control, *tgbotapi.BotAPI, tgbotapi.Update, error)) {
	self.UseErr(mapErrorHandlerFunc(handlers)...)
}

func (self *Router) HandleUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (err error) {
	err = self.traverseUpdate(ctx, bot, update)
	if err != nil {
		err = self.traverseError(ctx, bot, update, err)
	}

	return
}

func (self *Router) traverseUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (err error) {
	// the complete will notify all controls in all handlers
	// that traversal was completed
	complete := make(chan struct{})

	// initial next context for the first handler in the chain
	nextContext, cancel := context.WithCancel(ctx)

	logger := &Logger{fmt.Sprintf("[%d][%s] ", update.UpdateID, update.Message.From.String())}

handling:
	for _, handler := range self.handlers {
		ctrl := NewControl(complete, nextContext, logger)
		go handler.Serve(ctrl, bot, update)

		select {
		case <-nextContext.Done():
			err = nextContext.Err()
			break handling

		case signal := <-ctrl.signal:
			switch signal {
			case s_NEXT:
				nextContext = ctrl.NextContext()
				continue handling

			case s_ERROR:
				err = ctrl.Error()
				cancel()
				break handling

			case s_STOP:
				err = ctrl.Error()
				cancel()
				break handling

			default:
				err = fmt.Errorf("telegram: unknown handler signal code: %d", signal)
				cancel()
				break handling
			}
		}
	}

	close(complete)

	return
}

func (self *Router) traverseError(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) (fatal error) {
	// the complete will notify all controls in all handlers
	// that traversal was completed
	complete := make(chan struct{})

	// initial next context for the first handler in the chain
	nextContext, cancel := context.WithCancel(ctx)

handling:
	for _, handler := range self.errorHandlers {
		ctrl := NewControl(complete, nextContext)
		go handler.ServeError(ctrl, bot, update, err)

		select {
		case <-nextContext.Done():
			fatal = nextContext.Err()
			break handling

		case signal := <-ctrl.signal:
			switch signal {
			case s_NEXT:
				// next means that error handler could not fix anything
				nextContext = ctrl.NextContext()
				continue handling

			case s_ERROR:
				// error in error handling means
				// that there is fatal error =)
				fatal = ctrl.Error()
				cancel()
				break handling

			case s_STOP:
				// do not read the stop error
				// cause in error handling this means
				// that error was handled OK
				cancel()
				break handling

			default:
				fatal = fmt.Errorf("telegram: unknown error handler signal code: %d", signal)
				cancel()
				break handling
			}
		}
	}

	close(complete)

	return
}
