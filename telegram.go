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
	SSL    SSL
	Listen Listen
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

	switch {
	case self.config.WebHook != nil:
		config := self.config.WebHook

		if _, err := self.bot.SetWebhook(tgbotapi.WebhookConfig{URL: &config.URL, Certificate: config.SSL.Cert}); err != nil {
			return err
		}

		self.bot.ListenForWebhook("/" + config.URL.Path)
		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Listen.Addr, config.Listen.Port), config.SSL.Cert, config.SSL.Key, nil)
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

	return NewByBot(bot)
}

func (self *Application) Listen() error {
	for update := range self.bot.Updates {
		ctx := context.Background()
		go self.HandleUpdate(ctx, self.bot, &update)
	}
}
