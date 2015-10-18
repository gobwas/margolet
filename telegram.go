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
	bot *tgbotapi.BotAPI
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

func NewByBot(bot *tgbotapi.BotAPI) (*Application, error) {
	return &Application{
		bot: bot,
	}, nil
}

func New(config Config) (app *Application, err error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return
	}
	bot.Debug = config.Debug

	// todo move it somehow in app.Listen()
	switch {
	case config.WebHook != nil:
		config := config.WebHook

		if _, err := bot.SetWebhook(tgbotapi.WebhookConfig{URL: &config.URL, Certificate: config.SSL.Cert}); err != nil {
			return nil, err
		}

		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Listen.Addr, config.Listen.Port), config.SSL.Cert, config.SSL.Key, nil)
		bot.ListenForWebhook("/" + config.URL.Path)

	case config.Polling != nil:
		config := config.Polling

		u := tgbotapi.NewUpdate(config.Offset)
		u.Timeout = config.Timeout

		err := bot.UpdatesChan(u)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Could not listen for updates: Polling or WebHook config should be set")
	}

	return NewByBot(bot)
}

func (self *Application) Listen() error {
	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.HandleUpdate(ctx, self.bot, update)
	}

	return nil
}
