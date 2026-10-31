package repository

import (
	"context"

	"github.com/LiddleChild/liddle-vpc/backend/internal/gristapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
)

type IngestedTransactionRepository struct {
	documentID  string
	gristClient gristapi.Client
	tableID     string
}

func NewIngestedTransactionRepository(documentID string, tableID string, gristClient gristapi.Client) IngestedTransactionRepository {
	return IngestedTransactionRepository{
		documentID:  documentID,
		gristClient: gristClient,
		tableID:     tableID,
	}
}

func (repo IngestedTransactionRepository) CreateIngestedTransactions(ctx context.Context, transactions []model.IngestedTransaction) error {
	records := make([]gristapi.RecordAdd[ingestedTransactionFields], 0, len(transactions))
	for _, transaction := range transactions {
		records = append(records, gristapi.RecordAdd[ingestedTransactionFields]{
			Fields: ingestedTransactionFields{
				App:          transaction.App,
				Type:         transaction.Type,
				Date:         transaction.Date,
				Amount:       transaction.Amount,
				Currency:     transaction.Currency,
				Notification: transaction.Notification,
			},
		})
	}

	return repo.gristClient.AddTableRecords(ctx, repo.documentID, repo.tableID, records)
}
