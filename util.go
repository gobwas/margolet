package telegram

import "github.com/Syfaro/telegram-bot-api"

func mapRouteHandler(pattern string, handlers []Handler) []Handler {
	var mapped []Handler
	for _, handler := range handlers {
		mapped = append(mapped, NewRoute(pattern, handler))
	}

	return mapped
}

func mapHandlerFunc(handlers []func(*tgbotapi.BotAPI, *tgbotapi.Update, *Control)) []Handler {
	var mapped []Handler
	for _, handler := range handlers {
		mapped = append(mapped, HandlerFunc(handler))
	}

	return mapped
}

func mapErrorHandlerFunc(handlers []func(*tgbotapi.BotAPI, *tgbotapi.Update, error, *Control)) []ErrorHandler {
	var mapped []ErrorHandler
	for _, handler := range handlers {
		mapped = append(mapped, ErrorHandlerFunc(handler))
	}

	return mapped
}
