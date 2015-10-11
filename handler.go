package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

// main handler
type Handler interface {
	Serve(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update)
}

type HandlerFunc func(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update)

// main error handler
type ErrorHandler interface {
	ServeError(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update, error)
}

type ErrorHandlerFunc func(context.Context, *Control, *tgbotapi.BotAPI, *tgbotapi.Update, error)

func (self ErrorHandlerFunc) ServeError(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error) {
	self(ctx, ctrl, bot, update, err)
}

func (self HandlerFunc) Serve(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	self(ctx, ctrl, bot, update)
}
