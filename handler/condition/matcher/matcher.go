package matcher

import (
	"gopkg.in/telegram-bot-api.v2"
)

type Matcher interface {
	Match(message tgbotapi.Message) bool
}

type MatcherFunc func(message tgbotapi.Message) bool

func (self MatcherFunc) Match(message tgbotapi.Message) bool {
	return self(message)
}
