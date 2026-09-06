package grist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (c Client) QueryWithParams[T any](ctx context.Context, req QueryWithParamsRequest) (QueryWithParamsResponse[T], error) {
	url := c.baseURL.JoinPath(fmt.Sprintf("/api/docs/%s/sql", req.DocID)).String()

	body := queryWithParamsRequest{
		SQL:     req.SQL,
		Args:    req.Args,
		Timeout: nil,
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buffer)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	token := fmt.Sprintf("Bearer %s", c.apiKey)
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	if httpResp.StatusCode/100 != 2 {
		var gristErr Error
		if err := json.NewDecoder(httpResp.Body).Decode(&gristErr); err != nil {
			return QueryWithParamsResponse[T]{}, err
		}

		return QueryWithParamsResponse[T]{}, gristErr
	}

	var resp queryWithParamsResponse[T]
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return QueryWithParamsResponse[T]{}, err
	}

	return QueryWithParamsResponse[T]{
		Records: resp.Records,
	}, nil
}
