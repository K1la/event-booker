package service

import (
	"context"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/google/uuid"
)

func (s *Service) ConfirmBookingPayment(ctx context.Context, eventID uuid.UUID) error {
	return s.db.ConfirmBookingPayment(ctx, eventID)
}
func (s *Service) CancelBooking(ctx context.Context, booking *dto.QueueMessage) error {
	return s.db.CancelBooking(ctx, booking)
}
