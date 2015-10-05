package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
	"regexp"
)

var regexp_command = regexp.MustCompile(`^\/[a-zA-Z0-9_]+$`)

//var slug = regexp.MustCompile(`\{\}`)

// @see https://core.telegram.org/bots#commands
func IsValidCommand(command string) bool {
	return regexp_command.Match([]byte(command))
}

func MatchPattern(pattern string, text string) bool {
	return pattern == text
}

type Route struct {
	pattern string
	handler Handler
}

func (self *Route) Serve(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
	if MatchPattern(self.pattern, update.Message.Text) {
		self.handler.Serve(ctx, bot, update, ctrl)
		return
	}

	ctrl.Next()
}

func NewRoute(pattern string, handler Handler) *Route {
	if !IsValidCommand(pattern) {
		panic(fmt.Sprintf("telegram: invalid command syntax: %q", pattern))
	}

	return &Route{pattern, handler}
}
