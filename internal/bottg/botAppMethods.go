package bottg

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleStart --- /start ---
func (app *BotApp) HandleStart(ctx *th.Context, update telego.Update) error {
	userID := update.Message.From.ID
	userName := update.Message.From.FirstName

	app.lock.Lock()
	app.users[userID] = User{Name: userName, Scenario: ScenarioNone, ConvState: StateDefault}
	app.lock.Unlock()

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "Пройти онбординг", CallbackData: "onboarding"},
				{Text: "Информация", CallbackData: "info"},
			},
		},
	}

	_, _ = app.bot.SendMessage(
		ctx,
		tu.Message(update.Message.Chat.ChatID(), fmt.Sprintf("Привет, %s! 👋 Выберите действие:", userName)).
			WithReplyMarkup(keyboard),
	)

	return nil
}

// HandleCallback --- CallbackQuery ---
func (app *BotApp) HandleCallback(ctx *th.Context, cq telego.CallbackQuery) error {
	userID := cq.From.ID
	userName := cq.From.FirstName

	app.lock.Lock()
	user := app.users[userID]

	var chatID telego.ChatID
	if cq.Message != nil && cq.Message.Message() != nil {
		chatID = tu.ID(cq.Message.Message().Chat.ID)
	} else {
		// Сообщение недоступно
		return nil
	}

	switch cq.Data {
	case "onboarding":
		user.Scenario = ScenarioOnboarding
		user.ConvState = StateAskEmail
		app.users[userID] = user
		app.lock.Unlock()

		_, _ = app.bot.SendMessage(
			ctx,
			tu.Message(chatID, fmt.Sprintf("Отлично, %s! Введи, пожалуйста, свою почту:", userName)),
		)
	case "info":
		user.Scenario = ScenarioInfo
		user.ConvState = StateDefault
		app.users[userID] = user
		app.lock.Unlock()

		_, _ = app.bot.SendMessage(
			ctx,
			tu.Message(chatID, "Какая-то инфа"),
		)
		_, _ = app.bot.SendMessage(
			ctx,
			tu.Message(chatID, "Выберите действие через /start"),
		)
	default:
		app.lock.Unlock()
	}

	return nil
}

// HandleMessage --- Сообщения ---
func (app *BotApp) HandleMessage(ctx *th.Context, msg telego.Message) error {
	userID := msg.From.ID

	app.lock.RLock()
	user := app.users[userID]
	app.lock.RUnlock()

	switch user.Scenario {
	case ScenarioOnboarding:
		handleOnboarding(ctx, msg, app.bot, &user)
	case ScenarioInfo:
		handleInfo(ctx, msg, app.bot, &user)
	default:
		_, _ = app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), "Выберите действие через /start"))
	}

	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	return nil
}
