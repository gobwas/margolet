package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
)

type Conditional struct {
	matcher Matcher
	handler Handler
}

func (self *Conditional) Serve(ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if match, ok := self.matcher.Match(update.Message); ok {
		ctrl.WithValue(MATCH, *match)
		self.handler.Serve(ctrl, bot, update)
		return
	}

	ctrl.Next()
}
