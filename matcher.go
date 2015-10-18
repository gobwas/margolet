package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"regexp"
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

type Equal struct {
	pattern string
}

func (self Equal) Match(message tgbotapi.Message) (match *Match, ok bool) {
	if self.pattern == message.Text {
		match = &Match{Message: message}
		ok = true
	}

	return
}

type RegExpMatch struct {
	Match
	Submatch []string
}

type RegExp struct {
	pattern *regexp.Regexp
}

func (self RegExp) Match(message tgbotapi.Message) (match *Match, ok bool) {
	if self.pattern.MatchString(message.Text) {
		match = &Match{Message: message}
		for _, sub := range self.pattern.FindStringSubmatch(message.Text)[1:] {
			match.Slugs = append(match.Slugs, Slug{Value: sub})
		}
		ok = true
	}

	return
}
