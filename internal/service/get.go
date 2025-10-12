package service

import (
	"context"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetEventByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error) {
	return s.db.GetEventByID(ctx, eventID)
}
func (s *Service) GetEvents(ctx context.Context) ([]*model.Event, error) {
	return s.db.GetEvents(ctx)
}

func (s *Service) GetBookingByID(ctx context.Context, eventID uuid.UUID) (*dto.Booking, error) {
	return s.GetBookingByID(ctx, eventID)
}
