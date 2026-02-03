package models

import "time"

type RefreshRequest struct {
	ID        string
	Pair      CurrencyPair
	ValueE6   *ValueE6
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	ErrorMsg  *string
}

type Status string

const (
	Processing Status = "processing"
	Success    Status = "success"
	Error      Status = "error"
)

func (s Status) String() string {
	return string(s)
}
