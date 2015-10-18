package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/franela/goblin"
	"github.com/gobwas/telegram/matcher"
	"testing"
)

func TestRoute(t *testing.T) {
	g := Goblin(t)

	g.Describe("Route", func() {

		g.It("Should call valid pattern", func() {
			var called int

			route := Conditional{matcher.Equal{"/pattern"}, HandlerFunc(func(ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				called++
			})}

			route.Serve(&Control{}, &tgbotapi.BotAPI{}, &tgbotapi.Update{Message: tgbotapi.Message{Text: "/pattern"}})

			g.Assert(called).Eql(1)
		})

		g.It("Should set context's MATCH key", func() {
			var val matcher.Match

			route := Conditional{matcher.Equal{"/abc"}, HandlerFunc(func(ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				if v, ok := ctrl.Context().Value(MATCH).(matcher.Match); ok {
					val = v
				}
			})}

			route.Serve(&Control{}, &tgbotapi.BotAPI{}, &tgbotapi.Update{Message: tgbotapi.Message{Text: "/abc"}})

			g.Assert(val.Message.Text).Eql("/abc")
		})

	})
}