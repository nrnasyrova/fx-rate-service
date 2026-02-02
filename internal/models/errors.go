package models

import "errors"

var (
	ErrInvalidCurrencyCode = errors.New("invalid currency code")
	ErrSameCurrencies      = errors.New("currency pair must have different currencies")
	ErrNotFound            = errors.New("not found")
	ErrInvalidPairFormat   = errors.New("invalid currency pair format: expected BASE/QUOTE")
	ErrServiceOverloaded   = errors.New("service overloaded")
)
