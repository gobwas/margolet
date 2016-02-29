package matcher

import "github.com/go-telegram-bot-api/telegram-bot-api"

type AnyOf struct {
	Matchers []Matcher
}

func (self AnyOf) Match(message tgbotapi.Message) (ok bool) {
	for _, matcher := range self.Matchers {
		if matcher.Match(message) {
			return
		}
	}

	return
}
