package bottg

import (
	"context"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func BotRun(ctx *context.Context, bot *telego.Bot) {

	updates, _ := bot.UpdatesViaLongPolling(*ctx, nil)

	for update := range updates {
		if update.Message != nil {
			chatID := tu.ID(update.Message.Chat.ID)
			keyboard := tu.Keyboard(
				tu.KeyboardRow(
					tu.KeyboardButton("Старт"),
					tu.KeyboardButton("Отправить на подтверждение"),
				),
			)
			message := tu.Message(
				chatID,
				"Hello Bro",
			).WithReplyMarkup(keyboard)
			_, _ = bot.SendMessage(*ctx, message)
		}
	}
}
