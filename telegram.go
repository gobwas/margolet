package telegram

import "github.com/Syfaro/telegram-bot-api"

type Application struct {
	Router
	bot *tgbotapi.BotAPI
}

func NewApplication(api *tgbotapi.BotAPI) *Application {
	return &Application{
		bot: api,
	}
}

func (self *Application) Listen() {
	for update := range self.bot.Updates {
		go self.OnUpdate(self.bot, &update)
	}
}
