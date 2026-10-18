package gristapi

type Record[T any] struct {
	Fields T `json:"fields"`
}

type TableRecord struct {
	ID     int            `json:"id"`
	Fields map[string]any `json:"fields"`
}

type TableColumn struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

type tableRecordsResponse struct {
	Records []TableRecord `json:"records"`
}

type tableColumnsResponse struct {
	Columns []TableColumn `json:"columns"`
}

type queryWithParamsRequest struct {
	SQL     string `json:"sql"`
	Args    []any  `json:"args,omitempty"`
	Timeout *int   `json:"timeout,omitempty"`
}

type queryWithParamsResponse[T any] struct {
	Statement string      `json:"statement"`
	Records   []Record[T] `json:"records"`
}
