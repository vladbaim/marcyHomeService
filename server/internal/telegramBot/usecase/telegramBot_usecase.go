package usecase

import (
	"context"
	"fmt"
	"marcyHomeService/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramBotUsecase struct {
	telegramBot    *tgbotapi.BotAPI
	sensorDataRepo domain.SensorDataRepository
}

func NewTelegramBotUsecase(telegramBot *tgbotapi.BotAPI, sensorDataRepo domain.SensorDataRepository) domain.TelegramBotUsecase {
	return &telegramBotUsecase{
		telegramBot:    telegramBot,
		sensorDataRepo: sensorDataRepo,
	}
}

func (t *telegramBotUsecase) SendChatMessage(ctx context.Context, chatId int64, message string) (sendedMessage tgbotapi.Message, err error) {
	sensorData, err := t.sensorDataRepo.GetLast(ctx)
	if err != nil {
		return
	}

	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		sensorData.CreatedAt.Year(), sensorData.CreatedAt.Month(), sensorData.CreatedAt.Day(),
		sensorData.CreatedAt.Hour(), sensorData.CreatedAt.Minute(), sensorData.CreatedAt.Second())
	sendedMessage, err = t.telegramBot.Send(
		tgbotapi.NewMessage(
			chatId,
			message+
				fmt.Sprintf("\nТемпература: %.1fC \n Влажность: %.1f%% \n Уровень CO2: %d", sensorData.Temperature, sensorData.Humidity, sensorData.CarbonDioxide)+
				"\nВремя показателей: "+
				formattedTime))
	return
}
