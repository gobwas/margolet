package matcher

import (
	"gopkg.in/telegram-bot-api.v2"
)

type Source int

const (
	SourceText = iota
	SourceQuery
)

type Matcher interface {
	Match(message tgbotapi.Update) bool
}

type MatcherFunc func(message tgbotapi.Update) bool

func (self MatcherFunc) Match(message tgbotapi.Update) bool {
	return self(message)
}
