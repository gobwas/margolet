package matcher

import (
	"gopkg.in/telegram-bot-api.v2"
	"regexp"
)

type RegExp struct {
	Pattern *regexp.Regexp
}

func (self RegExp) Match(message tgbotapi.Message) (ok bool) {
	return self.Pattern.MatchString(message.Text)
}
