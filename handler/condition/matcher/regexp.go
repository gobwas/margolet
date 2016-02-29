package matcher

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
)

type RegExp struct {
	Pattern *regexp.Regexp
}

func (self RegExp) Match(message tgbotapi.Message) (ok bool) {
	return self.Pattern.MatchString(message.Text)
}
