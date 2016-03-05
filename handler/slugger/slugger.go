package slugger

import (
	"github.com/gobwas/telegram"
	"gopkg.in/telegram-bot-api.v2"
	"reflect"
	"strings"
)

type Call struct {
	Args  []string
	Query []string
}

type Slugger struct {
}

const methodPrefix = "/"
const argsSeparator = " "

func (r *Slugger) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var call Call

	text := strings.TrimPrefix(update.Message.Text, methodPrefix)
	fill(&call.Args, text)
	fill(&call.Query, update.InlineQuery.Query)

	ctrl.NextWithValue(reflect.TypeOf(call), call)
	ctrl.Next()
}

func fill(dest *[]string, s string) {
	if s == "" {
		return
	}

	args := strings.Split(s, argsSeparator)
	d := make([]string, len(args))
	for i, arg := range args {
		d[i] = arg
	}
	*dest = d
}
