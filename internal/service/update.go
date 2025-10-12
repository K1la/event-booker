package service

import (
	"context"
	"github.com/google/uuid"
)

func (s *Service) ConfirmBookingPayment(ctx context.Context, bookingID uuid.UUID) error {
	return s.db.ConfirmBookingPayment(ctx, bookingID)
}
func (s *Service) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	return s.db.CancelBooking(ctx, bookingID)
}
