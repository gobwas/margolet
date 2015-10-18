package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/franela/goblin"
	"regexp"
	"testing"
)

func TestMatcher(t *testing.T) {
	g := Goblin(t)

	g.Describe("RegExp", func() {

		g.It("should match ok", func() {
			for _, test := range []struct {
				pattern *regexp.Regexp
				ok      bool
				text    tgbotapi.Message
				match   *Match
			}{
				{
					regexp.MustCompile(`/a/([a-z]+)`),
					true,
					tgbotapi.Message{Text: `/a/b`},
					&Match{
						Message: tgbotapi.Message{Text: `/a/b`},
						Slugs: []Slug{
							Slug{Value: `b`},
						},
					},
				},
			} {
				matcher := RegExp{test.pattern}
				match, ok := matcher.Match(test.text)

				g.Assert(ok).Equal(test.ok)

				if test.ok {
					g.Assert(match).Equal(test.match)
				}
			}
		})

	})
}
