package matcher

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Matcher interface {
	Match(message tgbotapi.Message) bool
}

type MatcherFunc func(message tgbotapi.Message) bool

func (self MatcherFunc) Match(message tgbotapi.Message) bool {
	return self(message)
}
