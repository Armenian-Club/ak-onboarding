package bottg

import (
	"fmt"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strconv"
	"strings"
)

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
	_, err := app.bot.SendMessage(ctx, tu.Message(update.Message.Chat.ChatID(),
		"Привет, "+userName+" 👋! Выберите действие:").WithReplyMarkup(keyboard))
	if err != nil {
		return err
	}
	return nil
}

// HandleCallback --- обработка CallbackQuery
func (app *BotApp) HandleCallback(ctx *th.Context, cq telego.CallbackQuery) error {
	userID := cq.From.ID
	userName := cq.From.FirstName

	app.lock.Lock()
	user, ok := app.users[userID]
	if !ok {
		// Если пользователя нет в map, создаём его
		user = User{Name: userName, Scenario: ScenarioNone, ConvState: StateDefault, Username: cq.From.Username}
		app.users[userID] = user
	}
	app.lock.Unlock()

	// Определяем ChatID для сообщений
	var chatID telego.ChatID
	if cq.Message != nil {
		chatID = tu.ID(cq.Message.Message().Chat.ID)
	} else {
		// Сообщение недоступно, пропускаем
		return nil
	}

	// Обрабатываем callback
	switch {
	case cq.Data == "onboarding":
		err := app.caseOnbording(ctx, user, userID, chatID, userName)
		if err != nil {
			return err
		}
	case cq.Data == "info":
		err := app.caseInfo(ctx, user, userID, chatID)
		if err != nil {
			return err
		}
	case strings.HasPrefix(cq.Data, "approve_"):
		err := app.caseApprove(ctx, cq, chatID)
		if err != nil {
			return err
		}
	case strings.HasPrefix(cq.Data, "reject_"):
		err := app.caseReject(ctx, cq, chatID)
		if err != nil {
			return err
		}
	default:
		log.Printf("⚠️ Неизвестный callback: %s", cq.Data)
	}

	return nil
}

// HandleMessage --- обработка сообщений
func (app *BotApp) HandleMessage(ctx *th.Context, msg telego.Message) error {
	userID := msg.From.ID

	app.lock.RLock()
	user, ok := app.users[userID]
	app.lock.RUnlock()

	if !ok {
		// Если пользователя нет, создаём с дефолтными значениями
		user = User{
			Name:      msg.From.FirstName,
			Username:  msg.From.Username,
			Scenario:  ScenarioNone,
			ConvState: StateDefault,
		}
		app.lock.Lock()
		app.users[userID] = user
		app.lock.Unlock()
	}

	switch user.Scenario {
	case ScenarioOnboarding:
		err := app.handleOnboarding(ctx, msg, &user)
		if err != nil {
			return err
		}
	case ScenarioInfo:
		app.handleInfo(ctx, msg, app.bot, &user)
	default:
		_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), "Выберите действие через /start"))
		if err != nil {
			return err
		}
	}
	// Обновляем пользователя в map после изменения состояния
	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	return nil
}

// --- Вспомогательные функции ---

func (app *BotApp) safeEditMarkup(ctx *th.Context, chatID telego.ChatID, msgID int, markup *telego.InlineKeyboardMarkup) error {
	_, err := app.bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   msgID,
		ReplyMarkup: markup,
	})
	if err != nil {
		return err
	}
	return nil
}

// сброс состояния пользователя (после завершения онбординга)
func (app *BotApp) resetUser(userID int64) {
	app.lock.Lock()
	defer app.lock.Unlock()
	delete(app.users, userID)
	log.Printf("Пользователь %d удалён из map (resetUser)", userID)
}

// --- Callback кейсы ---

func (app *BotApp) caseOnbording(ctx *th.Context, user User, userID int64, chatID telego.ChatID, userName string) error {
	user.Scenario = ScenarioOnboarding
	user.ConvState = StateAskEmail

	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	_, err := app.bot.SendMessage(ctx, tu.Message(chatID, fmt.Sprintf("Отлично, %s! Введи, пожалуйста, свою почту:", userName)))
	if err != nil {
		return err
	}
	return nil
}

func (app *BotApp) caseInfo(ctx *th.Context, user User, userID int64, chatID telego.ChatID) error {
	user.Scenario = ScenarioInfo
	user.ConvState = StateDefault

	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	_, err := app.bot.SendMessage(ctx, tu.Message(chatID, "Какая-то инфа"))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(chatID, "Выберите действие через /start"))
	if err != nil {
		return err
	}
	return nil
}

func (app *BotApp) caseApprove(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) error {
	targetIDStr := strings.TrimPrefix(cq.Data, "approve_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		return err
	}

	// ✅ убираем кнопки у сообщения админа
	if cq.Message != nil {
		err = app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
		if err != nil {
			return err
		}
	}

	// уведомляем админа
	_, err = app.bot.SendMessage(ctx, tu.Message(chatID, "✅ Пользователь подтверждён."))
	if err != nil {
		return err
	}
	//ONBOARDING
	err = app.onboarder.Onboard(ctx, app.users[targetID].Email, app.users[targetID].Gmail)
	if err != nil {
		return err
	}

	// уведомляем пользователя
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "🎉 Твой онбординг подтверждён!"))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))
	if err != nil {
		return err
	}
	// ❗ удаляем юзера
	app.resetUser(targetID)
	return nil
}

func (app *BotApp) caseReject(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) error {
	targetIDStr := strings.TrimPrefix(cq.Data, "reject_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		return err
	}

	// ❌ убираем кнопки у сообщения админа
	if cq.Message != nil {
		err = app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
		if err != nil {
			return err
		}
	}

	// уведомляем админа
	_, err = app.bot.SendMessage(ctx, tu.Message(chatID, "❌ Пользователь отклонён."))
	if err != nil {
		return err
	}
	// уведомляем пользователя
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "❌ Администратор отклонил онбординг."))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))
	if err != nil {
		return err
	}

	// ❗ удаляем юзера
	app.resetUser(targetID)
	return nil
}

func AdminIdParse() int64 {
	adminIdInt, err := strconv.ParseInt(config.AdminID, 10, 64)
	if err != nil {
		panic(err)
	}
	return adminIdInt
}
