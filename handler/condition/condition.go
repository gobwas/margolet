package condition

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/condition/matcher"
)

type Condition struct {
	Matcher matcher.Matcher
	Handler telegram.Handler
}

func (self Condition) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if self.Matcher.Match(update.Message) {
		self.Handler.Serve(ctrl, bot, update)
		return
	}

	ctrl.Next()
}
