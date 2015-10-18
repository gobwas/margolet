package matcher

import "github.com/Syfaro/telegram-bot-api"

type AnyOf struct {
	Matchers []Matcher
}

func (self AnyOf) Match(message tgbotapi.Message) (match *Match, ok bool) {
	for _, matcher := range self.Matchers {
		if match, ok = matcher.Match(message); ok {
			return
		}
	}

	return
}
