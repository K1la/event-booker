package service

import (
	"context"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/model"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) CreateEvent(ctx context.Context, event *dto.CreateEvent) (*model.Event, error) {
	return s.db.CreateEvent(ctx, event)
}

func (s *Service) CreateBooking(ctx context.Context, booking *dto.CreateBooking) (*model.Booking, error) {
	createBooking, err := s.db.CreateBooking(ctx, booking)
	if err != nil {
		return nil, err
	}

	var msg dto.QueueMessage
	msg.BookingID = createBooking.ID
	msg.PlacesCount = createBooking.PlacesCount

	if err = s.rbmq.Publish(msg); err != nil {
		return nil, err
	}

	zlog.Logger.Info().Msgf("successfule publish to queue & created Booking: %v", createBooking)
	return createBooking, nil
}
