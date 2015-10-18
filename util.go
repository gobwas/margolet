package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gobwas/telegram/matcher"
)

func mapRouteHandler(pattern string, handlers []Handler) []Handler {
	var mapped []Handler
	for _, handler := range handlers {
		mapped = append(mapped, &Conditional{matcher.Equal{pattern}, handler})
	}

	return mapped
}

func mapHandlerFunc(handlers []func(*Control, *tgbotapi.BotAPI, *tgbotapi.Update)) []Handler {
	var mapped []Handler
	for _, handler := range handlers {
		mapped = append(mapped, HandlerFunc(handler))
	}

	return mapped
}

func mapErrorHandlerFunc(handlers []func(*Control, *tgbotapi.BotAPI, *tgbotapi.Update, error)) []ErrorHandler {
	var mapped []ErrorHandler
	for _, handler := range handlers {
		mapped = append(mapped, ErrorHandlerFunc(handler))
	}

	return mapped
}
