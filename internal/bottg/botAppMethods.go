package bottg

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
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
		app.lock.Unlock()
		return nil
	}

	switch {
	case cq.Data == "onboarding":
		user.Scenario = ScenarioOnboarding
		user.ConvState = StateAskEmail
		app.users[userID] = user
		app.lock.Unlock()

		_, _ = app.bot.SendMessage(
			ctx,
			tu.Message(chatID, fmt.Sprintf("Отлично, %s! Введи, пожалуйста, свою почту:", userName)),
		)

	case cq.Data == "info":
		user.Scenario = ScenarioInfo
		user.ConvState = StateDefault
		app.users[userID] = user
		app.lock.Unlock()

		_, _ = app.bot.SendMessage(ctx, tu.Message(chatID, "Какая-то инфа"))
		_, _ = app.bot.SendMessage(ctx, tu.Message(chatID, "Выберите действие через /start"))

	case strings.HasPrefix(cq.Data, "approve_"):
		app.lock.Unlock()
		targetIDStr := strings.TrimPrefix(cq.Data, "approve_")
		targetID, _ := strconv.ParseInt(targetIDStr, 10, 64)

		// ✅ Убираем кнопки у сообщения админа
		_, _ = app.bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   cq.Message.GetMessageID(),
			ReplyMarkup: nil, // убираем клавиатуру
		})

		// уведомляем админа
		_, _ = app.bot.SendMessage(ctx, tu.Message(chatID, "✅ Пользователь подтверждён."))

		// уведомляем пользователя
		_, _ = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "🎉 Твой онбординг подтверждён!"))
		_, _ = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))

	case strings.HasPrefix(cq.Data, "reject_"):
		app.lock.Unlock()
		targetIDStr := strings.TrimPrefix(cq.Data, "reject_")
		targetID, _ := strconv.ParseInt(targetIDStr, 10, 64)

		// ❌ Убираем кнопки у сообщения админа
		_, _ = app.bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   cq.Message.GetMessageID(),
			ReplyMarkup: nil,
		})

		// уведомляем админа
		_, _ = app.bot.SendMessage(ctx, tu.Message(chatID, "❌ Пользователь отклонён."))

		// уведомляем пользователя
		_, _ = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "❌ Администратор отклонил онбординг."))
		_, _ = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))
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
