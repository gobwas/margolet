package telegram

import (
	. "github.com/franela/goblin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram/matcher"
	"testing"
)

func TestRoute(t *testing.T) {
	g := Goblin(t)

	g.Describe("Route", func() {

		g.It("Should call valid pattern", func() {
			var called int

			route := Condition{matcher.Equal{"/pattern"}, HandlerFunc(func(ctrl *Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
				called++
			})}

			route.Serve(&Control{}, &tgbotapi.BotAPI{}, tgbotapi.Update{Message: tgbotapi.Message{Text: "/pattern"}})

			g.Assert(called).Eql(1)
		})

	})
}
