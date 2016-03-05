package matcher

import (
	"github.com/gobwas/glob"
	"gopkg.in/telegram-bot-api.v2"
)

type Glob struct {
	Pattern glob.Glob
}

func (self Glob) Match(message tgbotapi.Message) bool {
	return self.Pattern.Match(message.Text)
}
