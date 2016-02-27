package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// main handler
type Handler interface {
	Serve(*Control, *tgbotapi.BotAPI, tgbotapi.Update)
}

type HandlerFunc func(*Control, *tgbotapi.BotAPI, tgbotapi.Update)

func (self HandlerFunc) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	self(ctrl, bot, update)
}

// main error handler
type ErrorHandler interface {
	ServeError(*Control, *tgbotapi.BotAPI, tgbotapi.Update, error)
}

type ErrorHandlerFunc func(*Control, *tgbotapi.BotAPI, tgbotapi.Update, error)

func (self ErrorHandlerFunc) ServeError(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {
	self(ctrl, bot, update, err)
}
