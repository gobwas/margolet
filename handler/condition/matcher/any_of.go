package matcher

import "gopkg.in/telegram-bot-api.v2"

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
