package service

import (
	"context"
)

type DBRepo interface {
}

type RabbitMQ interface {
	Publish(booking dto.QueueMessage) error
	Consume(ctx context.Context) (<-chan []byte, error)
}

type Sender interface {
	SendToTelegram(telegramId int, text string) error
}
