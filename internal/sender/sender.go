package sender

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wb-go/wbf/zlog"
)

type TelegramSender struct {
	botApi *tgbotapi.BotAPI
}

func New() *TelegramSender {
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("could not connect to telegram api: %v", err)
	}

	botApi.Debug = false

	return &TelegramSender{
		botApi: botApi,
	}
}

func (t *TelegramSender) SendToTelegram(telegramId int, text string) error {
	msg := tgbotapi.NewMessage(int64(telegramId), text)
	_, err := t.botApi.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send message to telegram user: %w", err)
	}

	zlog.Logger.Info().Msgf("message to telegram user with id: %d was sent successfylly: ", telegramId)
	return nil
}
