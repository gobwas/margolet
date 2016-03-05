package matcher

import (
	"gopkg.in/telegram-bot-api.v2"
	"regexp"
)

type RegExp struct {
	Source  Source
	Pattern *regexp.Regexp
}

func (self RegExp) Match(update tgbotapi.Update) (ok bool) {
	switch self.Source {
	case SourceText:
		return self.Pattern.MatchString(update.Message.Text)
	case SourceQuery:
		return self.Pattern.MatchString(update.InlineQuery.Query)
	}

	return false
}
