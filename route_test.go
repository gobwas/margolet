package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/franela/goblin"
	"golang.org/x/net/context"
	"testing"
)

func TestRoute(t *testing.T) {
	g := Goblin(t)

	g.Describe("Route", func() {

		g.It("Should pass valid pattern", func() {
			NewRoute("/pattern", HandlerFunc(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {

			}))
		})

		g.It("Should call valid pattern", func() {
			var called int

			route := NewRoute("/pattern", HandlerFunc(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				called++
			}))

			route.Serve(context.Background(), &tgbotapi.BotAPI{}, &tgbotapi.Update{Message: tgbotapi.Message{Text: "/pattern"}}, &Control{})

			g.Assert(called).Eql(1)
		})

	})
}
