package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/franela/goblin"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	g := Goblin(t)

	g.Describe("Router", func() {
		var bot tgbotapi.BotAPI
		var update tgbotapi.Update

		g.BeforeEach(func() {
			bot = tgbotapi.BotAPI{}
			update = tgbotapi.Update{}
		})

		g.It("Should UseFunc() okay", func() {
			router := NewRouter()
			router.UseFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				fmt.Println("Im okay")
				ctrl.Next()
				//				if err := ctrl.Next(); err != nil {
				//					fmt.Println(err)
				//				} else {
				//					fmt.Println("Next is ok")
				//				}
			})

			router.OnUpdate(&bot, &update)
			time.Sleep(1 * time.Millisecond)
		})

		g.It("Should call in chain", func() {
			router := NewRouter()
			router.UseFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				fmt.Println("A called")
				ctrl.Error(fmt.Errorf("A produced error"))
			})
			router.UseFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				fmt.Println("B called")
				ctrl.Next()
			})
			router.UseErrFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				fmt.Println("A error fixer called", err)
				ctrl.Error(fmt.Errorf("Could not fix A's error"))
			})
			router.UseErrFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				fmt.Println("B error fixer called", err)
				ctrl.Next()
			})
			router.UseErrFunc(func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				fmt.Println("C error fixer called", err)
				ctrl.Error(err)
			})

			err := router.OnUpdate(&bot, &update)
			if err != nil {
				fmt.Println("OnUpdate() err", err)
			} else {
				fmt.Println("OnUpdate() succ")
			}
			time.Sleep(1 * time.Millisecond)
		})

	})
}
