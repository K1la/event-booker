package handler

import (
	"errors"
	"github.com/K1la/event-booker/internal/api/response"
	"github.com/K1la/event-booker/internal/dto"
	"github.com/K1la/event-booker/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

// events/
func (h *Handler) CreateEvent(c *ginext.Context) {
	var createEvent dto.CreateEvent
	if err := c.ShouldBindJSON(&createEvent); err != nil {
		zlog.Logger.Error().Err(err).Msg("Bind json failed")
		response.Internal(c, err)
		return
	}

	event, err := h.service.CreateEvent(c.Request.Context(), &createEvent)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("CreateEvent failed")
		response.Internal(c, err)
		return
	}

	zlog.Logger.Info().Interface("event", event).Msg("CreateEvent success")
	response.OK(c, event)
}

// events/:id/book
func (h *Handler) CreateBooking(c *ginext.Context) {
	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("missing or invalid event id")
		response.BadRequest(c, err)
		return
	}

	var booking dto.CreateBooking
	booking.EventID = eventID
	if err = c.ShouldBindJSON(&booking); err != nil {
		zlog.Logger.Error().Err(err).Msg("bind json failed")
		response.BadRequest(c, err)
		return
	}

	booked, err := h.service.CreateBooking(c.Request.Context(), &booking)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("CreateBooking failed")
		response.Internal(c, err)
		return
	}

	zlog.Logger.Info().Interface("booked", booked).Msg("CreateBooking success")
	response.OK(c, booked)

}

// events/:id/confirm
func (h *Handler) ConfirmBookingPayment(c *ginext.Context) {
	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("missing or invalid event id")
		response.BadRequest(c, err)
		return
	}

	if err = h.service.ConfirmBookingPayment(c.Request.Context(), eventID); err != nil {
		if errors.Is(err, repository.ErrBookingNotFoundOrAlreadyConfirmed) {
			zlog.Logger.Error().Err(err).Msg("booking not found or already confirmed")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("ConfirmBookingPayment failed")
		response.Internal(c, err)
		return
	}

	zlog.Logger.Info().Interface("eventID", eventID).Msg("ConfirmBookingPayment success")
	response.OK(c, ginext.H{"status": "payment confirmed"})
}

// events/:id
func (h *Handler) CancelBooking(c *ginext.Context) {
	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("missing or invalid event id")
		response.BadRequest(c, err)
		return
	}

	if err = h.service.CancelBooking(c.Request.Context(), eventID); err != nil {
		if errors.Is(err, repository.ErrBookingNotFoundOrAlreadyCancelled) {
			zlog.Logger.Error().Err(err).Msg("booking not found or already canceled")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("Cancel booking failed")
		response.Internal(c, err)
		return
	}

	zlog.Logger.Info().Interface("eventID", eventID).Msg("Cancel booking success")
	response.OK(c, ginext.H{"status": "booking successfully cancelled"})
}
