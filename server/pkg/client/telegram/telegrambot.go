package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotClient interface {
	SendToChat(chatId string, message string) error
}

func NewTelegramBotClient(appUrl string, telegramBotToken string) (bot *tgbotapi.BotAPI, whUrl string, err error) {
	bot, err = tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	whUrl = "/api/telegram/" + bot.Token
	wh, _ := tgbotapi.NewWebhook(appUrl + whUrl)

	_, err = bot.Request(wh)
	if err != nil {
		return
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return
}
