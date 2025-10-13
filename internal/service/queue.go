package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/repository"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) StartWorker(ctx context.Context) {
	zlog.Logger.Info().Msg("started worker to cancel not paid bookings")
	go func() {
		s.consumeMessages(ctx)
	}()
}

func (s *Service) consumeMessages(ctx context.Context) {
	msgs, err := s.rbmq.Consume(ctx)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to consume messages from queue")
		return
	}

	for msg := range msgs {
		if err = s.handleQueueMessage(ctx, msg); err != nil {
			zlog.Logger.Error().Err(err).Msg("failed to handle message from queue")
			continue
		}
		zlog.Logger.Info().Msgf("processed message from queue")
	}
}

func (s *Service) handleQueueMessage(ctx context.Context, msgData []byte) error {
	var msg dto.QueueMessage
	if err := json.Unmarshal(msgData, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal queueMessage: %w", err)
	}

	bookingInfo, err := s.db.GetBookingByID(ctx, msg.BookingID)
	if err != nil {
		return err
	}

	if bookingInfo.Status == repository.StatusConfirmed {
		zlog.Logger.Info().Msgf("booking is already confirmed, id: %s", msg.BookingID)
		return nil
	}

	if err = s.db.CancelBooking(ctx, &msg); err != nil {
		return err
	}

	if bookingInfo.TelegramID != 0 {
		tgMsg := fmt.Sprintf("Your booking to event (%v) was cancelled due to unpaid status", bookingInfo.EventTitle)
		if err = s.sender.SendToTelegram(bookingInfo.TelegramID, tgMsg); err != nil {
			return err
		}
	}

	return nil
}
