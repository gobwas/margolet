package matcher

import "github.com/go-telegram-bot-api/telegram-bot-api"

type Equal struct {
	Pattern string
}

func (self Equal) Match(message tgbotapi.Message) bool {
	return self.Pattern == message.Text
}
