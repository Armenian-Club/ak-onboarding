package bottg

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"net/mail"
	"strconv"
	"strings"
)

// --- Обработка сценария Onboarding ---
func (app *BotApp) handleOnboarding(ctx *th.Context, msg telego.Message, user *User) error {
	var text string
	switch user.ConvState {
	case StateAskEmail:
		addr, uncorrectEmail := mail.ParseAddress(msg.Text)
		user.Username = msg.From.Username
		if uncorrectEmail != nil {
			text = "Неправильный формат почты, попробуйте ещё раз."
		} else if strings.HasSuffix(addr.Address, "@gmail.com") {
			user.Gmail = addr.Address
			user.ConvState = StateConfirm
			text = fmt.Sprintf("Вы указали почту(ы): %s %s. Верно? (Да/Нет)", user.Email, user.Gmail)

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
			if user.Email == "" {
				user.Email = addr.Address
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
		// создаём объект для удаления кнопок у пользователя
		removeKeyboard := &telego.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		}

		if msg.Text == "Да" {
			user.ConvState = StateWaitAdmin
			text = "Спасибо! Отправил запрос администратору для подтверждения, ожидай ответа."

			// Убираем кнопки у пользователя
			_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text).WithReplyMarkup(removeKeyboard))
			if err != nil {
				return err
			}
			// Отправляем админу заявку

			adminText := fmt.Sprintf(
				"Пользователь @%s хочет пройти онбординг",
				user.Username,
			)

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

		} else if msg.Text == "Нет" {
			// Пользователь сказал "Нет" → возвращаемся на ввод
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
	default:
		panic("unhandled default case")
	}

	_, err := app.bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text))
	if err != nil {
		return err
	}
	return nil
}

// --- Обработка сценария Info ---
func (app *BotApp) handleInfo(ctx *th.Context, msg telego.Message, bot *telego.Bot, user *User) {

	//тут можно задать любые действия, но тут пока пусто

	user.Scenario = ScenarioNone
	user.ConvState = StateDefault
}
