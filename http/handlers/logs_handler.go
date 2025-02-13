package handlers

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// Local Packages
	errors "learn-go/errors"
	models "learn-go/models/logs"
)

type LogsService interface {
	InsertLog(context.Context, models.LogModel) (string, error)
}

type LogsHandler struct {
	svc LogsService
}

func NewLogsHandler(svc LogsService) *LogsHandler {
	return &LogsHandler{svc: svc}
}

func (h *LogsHandler) Insert(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var log models.LogModel
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		return nil, http.StatusBadRequest, errors.InvalidBodyErr(err)
	}
	if err := log.Validate(); err != nil {
		return nil, http.StatusBadRequest, errors.ValidationFailedErr(err)
	}

	id, err := h.svc.InsertLog(r.Context(), log)
	if err == nil {
		message := fmt.Sprintf("Log Inserted Successfully :: %s", id)
		return map[string]string{"message": message}, http.StatusOK, nil
	}
	return
}
