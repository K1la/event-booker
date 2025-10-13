package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/google/uuid"
)

func (r *Postgres) GetEventByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error) {
	query := `SELECT * FROM events WHERE id = $1`

	var event model.Event
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.EventAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get event from db: %w", err)
	}
	return &event, nil
}

func (r *Postgres) GetEvents(ctx context.Context) ([]*model.Event, error) {
	query := `
	SELECT 
		e.id,
		e.title,
		e.total_seats,
		e.available_seats,
		e.event_at,
		e.created_at,
		e.updated_at,
		COALESCE(json_agg(json_build_object(
			'id', b.id,
			'event_id', b.event_id,
			'status', b.status,
			'telegram_id', b.telegram_id,
			'created_at', b.created_at,
			'updated_at', b.updated_at
			)) FILTER (WHERE b.id IS NOT NULL), '[]'
		) AS bookings
	FROM events e
	LEFT JOIN bookings b ON b.event_id = e.id
	GROUP BY e.id
	ORDER BY e.created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get events from db: %w", err)
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var e model.Event
		var bookingsJSON []byte

		if err = rows.Scan(
			&e.ID,
			&e.Title,
			&e.TotalSeats,
			&e.AvailableSeats,
			&e.EventAt,
			&e.CreatedAt,
			&e.UpdatedAt,
			&bookingsJSON,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Декодируем JSON → []Booking
		if err = json.Unmarshal(bookingsJSON, &e.Bookings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bookings: %w", err)
		}

		events = append(events, &e)
	}

	if len(events) == 0 {
		return nil, ErrEventsNotFound
	}

	return events, nil
}

func (r *Postgres) GetBookingByID(ctx context.Context, id uuid.UUID) (*dto.Booking, error) {
	query := `SELECT b.*, e.title
	FROM bookings b
	JOIN events e on e.id = b.event_id
	WHERE b.id = $1`

	var booking dto.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.Status,
		&booking.TelegramID,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.EventTitle,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSuchBooking
		}
		return nil, fmt.Errorf("failed to get booking from db: %w", err)
	}
	return &booking, nil
}
