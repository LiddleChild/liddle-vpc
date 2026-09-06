package main

import (
	"context"
	"fmt"

	"github.com/LiddleChild/liddle-vpc/gristguard/grist"
)

type Repository struct {
	documentID  string
	gristClient grist.Client
}

func NewRepository(documentID string, gristClient grist.Client) Repository {
	return Repository{
		documentID:  documentID,
		gristClient: gristClient,
	}
}

type listAllEndpointsRecord struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

func (repo Repository) ListAllEndpoints(ctx context.Context) ([]string, error) {
	query := `
	select method, endpoint
	from grist_api_permissions
	`

	req := grist.QueryWithParamsRequest{
		DocID: repo.documentID,
		SQL:   query,
	}

	resp, err := repo.gristClient.QueryWithParams[listAllEndpointsRecord](ctx, req)
	if err != nil {
		return nil, err
	}

	endpoints := make([]string, 0, len(resp.Records))
	for _, record := range resp.Records {
		endpoints = append(endpoints, fmt.Sprintf("%s %s", record.Fields.Method, record.Fields.Endpoint))
	}

	return endpoints, nil
}

type ListEndpointsByNameAndKeyParams struct {
	Name string
	Key  string
}

type listEndpointsByNameAndKeyRecord struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

func (repo Repository) ListEndpointsByNameAndKey(ctx context.Context, params ListEndpointsByNameAndKeyParams) ([]string, error) {
	query := `
	select
		p.method,
		p.endpoint
	from grist_api_permissions p
	join api_keys k
	join json_each(k.permissions) kp
		on kp.value = p.permission
	where k.name = ?
		and k.key = ?
	`

	req := grist.QueryWithParamsRequest{
		DocID: repo.documentID,
		SQL:   query,
		Args: []any{
			any(params.Name),
			any(params.Key),
		},
	}

	resp, err := repo.gristClient.QueryWithParams[listEndpointsByNameAndKeyRecord](ctx, req)
	if err != nil {
		return nil, err
	}

	endpoints := make([]string, 0, len(resp.Records))
	for _, record := range resp.Records {
		endpoints = append(endpoints, fmt.Sprintf("%s %s", record.Fields.Method, record.Fields.Endpoint))
	}

	return endpoints, nil
}
