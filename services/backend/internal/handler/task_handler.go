package handler

import (
	"net/http"

	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
)

type TaskHandler struct {
	svc service.CurrencySyncService
}

func NewTaskHandler(svc service.CurrencySyncService) TaskHandler {
	return TaskHandler{
		svc: svc,
	}
}

func (h TaskHandler) TriggerCurrencyTask(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Run(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}
