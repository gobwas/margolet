package matcher

import "github.com/go-telegram-bot-api/telegram-bot-api"

type Equal struct {
	Pattern string
}

func (self Equal) Match(message tgbotapi.Message) (match *Match, ok bool) {
	if self.Pattern == message.Text {
		match = &Match{Message: message}
		ok = true
	}

	return
}
