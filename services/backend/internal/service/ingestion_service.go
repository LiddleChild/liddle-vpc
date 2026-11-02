package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
	"github.com/LiddleChild/liddle-vpc/backend/internal/resources"
)

var ErrInvalidIngestionRequest = errors.New("invalid ingestion request")

type IngestionRepository interface {
	CreateIngestedTransactions(ctx context.Context, transactions []model.IngestedTransaction) error
}

type IngestionService struct {
	repo IngestionRepository
}

func NewIngestionService(repo IngestionRepository) IngestionService {
	return IngestionService{repo: repo}
}

func (svc IngestionService) GetModelIngestionInstruction() string {
	return resources.IngestionInstruction
}

func (svc IngestionService) IngestTransaction(ctx context.Context, request IngestionRequest) error {
	if len(request.Transactions) == 0 {
		return fmt.Errorf("%w: transactions must not be empty", ErrInvalidIngestionRequest)
	}

	now := time.Now()

	records := make([]model.IngestedTransaction, 0, len(request.Transactions))
	for _, item := range request.Transactions {
		amount, err := strconv.ParseFloat(item.Transaction.Amount, 10)
		if err != nil {
			return err
		}

		records = append(records, model.IngestedTransaction{
			App:          item.Transaction.App,
			Type:         item.Transaction.Type,
			Date:         item.Transaction.Date,
			Amount:       amount,
			Currency:     item.Transaction.Currency,
			Notification: item.Notification,
			CreatedAt:    item.Transaction.CreatedAt,
			IngestedAt:   now,
		})
	}

	return svc.repo.CreateIngestedTransactions(ctx, records)
}
