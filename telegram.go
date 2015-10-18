package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

type Application struct {
	Router
	bot    *tgbotapi.BotAPI
	config Config
}

type Config struct {
	Token      string
	Debug      bool
	UseWebHook bool
	Polling    Polling
	WebHook    WebHook
}

type Polling struct {
	Offset  int
	Timeout int
}

type WebHook struct {
	URL    url.URL
	SSL    SSL
	Server Server
}

type SSL struct {
	Cert string
	Key  string
}

type Server struct {
	Host string
	Port int
}

func New(config Config) (app *Application, err error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return
	}
	bot.Debug = config.Debug

	app = &Application{
		config: config,
		bot:    bot,
	}

	return
}

func (self *Application) Listen() error {
	if self.config.UseWebHook {
		config := self.config.WebHook

		if _, err := self.bot.SetWebhook(tgbotapi.WebhookConfig{URL: &config.URL, Certificate: config.SSL.Cert}); err != nil {
			return err
		}

		self.bot.ListenForWebhook("/" + config.URL.Path)
		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), config.SSL.Cert, config.SSL.Key, nil)
	} else {
		config := self.config.Polling
		u := tgbotapi.NewUpdate(config.Offset)
		u.Timeout = config.Timeout

		err := self.bot.UpdatesChan(u)
		if err != nil {
			return err
		}
	}

	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.HandleUpdate(ctx, self.bot, &update)
	}
}
