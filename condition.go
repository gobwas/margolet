package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gobwas/telegram/matcher"
)

type Condition struct {
	Matcher matcher.Matcher
	Handler Handler
}

func (self *Condition) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if match, ok := self.Matcher.Match(update.Message); ok {
		ctrl.WithValue(v_MATCH, *match)
		self.Handler.Serve(ctrl, bot, update)
		return
	}

	ctrl.Next()
}
