package matcher

import "gopkg.in/telegram-bot-api.v2"

type AnyOf struct {
	Matchers []Matcher
}

func (self AnyOf) Match(update tgbotapi.Update) (ok bool) {
	for _, matcher := range self.Matchers {
		if matcher.Match(update) {
			return
		}
	}

	return
}
