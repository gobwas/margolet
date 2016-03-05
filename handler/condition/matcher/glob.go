package matcher

import (
	"github.com/gobwas/glob"
	"gopkg.in/telegram-bot-api.v2"
)

type Glob struct {
	Source  Source
	Pattern glob.Glob
}

func (self Glob) Match(update tgbotapi.Update) bool {
	switch self.Source {
	case SourceText:
		self.Pattern.Match(update.Message.Text)
	case SourceQuery:
		return self.Pattern.Match(update.InlineQuery.Query)
	}

	return false
}
