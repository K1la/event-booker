package service

import (
	"context"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/google/uuid"
)

type DBRepo interface {
	CreateEvent(ctx context.Context, event *dto.CreateEvent) (*model.Event, error)
	CreateBooking(ctx context.Context, booking *dto.CreateBooking) (*model.Booking, error)
	ConfirmBookingPayment(ctx context.Context, bookingID uuid.UUID) error
	CancelBooking(ctx context.Context, booking *dto.QueueMessage) error

	//DeleteBooking(ctx context.Context, msg dto.QueueMessage) error

	GetEventByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error)
	GetEvents(ctx context.Context) ([]*model.Event, error)

	GetBookingByID(ctx context.Context, id uuid.UUID) (*dto.Booking, error)
}

type RabbitMQ interface {
	Publish(booking dto.QueueMessage) error
	Consume(ctx context.Context) (<-chan []byte, error)
}

type Sender interface {
	SendToTelegram(telegramId int, text string) error
}
