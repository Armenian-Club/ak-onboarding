package bottg

import (
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// --- Обработка сценария Onboarding ---
func (app *BotApp) handleOnboarding(ctx *th.Context, msg telego.Message, user *User) error {
	var text string

	switch user.ConvState {
	case StateAskEmail:
		addr, err := mail.ParseAddress(msg.Text)
		user.Username = msg.From.Username

		if err != nil {
			text = incorrectEmailText
		} else if strings.HasSuffix(addr.Address, "@gmail.com") {
			user.Gmail = addr.Address

			if user.Email == "" {
				user.Email = addr.Address
			}

			user.ConvState = StateConfirm
			text = fmt.Sprintf(emailCheckText, user.Email, user.Gmail)

			keyboard := &telego.ReplyKeyboardMarkup{
				Keyboard: [][]telego.KeyboardButton{
					{{Text: "Да"}, {Text: "Нет"}},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			}

			_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text).WithReplyMarkup(keyboard))
			if err != nil {
				return err
			}

			log.Printf("Got %s: Email: %s, Gmail: %s", user.Name, user.Email, user.Gmail)
			return nil
		} else {
			if user.Email == "" {
				user.Email = addr.Address
			}
			text = "Для работы сервисов Google введите, пожалуйста, Gmail."
		}

	case StateConfirm:
		removeKeyboard := &telego.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		}

		if msg.Text == "Да" {
			user.ConvState = StateWaitAdmin

			text = "Спасибо! Отправил запрос администратору для подтверждения, ожидай ответа."

			_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text).WithReplyMarkup(removeKeyboard))
			if err != nil {
				return err
			}

			adminText := fmt.Sprintf(userWantOnboardText, user.Username)

			keyboard := &telego.InlineKeyboardMarkup{
				InlineKeyboard: [][]telego.InlineKeyboardButton{
					{
						{Text: "✅ Подтвердить", CallbackData: "approve_" + strconv.FormatInt(msg.Chat.ID, 10)},
						{Text: "❌ Отклонить", CallbackData: "reject_" + strconv.FormatInt(msg.Chat.ID, 10)},
					},
				},
			}

			_, err = app.bot.SendMessage(ctx, tu.Message(tu.ID(app.adminID), adminText).WithReplyMarkup(keyboard))
			if err != nil {
				return err
			}

			return nil
		}

		if msg.Text == "Нет" {
			user.Email = ""
			user.Gmail = ""
			user.ConvState = StateAskEmail
			text = "Хорошо, давайте попробуем ещё раз. Введите почту:"

			_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text).WithReplyMarkup(removeKeyboard))
			if err != nil {
				return err
			}

			return nil
		}

		text = "Пожалуйста, выберите: Да или Нет."

	default:
		panic("unhandled default case")
	}

	_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text))
	if err != nil {
		return err
	}

	return nil
}
