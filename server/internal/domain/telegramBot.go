package domain

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotUsecase interface {
	SendChatMessage(ctx context.Context, chatId int64, message string) (tgbotapi.Message, error)
}
