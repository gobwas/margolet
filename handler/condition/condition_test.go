package condition

import (
	. "github.com/franela/goblin"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/condition/matcher"
	"gopkg.in/telegram-bot-api.v2"
	"testing"
)

func TestRoute(t *testing.T) {
	g := Goblin(t)

	g.Describe("Route", func() {

		g.It("Should call valid pattern", func() {
			var called int

			route := Condition{matcher.Equal{"/pattern"}, telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
				called++
			})}

			route.Serve(&telegram.Control{}, &tgbotapi.BotAPI{}, tgbotapi.Update{Message: tgbotapi.Message{Text: "/pattern"}})

			g.Assert(called).Eql(1)
		})

	})
}
