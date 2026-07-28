package bottg

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (app *BotApp) HandleStart(ctx *th.Context, update telego.Update) error {
	userID := update.Message.From.ID
	userName := update.Message.From.FirstName

	app.lock.Lock()
	app.users[userID] = User{
		Name:      userName,
		Username:  update.Message.From.Username,
		Scenario:  ScenarioNone,
		ConvState: StateDefault,
	}
	app.lock.Unlock()

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "Пройти онбординг", CallbackData: "onboarding"},
				{Text: "Инструкции", CallbackData: "info"},
			},
		},
	}

	_, err := app.bot.SendMessage(
		ctx,
		tu.Message(update.Message.Chat.ChatID(), "Привет, "+userName+" 👋! Выберите действие:").WithReplyMarkup(keyboard),
	)
	if err != nil {
		return err
	}

	_, err = app.bot.SendMessage(ctx, tu.Message(update.Message.Chat.ChatID(), supportText+config.SysadminTag))
	return err
}

// HandleCallback --- обработка CallbackQuery
func (app *BotApp) HandleCallback(ctx *th.Context, cq telego.CallbackQuery) error {
	userID := cq.From.ID
	userName := cq.From.FirstName

	app.lock.Lock()
	_, ok := app.users[userID]
	if !ok {
		app.users[userID] = User{
			Name:      userName,
			Username:  cq.From.Username,
			Scenario:  ScenarioNone,
			ConvState: StateDefault,
		}
	}
	app.lock.Unlock()

	var chatID telego.ChatID
	if cq.Message != nil {
		chatID = tu.ID(cq.Message.Message().Chat.ID)
	} else {
		return nil
	}

	switch {
	case cq.Data == "onboarding":
		if err := app.caseOnboarding(ctx, userID, chatID, userName); err != nil {
			return err
		}
	case cq.Data == "info":
		if err := app.caseInfo(ctx, userID, chatID); err != nil {
			return err
		}
	case strings.HasPrefix(cq.Data, "approve_"):
		if err := app.caseApprove(ctx, cq, chatID); err != nil {
			return err
		}
	case strings.HasPrefix(cq.Data, "reject_"):
		if err := app.caseReject(ctx, cq, chatID); err != nil {
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
		if err := app.handleOnboarding(ctx, msg, &user); err != nil {
			return err
		}
	default:
		_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), chooseActionText))
		if err != nil {
			return err
		}
	}

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

func (app *BotApp) sendInfoMessages(ctx *th.Context, chatID telego.ChatID) error {
	messages := []string{
		startInfoPage,
		"Инструкции для настройки сервисов:",
		"Для настройки Mattermost: https://outline.armenianclub.org/s/9814ee83-3a0e-4e7d-872f-c767d2216558",
		"Для Google Drive: https://outline.armenianclub.org/s/30b3026a-b656-4b1f-9415-d775effdcf22",
		chooseActionText,
	}

	for _, text := range messages {
		_, err := app.bot.SendMessage(ctx, tu.Message(chatID, text))
		if err != nil {
			return err
		}
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

func (app *BotApp) caseOnboarding(ctx *th.Context, userID int64, chatID telego.ChatID, userName string) error {
	app.lock.Lock()
	user := app.users[userID]
	user.Scenario = ScenarioOnboarding
	user.ConvState = StateAskEmail
	app.users[userID] = user
	app.lock.Unlock()

	_, err := app.bot.SendMessage(ctx, tu.Message(chatID, fmt.Sprintf(getEmailText, userName)))
	if err != nil {
		return err
	}

	return nil
}

func (app *BotApp) caseInfo(ctx *th.Context, userID int64, chatID telego.ChatID) error {
	app.lock.Lock()
	user := app.users[userID]
	user.Scenario = ScenarioNone
	user.ConvState = StateDefault
	app.users[userID] = user
	app.lock.Unlock()

	return app.sendInfoMessages(ctx, chatID)
}

func (app *BotApp) caseApprove(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) error {
	targetIDStr := strings.TrimPrefix(cq.Data, "approve_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		return err
	}

	if cq.Message != nil {
		err = app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
		if err != nil {
			return err
		}
	}

	_, err = app.bot.SendMessage(ctx, tu.Message(chatID, adminApprovedUserText))
	if err != nil {
		return err
	}

	app.lock.RLock()
	targetUser, ok := app.users[targetID]
	app.lock.RUnlock()
	if !ok {
		return fmt.Errorf("user %d not found", targetID)
	}

	err = app.onboarder.Onboard(ctx, targetUser.Email, targetUser.Gmail)
	if err != nil {
		return err
	}

	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), userOnboardApproveText))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), startInfoPage))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), checkEmailText))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), instructionsForMM))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), instructionsForGD))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), chooseActionText))
	if err != nil {
		return err
	}

	app.resetUser(targetID)
	return nil
}

func (app *BotApp) caseReject(ctx *th.Context, cq telego.CallbackQuery, chatID telego.ChatID) error {
	targetIDStr := strings.TrimPrefix(cq.Data, "reject_")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		return err
	}

	if cq.Message != nil {
		err = app.safeEditMarkup(ctx, chatID, cq.Message.GetMessageID(), nil)
		if err != nil {
			return err
		}
	}

	_, err = app.bot.SendMessage(ctx, tu.Message(chatID, adminRejectUserText))
	if err != nil {
		return err
	}

	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), userOnboardRejectText))
	if err != nil {
		return err
	}
	_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(targetID), chooseActionText))
	if err != nil {
		return err
	}

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
