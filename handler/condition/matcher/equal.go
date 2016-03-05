package matcher

import "gopkg.in/telegram-bot-api.v2"

type Equal struct {
	Pattern string
}

func (self Equal) Match(message tgbotapi.Message) bool {
	return self.Pattern == message.Text
}
