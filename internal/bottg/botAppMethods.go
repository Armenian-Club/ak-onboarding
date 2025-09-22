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

	app.safeSend(ctx, tu.Message(update.Message.Chat.ChatID(),
		"Привет, "+userName+" 👋! Выберите действие:").WithReplyMarkup(keyboard))

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
		app.caseOnbording(ctx, user, userID, chatID, userName)
	case cq.Data == "info":
		app.caseInfo(ctx, user, userID, chatID)
	case strings.HasPrefix(cq.Data, "approve_"):
		app.caseApprove(ctx, cq, chatID)
	case strings.HasPrefix(cq.Data, "reject_"):
		app.caseReject(ctx, cq, chatID)
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
		app.handleOnboarding(ctx, msg, &user)
	case ScenarioInfo:
		app.handleInfo(ctx, msg, app.bot, &user)
	default:
		app.safeSend(ctx, tu.Message(msg.Chat.ChatID(), "Выберите действие через /start"))
	}
	// Обновляем пользователя в map после изменения состояния
	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	return nil
}

// --- Вспомогательные функции ---

// безопасная отправка сообщений
func (app *BotApp) safeSend(ctx *th.Context, msg *telego.SendMessageParams) {
	_, err := app.bot.SendMessage(ctx, msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения (%v): %v", msg.Text, err)
	}
}

// безопасное редактирование разметки (убираем кнопки)
func (app *BotApp) safeEditMarkup(ctx *th.Context, chatID telego.ChatID, msgID int, markup *telego.InlineKeyboardMarkup) {
	_, err := app.bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   msgID,
		ReplyMarkup: markup,
	})
	if err != nil {
		log.Printf("Ошибка при редактировании клавиатуры: %v", err)
	}
}

// сброс состояния пользователя (после завершения онбординга)
func (app *BotApp) resetUser(userID int64) {
	app.lock.Lock()
	defer app.lock.Unlock()
	delete(app.users, userID)
	log.Printf("Пользователь %d удалён из map (resetUser)", userID)
}

// --- Callback кейсы ---

func (app *BotApp) caseOnbording(ctx *th.Context, user User, userID int64, chatID telego.ChatID, userName string) {
	user.Scenario = ScenarioOnboarding
	user.ConvState = StateAskEmail

	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	app.safeSend(ctx, tu.Message(chatID, fmt.Sprintf("Отлично, %s! Введи, пожалуйста, свою почту:", userName)))
}

func (app *BotApp) caseInfo(ctx *th.Context, user User, userID int64, chatID telego.ChatID) {
	user.Scenario = ScenarioInfo
	user.ConvState = StateDefault

	app.lock.Lock()
	app.users[userID] = user
	app.lock.Unlock()

	app.safeSend(ctx, tu.Message(chatID, "Какая-то инфа"))
	app.safeSend(ctx, tu.Message(chatID, "Выберите действие через /start"))
}

func (app *BotApp) caseApprove(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) {
	targetIDStr := strings.TrimPrefix(cq.Data, "approve_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка парсинга ID в approve:", err)
		return
	}

	// ✅ убираем кнопки у сообщения админа
	if cq.Message != nil {
		app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
	}

	// уведомляем админа
	app.safeSend(ctx, tu.Message(chatID, "✅ Пользователь подтверждён."))

	//ONBOARDING
	err = app.onboarder.Onboard(ctx, app.users[targetID].Email, app.users[targetID].Gmail)
	if err != nil {
		log.Fatal("Ошибка онбординга: " + err.Error())
	}

	// уведомляем пользователя
	app.safeSend(ctx, tu.Message(tu.ID(targetID), "🎉 Твой онбординг подтверждён!"))
	app.safeSend(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))

	// ❗ удаляем юзера
	app.resetUser(targetID)
}

func (app *BotApp) caseReject(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) {
	targetIDStr := strings.TrimPrefix(cq.Data, "reject_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка парсинга ID в reject:", err)
		return
	}

	// ❌ убираем кнопки у сообщения админа
	if cq.Message != nil {
		app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
	}

	// уведомляем админа
	app.safeSend(ctx, tu.Message(chatID, "❌ Пользователь отклонён."))

	// уведомляем пользователя
	app.safeSend(ctx, tu.Message(tu.ID(targetID), "❌ Администратор отклонил онбординг."))
	app.safeSend(ctx, tu.Message(tu.ID(targetID), "Выберите действие через /start"))

	// ❗ удаляем юзера
	app.resetUser(targetID)

}

func AdminIdParse() int64 {
	adminIdInt, err := strconv.ParseInt(config.AdminID, 10, 64)
	if err != nil {
		log.Fatal("Неправильный adminChatID в конфиге: " + err.Error())
	}
	return adminIdInt
}
