package handler

import (
	"context"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/google/uuid"
)

type ServiceI interface {
	CreateEvent(ctx context.Context, event *dto.CreateEvent) (*model.Event, error)
	CreateBooking(ctx context.Context, booking *dto.CreateBooking) (*model.Booking, error)
	ConfirmBookingPayment(ctx context.Context, eventID uuid.UUID) error
	CancelBooking(ctx context.Context, bookingID *dto.QueueMessage) error

	GetEventByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error)
	GetEvents(ctx context.Context) ([]*model.Event, error)
}
