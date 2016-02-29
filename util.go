package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/glob"
)

func mapRouteHandler(pattern string, handlers []Handler) []Handler {
	var mapped []Handler

	g := glob.MustCompile(pattern)
	for _, handler := range handlers {
		mapped = append(mapped, HandlerFunc(func(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
			if g.Match(update.Message.Text) {
				handler.Serve(ctrl, bot, update)
			} else {
				ctrl.Next()
			}
		}))
	}

	return mapped
}

func mapHandlerFunc(handlers []func(*Control, *tgbotapi.BotAPI, tgbotapi.Update)) []Handler {
	var mapped []Handler
	for _, handler := range handlers {
		mapped = append(mapped, HandlerFunc(handler))
	}

	return mapped
}

func mapErrorHandlerFunc(handlers []func(*Control, *tgbotapi.BotAPI, tgbotapi.Update, error)) []ErrorHandler {
	var mapped []ErrorHandler
	for _, handler := range handlers {
		mapped = append(mapped, ErrorHandlerFunc(handler))
	}

	return mapped
}
