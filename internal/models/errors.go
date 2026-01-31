package models

import "errors"

var (
	ErrInvalidCurrencyCode = errors.New("invalid currency code")
	ErrSameCurrencies      = errors.New("currency pair must have different currencies")
)
