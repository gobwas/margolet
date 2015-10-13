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

		g.It("Should call valid pattern", func() {
			var called int

			route := Route{Equal{"/pattern"}, HandlerFunc(func(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				called++
			})}

			route.Serve(context.Background(), &Control{}, &tgbotapi.BotAPI{}, &tgbotapi.Update{Message: tgbotapi.Message{Text: "/pattern"}})

			g.Assert(called).Eql(1)
		})

		g.It("Should set context 'route' key", func() {
			var val Match

			route := Route{Equal{"/abc"}, HandlerFunc(func(ctx context.Context, ctrl *Control, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				if v, ok := ctx.Value(ROUTE).(Match); ok {
					val = v
				}
			})}

			route.Serve(context.Background(), &Control{}, &tgbotapi.BotAPI{}, &tgbotapi.Update{Message: tgbotapi.Message{Text: "/abc"}})

			g.Assert(val).Eql(Match{Text: "/abc"})
		})

	})
}
