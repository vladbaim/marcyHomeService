package http

import (
	"encoding/json"
	"fmt"
	"log"
	"marcyHomeService/internal/domain"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/julienschmidt/httprouter"
)

type TelegramBotHandler struct {
	TUsecase domain.TelegramBotUsecase
}

// NewTelegramBotHandler will initialize the articles/ resources endpoint
func NewTelegramBotHandler(tu domain.TelegramBotUsecase, router *httprouter.Router, whUrl string) {
	handler := &TelegramBotHandler{
		TUsecase: tu,
	}

	router.HandlerFunc(http.MethodPost, whUrl, handler.HandleTelegramWh)
}

func (t *TelegramBotHandler) HandleTelegramWh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var update tgbotapi.Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Telegram WebHook From: %+v Text: %+v\n", update.Message.From, update.Message.Text)

	t.TUsecase.SendChatMessage(r.Context(), update.Message.Chat.ID, fmt.Sprintf("Привет, %s! Вот такая погода)", update.Message.From))
}
