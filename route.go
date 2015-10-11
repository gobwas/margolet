package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

type Route struct {
	matcher Matcher
	handler Handler
}

func (self *Route) Serve(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if match, ok := self.matcher.Match(update.Message.Text); ok {
		self.handler.Serve(context.WithValue(ctx, "route", *match), ctrl, bot, update)
		return
	}

	ctrl.Next()
}
