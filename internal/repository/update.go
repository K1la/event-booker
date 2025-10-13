package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func (r *Postgres) ConfirmBookingPayment(ctx context.Context, bookingID uuid.UUID) error {
	query := `UPDATE bookings
	SET status = 'confirmed',
	    updated_at = NOW()
	WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, bookingID).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBookingNotFoundOrAlreadyConfirmed
		}
		return fmt.Errorf("failed to confirm booking payment: %w", err)
	}

	return nil
}
func (r *Postgres) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cancelQuery := `
		UPDATE bookings
		SET status = 'cancelled',
		    updated_at = NOW()
		WHERE id = $1
		RETURNING event_id;
    `

	var eventID uuid.UUID
	err = tx.QueryRowContext(ctx, cancelQuery, bookingID).Scan(&eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBookingNotFoundOrAlreadyCancelled
		}
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	updateEventQuery := `
 		UPDATE events
		SET available_seats = available_seats + 1,
		    updated_at = NOW()
 		WHERE id = $1;
	`
	_, err = tx.ExecContext(ctx, updateEventQuery, eventID)
	if err != nil {
		return fmt.Errorf("failed to update event seats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
