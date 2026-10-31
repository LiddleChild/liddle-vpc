package handler

import (
	"net/http"

	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
)

type APIHandler struct {
	ingestionSvc service.IngestionService
}

func NewAPIHandler(ingestionSvc service.IngestionService) APIHandler {
	return APIHandler{
		ingestionSvc: ingestionSvc,
	}
}

func (h APIHandler) GetModelIngestionInstruction(w http.ResponseWriter, _ *http.Request) {
	instruction := h.ingestionSvc.GetModelIngestionInstruction()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(instruction))
}
