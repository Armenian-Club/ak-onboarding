package bottg

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"net/mail"
	"strings"
)

// --- Обработка сценария Onboarding ---
func handleOnboarding(ctx *th.Context, msg telego.Message, bot *telego.Bot, user *User) {
	var text string

	switch user.ConvState {
	case StateAskEmail:
		addr, err := mail.ParseAddress(msg.Text)
		if err != nil {
			text = "Неправильный формат почты, попробуйте ещё раз."
		} else if strings.HasSuffix(addr.Address, "@gmail.com") {
			user.Gmail = addr.Address
			user.ConvState = StateConfirm
			text = fmt.Sprintf("Вы указали почту(ы): %s %s. Верно? (Да/Нет)", user.Gmail, user.Email)

			keyboard := &telego.ReplyKeyboardMarkup{
				Keyboard: [][]telego.KeyboardButton{
					{{Text: "Да"}, {Text: "Нет"}},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			}

			_, _ = bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text).WithReplyMarkup(keyboard))
			return
		} else {
			user.Email = addr.Address
			text = "Для работы сервисов Google введите, пожалуйста, Gmail."
		}

	case StateConfirm:
		if msg.Text == "Да" {
			text = "Спасибо! Онбординг завершен ✅\nВыберите действие через /start"
			user.ConvState = StateDefault
			user.Scenario = ScenarioNone
		} else {
			text = "Хорошо, давайте попробуем ещё раз. Введите Gmail:"
			user.ConvState = StateAskEmail
		}
	}

	if text != "" {
		_, _ = bot.SendMessage(ctx, tu.Message(msg.Chat.ChatID(), text))
	}
}

// --- Обработка сценария Info ---
func handleInfo(ctx *th.Context, msg telego.Message, bot *telego.Bot, user *User) {

	//тут можно задать любые действия, но тут пока пусто

	user.Scenario = ScenarioNone
	user.ConvState = StateDefault
}
