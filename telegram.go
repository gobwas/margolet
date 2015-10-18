package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

type Application struct {
	Router
	bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) (app *Application, err error) {
	app = &Application{
		bot: bot,
	}

	return
}

func (self *Application) Listen() {
	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.HandleUpdate(ctx, self.bot, &update)
	}
}
