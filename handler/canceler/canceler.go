package canceler

import (
	"github.com/gobwas/telegram"
	"gopkg.in/telegram-bot-api.v2"
	"time"
)

type Canceler struct {
	Timeout time.Duration
}

func (c *Canceler) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctrl.NextWithTimeout(c.Timeout)
	ctrl.Next()
}
