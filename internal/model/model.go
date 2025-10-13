package model

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	EventAt        time.Time `json:"event_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Bookings       []Booking `json:"bookings,omitempty"`
}

type Booking struct {
	ID          uuid.UUID `json:"id"`
	EventID     uuid.UUID `json:"event_id"`
	PlacesCount int       `json:"places_count"`
	Status      string    `json:"status"`
	TelegramID  int       `json:"telegram_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
