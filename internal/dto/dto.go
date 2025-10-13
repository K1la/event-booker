package dto

import (
	"github.com/google/uuid"
	"time"
)

type Booking struct {
	ID         uuid.UUID `json:"id"`
	EventID    uuid.UUID `json:"event_id"`
	EventTitle string    `json:"event_title"`
	TelegramID int       `json:"telegram_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type QueueMessage struct {
	BookingID   uuid.UUID `json:"booking_id"`
	PlacesCount int       `json:"places_count"`
}

type CreateEvent struct {
	Title      string    `json:"title"`
	EventAt    time.Time `json:"event_at"`
	TotalSeats int       `json:"total_seats"`
}

type CreateBooking struct {
	EventID     uuid.UUID `json:"event_id,omitempty"`
	TelegramID  int       `json:"telegram_id"`
	PlacesCount int       `json:"places_count"`
}
