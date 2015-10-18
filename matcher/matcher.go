package matcher

import (
	"github.com/Syfaro/telegram-bot-api"
)

type Slug struct {
	Key, Value string
}

type Match struct {
	Message tgbotapi.Message
	Slugs   []Slug
}

type Matcher interface {
	Match(message tgbotapi.Message) (*Match, bool)
}

type MatcherFunc func(message tgbotapi.Message) (*Match, bool)

func (self MatcherFunc) Match(message tgbotapi.Message) (*Match, bool) {
	return self(message)
}
