package matcher

import (
	"github.com/Syfaro/telegram-bot-api"
	"regexp"
)

type RegExpMatch struct {
	Match
	Submatch []string
}

type RegExp struct {
	Pattern *regexp.Regexp
}

func (self RegExp) Match(message tgbotapi.Message) (match *Match, ok bool) {
	if self.Pattern.MatchString(message.Text) {
		match = &Match{Message: message}
		for _, sub := range self.Pattern.FindStringSubmatch(message.Text)[1:] {
			match.Slugs = append(match.Slugs, Slug{Value: sub})
		}
		ok = true
	}

	return
}
