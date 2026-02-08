package handler

import (
	"errors"
	"fmt"
	"github.com/K1la/event-booker/internal/api/response"
	"github.com/K1la/event-booker/internal/repository"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

func (h *Handler) GetEventByID(c *ginext.Context) {
	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("missing or invalid event id")
		response.BadRequest(c, err)
		return
	}

	event, err := h.service.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			zlog.Logger.Error().Err(err).Msg("event not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to get event")
		response.Internal(c, fmt.Errorf("failed to get event: %w", err))
		return
	}

	zlog.Logger.Info().Interface("event", event).Msg("got event")
	response.OK(c, event)
}

func (h *Handler) GetEvents(c *ginext.Context) {
	events, err := h.service.GetEvents(c.Request.Context())
	if err != nil {
		if errors.Is(err, repository.ErrEventsNotFound) {
			zlog.Logger.Error().Err(err).Msg("events not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("could not get all events")
		response.Internal(c, err)
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET all events")
	response.OK(c, events)
}

func parseUUIDParam(c *ginext.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Interface(param, idStr).Msg("failed to parse UUID")
		return uuid.Nil, fmt.Errorf("invalid %s", param)
	}

	if id == uuid.Nil {
		zlog.Logger.Warn().Interface(param, id).Msg("missing UUID")
		return uuid.Nil, fmt.Errorf("missing %s", param)
	}

	return id, nil
}
