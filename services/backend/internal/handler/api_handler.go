package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(instruction))
}

func (h APIHandler) IngestTransaction(w http.ResponseWriter, r *http.Request) {
	var request ingestionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("invalid ingestion request", "error", err.Error())
		writeError(w, http.StatusBadRequest, "invalid ingestion request")
		return
	}

	serviceRequest, err := request.toServiceRequest()
	if err != nil {
		slog.Error("invalid ingestion request", "error", err.Error())
		writeError(w, http.StatusBadRequest, "invalid ingestion request")
		return
	}

	if err := h.ingestionSvc.IngestTransaction(r.Context(), serviceRequest); err != nil {
		if errors.Is(err, service.ErrInvalidIngestionRequest) {
			slog.Error("invalid ingestion request", "error", err.Error())
			writeError(w, http.StatusBadRequest, "invalid ingestion request")
			return
		}

		slog.Error("failed to ingest transactions", "error", err)
		writeError(w, http.StatusBadGateway, "transactions could not be ingested")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
