package slugger

import (
	"github.com/gobwas/telegram"
	"gopkg.in/telegram-bot-api.v2"
	"reflect"
	"strings"
)

type Call struct {
	Args []string
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

	var call Call
	args := strings.Split(path, argsSeparator)
	call.Args = make([]string, len(args))
	for i, arg := range args {
		call.Args[i] = arg
	}

	ctrl.NextWithValue(reflect.TypeOf(call), call)
	ctrl.Next()
}
