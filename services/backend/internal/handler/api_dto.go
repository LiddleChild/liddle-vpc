package handler

import (
	"encoding/json"

	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
)

type ingestionRequest struct {
	Transactions []ingestionItem `json:"transactions"`
}

type ingestionItem struct {
	Version      string               `json:"version"`
	Notification json.RawMessage      `json:"notification"`
	Transaction  ingestionTransaction `json:"transaction"`
}

type ingestionTransaction struct {
	App      string `json:"app"`
	Type     string `json:"type"`
	Date     string `json:"date"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (dto ingestionRequest) toServiceRequest() (service.IngestionRequest, error) {
	transactions := make([]service.IngestionItem, 0, len(dto.Transactions))
	for _, item := range dto.Transactions {
		transactions = append(transactions, service.IngestionItem{
			Version:      item.Version,
			Notification: item.Notification.String(),
			Transaction: service.IngestedTransaction{
				App:      item.Transaction.App,
				Type:     item.Transaction.Type,
				Date:     item.Transaction.Date,
				Amount:   item.Transaction.Amount,
				Currency: item.Transaction.Currency,
			},
		})
	}

	return service.IngestionRequest{
		Transactions: transactions,
	}, nil
}
