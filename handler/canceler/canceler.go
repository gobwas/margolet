package canceler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"time"
)

type Canceler struct {
	Timeout time.Duration
}

func (c *Canceler) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctrl.NextWithTimeout(c.Timeout)
	ctrl.Next()
}
