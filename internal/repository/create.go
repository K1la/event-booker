package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/lib/pq"
)

func (r *Postgres) CreateEvent(ctx context.Context, event *dto.CreateEvent) (*model.Event, error) {
	query := `
	INSERT INTO events(title, event_at, total_seats, available_seats)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at
	`

	var createdEvent model.Event
	err := r.db.QueryRowContext(ctx, query, event.Title, event.EventAt, event.TotalSeats, event.TotalSeats).Scan(
		&createdEvent.ID, &createdEvent.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create event in db: %w", err)
	}

	createdEvent.EventAt = event.EventAt
	createdEvent.Title = event.Title
	createdEvent.TotalSeats = event.TotalSeats
	createdEvent.AvailableSeats = event.TotalSeats

	return &createdEvent, nil
}

func (r *Postgres) CreateBooking(ctx context.Context, booking *dto.CreateBooking) (*model.Booking, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	var createdBooking model.Booking
	bookingsQuery := `INSERT INTO bookings(event_id, status, telegram_id)
	VALUES ($1, $2, $3) RETURNING id, created_at`

	err = tx.QueryRowContext(ctx, bookingsQuery, booking.EventID, statusPending, booking.TelegramID).Scan(
		&createdBooking.ID, &createdBooking.CreatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrNoSuchEvent
		}

		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	EventsQuery := `UPDATE events
	SET available_seats = available_seats - $1
	WHERE id = $2 AND available_seats >= 0`
	result, err := tx.ExecContext(ctx, EventsQuery, booking.PlacesCount, booking.EventID) // booking.PlacesCount,
	if err != nil {
		return nil, fmt.Errorf("failed to update booking in event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to update booking in event: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrNoSeatsAvailable
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	createdBooking.EventID = booking.EventID
	createdBooking.TelegramID = booking.TelegramID
	createdBooking.PlacesCount = booking.PlacesCount
	createdBooking.Status = statusPending

	return &createdBooking, nil
}
