package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
)

type ExpenseHandler struct {
	svc service.ExpenseService
}

func NewExpenseHandler(svc service.ExpenseService) ExpenseHandler {
	return ExpenseHandler{
		svc: svc,
	}
}

func (h ExpenseHandler) GetExpenseCatalog(w http.ResponseWriter, r *http.Request) {
	catalog, err := h.svc.GetExpenseCatalog(r.Context())
	if err != nil {
		slog.Error("failed to get expense catalog", slog.String("error", err.Error()))
		writeError(w, http.StatusBadGateway, "expense catalog is unavailable")
		return
	}

	writeJSON(w, http.StatusOK, newCatalogResponse(catalog))
}

func (h ExpenseHandler) GetMonthlyExpenseSummary(w http.ResponseWriter, r *http.Request) {
	filters, err := parseExpenseFilters(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid query parameters")
		return
	}

	summary, err := h.svc.GetMonthlyExpenseSummary(r.Context(), filters)
	switch {
	case errors.Is(err, service.ErrInvalidExpenseRequest):
		writeError(w, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		slog.Error("failed to get monthly expense summary", slog.String("error", err.Error()))
		writeError(w, http.StatusBadGateway, "monthly expense summary is unavailable")
		return
	}

	writeJSON(w, http.StatusOK, newMonthlySummaryResponse(summary))
}

func parseExpenseFilters(query url.Values) (model.ExpenseFilters, error) {
	for _, key := range []string{"month", "accounts", "categories", "tags"} {
		if !query.Has(key) {
			return model.ExpenseFilters{}, errors.New("required query parameter is missing")
		}
	}

	monthValues := query["month"]
	if len(monthValues) != 1 || monthValues[0] == "" {
		return model.ExpenseFilters{}, errors.New("month must be provided once")
	}

	filters := model.ExpenseFilters{Month: monthValues[0]}
	accountValues, err := parseCSVFilters(query, "accounts")
	if err != nil {
		return model.ExpenseFilters{}, err
	}

	filters.Accounts = make([]int, 0, len(accountValues))
	for _, value := range accountValues {
		accountID, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return model.ExpenseFilters{}, errors.New("account ID must be an integer")
		}
		filters.Accounts = append(filters.Accounts, accountID)
	}

	filters.Categories, err = parseCSVFilters(query, "categories")
	if err != nil {
		return model.ExpenseFilters{}, err
	}

	filters.Tags, err = parseCSVFilters(query, "tags")
	if err != nil {
		return model.ExpenseFilters{}, err
	}

	return filters, nil
}

func parseCSVFilters(query url.Values, key string) ([]string, error) {
	values := query[key]
	if len(values) != 1 {
		return nil, errors.New("filter must be provided once")
	}

	if values[0] == "" {
		return []string{}, nil
	}

	reader := csv.NewReader(strings.NewReader(values[0]))
	filters, err := reader.ReadAll()
	if err != nil || len(filters) != 1 {
		return nil, errors.New("filter must be valid CSV")
	}

	for index, value := range filters[0] {
		filters[0][index] = strings.TrimSpace(value)
		if filters[0][index] == "" {
			return nil, errors.New("filter values cannot be empty")
		}
	}

	return filters[0], nil
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
