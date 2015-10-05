package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/franela/goblin"
	"golang.org/x/net/context"
	"testing"
	"time"
)

type Call struct {
	Time time.Time
	Args []interface{}
}

type Spy struct {
	CallCount int
	Calls     []Call
}

type WithCalls interface {
	GetCalls() []Call
}

func (s Spy) GetCalls() []Call {
	return s.Calls
}

type HandlerSpy struct {
	Spy
	Fn func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control)
}

type ErrorSpy struct {
	Spy
	Fn func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control)
}

func NewHandlerSpy(fn func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control)) *HandlerSpy {
	spy := HandlerSpy{}
	spy.Fn = func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
		spy.CallCount++
		spy.Calls = append(spy.Calls, Call{
			Time: time.Now(),
			Args: []interface{}{ctx, bot, update, ctrl},
		})
		fn(ctx, bot, update, ctrl)
	}

	return &spy
}

func NewErrorSpy(fn func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control)) *ErrorSpy {
	spy := ErrorSpy{}
	spy.Fn = func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
		spy.CallCount++
		spy.Calls = append(spy.Calls, Call{
			Time: time.Now(),
			Args: []interface{}{ctx, bot, update, err, ctrl},
		})
		fn(ctx, bot, update, err, ctrl)
	}

	return &spy
}

func (s Spy) CalledBefore(another WithCalls) bool {
	for _, call := range s.Calls {
		for _, other := range another.GetCalls() {
			if call.Time.After(other.Time) {
				return false
			}
		}
	}

	return true
}

func (s Spy) CalledAfter(another Spy) bool {
	for _, call := range s.Calls {
		for _, other := range another.GetCalls() {
			if call.Time.Before(other.Time) {
				return false
			}
		}
	}

	return true
}

func (s *Spy) IsNeverCalled() bool {
	return s.CallCount == 0
}

func TestRouter(t *testing.T) {
	g := Goblin(t)

	g.Describe("Router", func() {
		var bot tgbotapi.BotAPI
		var update tgbotapi.Update

		g.BeforeEach(func() {
			bot = tgbotapi.BotAPI{}
			update = tgbotapi.Update{}
		})

		g.It("Should UseFunc()  okay", func() {
			var called bool
			router := NewRouter()
			router.UseFunc(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				called = true
				ctrl.Next()
			})

			router.OnUpdate(context.Background(), &bot, &update)
			time.Sleep(1 * time.Millisecond)
			g.Assert(called).IsTrue()
		})

		g.It("Should call in chain", func() {
			AHandler := NewHandlerSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				ctrl.Next()
			})
			BHandler := NewHandlerSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				ctrl.Error(fmt.Errorf("BHandler error"))
			})
			CHandler := NewHandlerSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctrl *Control) {
				ctrl.Next()
			})
			AErrorHandler := NewErrorSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				ctrl.Error(err)
			})
			BErrorHandler := NewErrorSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				ctrl.Error(fmt.Errorf("BErrorHandler error"))
			})
			CErrorHandler := NewErrorSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				ctrl.Next()
			})
			DErrorHandler := NewErrorSpy(func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error, ctrl *Control) {
				ctrl.Next()
			})

			router := NewRouter()
			router.UseFunc(AHandler.Fn, BHandler.Fn, CHandler.Fn)
			router.UseErrFunc(AErrorHandler.Fn, BErrorHandler.Fn, CErrorHandler.Fn, DErrorHandler.Fn)

			err := router.OnUpdate(context.Background(), &bot, &update)

			if err != nil {
				g.Fail(err)
				return
			}

			time.Sleep(1 * time.Millisecond)

			g.Assert(AHandler.CallCount).Eql(1)
			g.Assert(BHandler.CallCount).Eql(1)
			g.Assert(CHandler.CallCount).Eql(0)
			g.Assert(AHandler.CalledBefore(BHandler)).IsTrue()

			g.Assert(AErrorHandler.CallCount).Eql(1)
			g.Assert(BErrorHandler.CallCount).Eql(1)
			g.Assert(CErrorHandler.CallCount).Eql(1)
			g.Assert(DErrorHandler.CallCount).Eql(0)
			g.Assert(AErrorHandler.CalledBefore(BErrorHandler)).IsTrue()
			g.Assert(BErrorHandler.CalledBefore(CErrorHandler)).IsTrue()

			if bErrArg, ok := BErrorHandler.Calls[0].Args[3].(error); !ok {
				fmt.Printf("%t", BErrorHandler.Calls[0].Args)
				g.Fail("Mismatched type: expected error")
			} else {
				g.Assert(bErrArg.Error()).Eql("BHandler error")
			}

			if cErrArg, ok := CErrorHandler.Calls[0].Args[3].(error); !ok {
				g.Fail("Mismatched type: expected error")
			} else {
				g.Assert(cErrArg.Error()).Eql("BErrorHandler error")
			}
		})

	})
}
