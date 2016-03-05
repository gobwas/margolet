package matcher

import "gopkg.in/telegram-bot-api.v2"

type Equal struct {
	Source  Source
	Pattern string
}

func (self Equal) Match(update tgbotapi.Update) bool {
	switch self.Source {
	case SourceText:
		return self.Pattern == update.Message.Text
	case SourceQuery:
		return self.Pattern == update.InlineQuery.Query
	}

	return false
}
