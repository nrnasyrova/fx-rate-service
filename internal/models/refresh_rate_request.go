package models

import "time"

type RefreshRateRequest struct {
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
