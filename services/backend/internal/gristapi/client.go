package gristapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL *url.URL
	apiKey  string
}

func NewClient(baseURL *url.URL, apiKey string) Client {
	return Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

type QueryWithParamsRequest struct {
	DocID string
	SQL   string
	Args  []any
}

type QueryWithParamsResponse[T any] struct {
	Records []Record[T]
}

func (client Client) ListTableRecords(ctx context.Context, docID, tableID string) ([]TableRecord, error) {
	requestURL := client.baseURL.JoinPath(fmt.Sprintf("/api/docs/%s/tables/%s/records", url.PathEscape(docID), url.PathEscape(tableID))).String()
	response, err := client.get(ctx, requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("grist records request failed with status %d", response.StatusCode)
	}
	var result tableRecordsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Records, nil
}

func (client Client) UpdateTableRecords(ctx context.Context, docID, tableID string, records []RecordUpdate) error {
	requestURL := client.baseURL.JoinPath(fmt.Sprintf("/api/docs/%s/tables/%s/records", url.PathEscape(docID), url.PathEscape(tableID))).String()
	body, err := json.Marshal(recordsUpdateRequest{Records: records})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		buffer, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("grist update request failed with status: %d, body: %s", response.StatusCode, string(buffer))
	}

	return nil
}

func (client Client) ListTableColumns(ctx context.Context, docID, tableID string) ([]TableColumn, error) {
	requestURL := client.baseURL.JoinPath(fmt.Sprintf("/api/docs/%s/tables/%s/columns", url.PathEscape(docID), url.PathEscape(tableID))).String()
	response, err := client.get(ctx, requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("grist columns request failed with status %d", response.StatusCode)
	}
	var result tableColumnsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Columns, nil
}

func (client Client) get(ctx context.Context, requestURL string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	return http.DefaultClient.Do(request)
}

func (client Client) QueryWithParams[T any](ctx context.Context, req QueryWithParamsRequest) (QueryWithParamsResponse[T], error) {
	requestURL := client.baseURL.JoinPath(fmt.Sprintf("/api/docs/%s/sql", req.DocID)).String()

	body := queryWithParamsRequest{
		SQL:     req.SQL,
		Args:    req.Args,
		Timeout: nil,
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, buffer)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	token := fmt.Sprintf("Bearer %s", client.apiKey)
	httpRequest.Header.Set("Authorization", token)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode/100 != 2 {
		return QueryWithParamsResponse[T]{}, fmt.Errorf("grist SQL request failed with status %d", httpResponse.StatusCode)
	}

	var response queryWithParamsResponse[T]
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	return QueryWithParamsResponse[T]{Records: response.Records}, nil
}
