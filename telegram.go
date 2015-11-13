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
	Token   string
	Debug   bool
	Polling *Polling
	WebHook *WebHook
}

type Polling struct {
	Offset  int
	Timeout int
}

type WebHook struct {
	URL    url.URL
	Listen Listen
	SSL    *SSL
}

type SSL struct {
	Cert string
	Key  string
}

type Listen struct {
	Addr string
	Port int
}

func New(config Config) (app *Application, err error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return
	}

	bot.Debug = config.Debug

	return &Application{
		bot:    bot,
		config: config,
	}, nil
}

func (self *Application) Listen() error {
	// todo move it somehow in app.Listen()
	switch {
	case self.config.WebHook != nil:
		config := self.config.WebHook

		// register hook
		webhookConfig := tgbotapi.WebhookConfig{
			URL:         &config.URL,
			Certificate: config.SSL.Cert,
		}
		if _, err := self.bot.SetWebhook(webhookConfig); err != nil {
			return err
		}

		// create server
		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Listen.Addr, config.Listen.Port), config.SSL.Cert, config.SSL.Key, nil)

		// listen path
		self.bot.ListenForWebhook("/" + config.URL.Path)

	case self.config.Polling != nil:
		config := self.config.Polling

		u := tgbotapi.NewUpdate(config.Offset)
		u.Timeout = config.Timeout

		err := self.bot.UpdatesChan(u)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Could not listen for updates: Polling or WebHook config should be set")
	}

	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.HandleUpdate(ctx, self.bot, update)
	}

	return nil
}
