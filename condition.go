package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram/matcher"
)

type Condition struct {
	Matcher matcher.Matcher
	Handler Handler
}

func (self Condition) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if self.Matcher.Match(update.Message) {
		self.Handler.Serve(ctrl, bot, update)
		return
	}

	ctrl.Next()
}
