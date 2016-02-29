package matcher

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/glob"
)

type Glob struct {
	Pattern glob.Glob
}

func (self Glob) Match(message tgbotapi.Message) bool {
	return self.Pattern.Match(message.Text)
}
