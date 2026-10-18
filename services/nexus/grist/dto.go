package grist

type Error struct {
	ErrorMessage string `json:"error"`
}

func (err Error) Error() string {
	return err.ErrorMessage
}

type Record[T any] struct {
	Fields T `json:"fields"`
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
