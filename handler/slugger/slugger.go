package slugger

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"reflect"
	"strings"
)

type Call struct {
	Method string
	Args   []string
}

type Slugger struct {
}

const methodPrefix = "/"
const argsSeparator = " "

func (r *Slugger) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := update.Message.Text
	path := strings.TrimPrefix(text, methodPrefix)

	// this is not a method
	if path == text {
		ctrl.Next()
		return
	}

	args := strings.Split(path, argsSeparator)
	var call Call
	if len(args) == 1 {
		call.Method = args[0]
	} else {
		call.Args = make([]string, len(args)-1)
		for i, arg := range args[1:] {
			call.Args[i] = arg
		}
	}

	ctrl.NextWithValue(reflect.TypeOf(call), call)
	ctrl.Next()
}
