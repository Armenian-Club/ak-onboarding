package bottg

import (
	"context"
	"sync"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// ConvState --- Состояния внутри сценария ---
type ConvState uint

const (
	StateDefault ConvState = iota
	StateAskEmail
	StateConfirm
)

// Scenario --- Сценарии ---
type Scenario string

const (
	ScenarioNone       Scenario = ""
	ScenarioOnboarding Scenario = "onboarding"
	ScenarioInfo       Scenario = "info"
)

// User --- Пользователь ---
type User struct {
	Name      string
	Scenario  Scenario
	ConvState ConvState
	Email     string
	Gmail     string
}

// BotApp --- Приложение бота ---
type BotApp struct {
	bot   *telego.Bot
	users map[int64]User
	lock  sync.RWMutex
}

// NewBotApp --- Конструктор ---
func NewBotApp(bot *telego.Bot) *BotApp {
	return &BotApp{
		bot:   bot,
		users: make(map[int64]User),
	}
}

// Run --- Запуск приложения ---
func (app *BotApp) Run(ctx context.Context) {
	updates, _ := app.bot.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(app.bot, updates)

	// Привязка методов
	bh.Handle(app.HandleStart, th.CommandEqual("start"))
	bh.HandleCallbackQuery(app.HandleCallback)
	bh.HandleMessage(app.HandleMessage)

	defer func() { _ = bh.Stop() }()
	_ = bh.Start()
}
