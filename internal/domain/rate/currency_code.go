package rate

import "strings"

type CurrencyCode string

const (
	USD CurrencyCode = "USD"
	EUR CurrencyCode = "EUR"
	MXN CurrencyCode = "MXN"
)

func NewCurrencyCode(s string) (CurrencyCode, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	switch CurrencyCode(s) {
	case USD, EUR, MXN:
		return CurrencyCode(s), nil
	default:
		return "", ErrInvalidCurrencyCode
	}
}

func (c CurrencyCode) String() string {
	return string(c)
}
