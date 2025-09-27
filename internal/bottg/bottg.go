package bottg

import (
	"context"
	"github.com/Armenian-Club/ak-onboarding/internal/app"
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
	StateWaitAdmin
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
	Username  string
	Name      string
	Scenario  Scenario
	ConvState ConvState
	Email     string
	Gmail     string
}

// BotApp --- Приложение бота ---
type BotApp struct {
	bot       *telego.Bot
	users     map[int64]User
	lock      sync.RWMutex
	adminID   int64
	onboarder app.Onboarder
}

// NewBotApp --- Конструктор ---
func NewBotApp(bot *telego.Bot, onboarder app.Onboarder) *BotApp {
	return &BotApp{
		bot:       bot,
		users:     make(map[int64]User),
		adminID:   AdminIdParse(),
		onboarder: onboarder,
	}
}

// Run --- Запуск приложения ---
func (app *BotApp) Run(ctx context.Context) error {
	updates, err := app.bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}
	bh, err := th.NewBotHandler(app.bot, updates)
	if err != nil {
		return err
	}
	// Привязка методов
	bh.Handle(app.HandleStart, th.CommandEqual("start"))
	bh.HandleCallbackQuery(app.HandleCallback)
	bh.HandleMessage(app.HandleMessage)

	defer func() { _ = bh.Stop() }()
	err = bh.Start()
	if err != nil {
		return err
	}
	return nil
}
