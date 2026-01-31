package models

import (
	"strings"
)

type CurrencyPair struct {
	baseCurrency  CurrencyCode
	quoteCurrency CurrencyCode
}

func NewCurrencyPair(baseCurrency, quoteCurrency string) (CurrencyPair, error) {
	bc, err := NewCurrencyCode(baseCurrency)
	if err != nil {
		return CurrencyPair{}, err
	}

	qc, err := NewCurrencyCode(quoteCurrency)
	if err != nil {
		return CurrencyPair{}, err
	}

	if bc == qc {
		return CurrencyPair{}, ErrSameCurrencies
	}

	return CurrencyPair{
		baseCurrency:  bc,
		quoteCurrency: qc,
	}, nil
}

func (p CurrencyPair) BaseCurrency() CurrencyCode {
	return p.baseCurrency
}

func (p CurrencyPair) QuoteCurrency() CurrencyCode {
	return p.quoteCurrency
}

func (p CurrencyPair) String() string {
	return p.baseCurrency.String() + "/" + p.quoteCurrency.String()
}

func ParseCurrencyPair(pairStr string) (CurrencyPair, error) {
	parts := strings.Split(pairStr, "/")
	if len(parts) != 2 {
		return CurrencyPair{}, ErrInvalidPairFormat
	}

	return NewCurrencyPair(parts[0], parts[1])
}
