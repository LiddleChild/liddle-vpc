package gristapi

type Record[T any] struct {
	Fields T `json:"fields"`
}

type TableRecord[T any] struct {
	ID     int `json:"id"`
	Fields T   `json:"fields"`
}

type TableColumn[T any] struct {
	ID     string `json:"id"`
	Fields T      `json:"fields"`
}

type RecordUpdate[T any] struct {
	ID     int `json:"id"`
	Fields T   `json:"fields"`
}

type recordsUpdateRequest[T any] struct {
	Records []RecordUpdate[T] `json:"records"`
}

type tableRecordsResponse[T any] struct {
	Records []TableRecord[T] `json:"records"`
}

type tableColumnsResponse[T any] struct {
	Columns []TableColumn[T] `json:"columns"`
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
