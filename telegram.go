package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

type Application struct {
	router Router
	bot    *tgbotapi.BotAPI
}

func NewApplication(api *tgbotapi.BotAPI) *Application {
	return &Application{
		bot: api,
	}
}

func (self *Application) Listen() {
	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.router.OnUpdate(ctx, self.bot, &update)
	}
}
