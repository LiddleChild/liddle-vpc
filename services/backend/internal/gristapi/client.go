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

func (c Client) ListTableRecords[T any](ctx context.Context, docID string, tableID string) ([]TableRecord[T], error) {
	response, err := c.request(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/docs/%s/tables/%s/records", url.PathEscape(docID), url.PathEscape(tableID)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result tableRecordsResponse[T]
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result.Records, nil
}

func (c Client) UpdateTableRecords[T any](ctx context.Context, docID string, tableID string, records []RecordUpdate[T]) error {
	body, err := json.Marshal(recordsUpdateRequest[T]{Records: records})
	if err != nil {
		return err
	}

	_, err = c.request(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/api/docs/%s/tables/%s/records", url.PathEscape(docID), url.PathEscape(tableID)),
		bytes.NewReader(body),
	)

	return err
}

func (c Client) ListTableColumns[T any](ctx context.Context, docID string, tableID string) ([]TableColumn[T], error) {
	response, err := c.request(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/docs/%s/tables/%s/columns", url.PathEscape(docID), url.PathEscape(tableID)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result tableColumnsResponse[T]
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result.Columns, nil
}

func (c Client) QueryWithParams[T any](ctx context.Context, req QueryWithParamsRequest) (QueryWithParamsResponse[T], error) {
	body := queryWithParamsRequest{
		SQL:     req.SQL,
		Args:    req.Args,
		Timeout: nil,
	}

	buffer, err := json.Marshal(body)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	responseBody, err := c.request(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/docs/%s/sql", url.PathEscape(req.DocID)),
		bytes.NewReader(buffer),
	)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	var response queryWithParamsResponse[T]
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	return QueryWithParamsResponse[T]{Records: response.Records}, nil
}

func (c Client) request(ctx context.Context, method string, path string, body io.Reader) ([]byte, error) {
	requestURL := c.baseURL.JoinPath(path).String()

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read grist response body: %w", err)
	}

	if httpResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("grist request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	return respBody, nil
}
