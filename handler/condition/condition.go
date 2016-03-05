package condition

import (
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/condition/matcher"
	"gopkg.in/telegram-bot-api.v2"
)

type Condition struct {
	Matcher matcher.Matcher
	Handler telegram.Handler
}

func (self Condition) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if self.Matcher.Match(update) {
		self.Handler.Serve(ctrl, bot, update)
		return
	}

	ctrl.Next()
}
