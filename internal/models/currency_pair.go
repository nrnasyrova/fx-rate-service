package models

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

	if baseCurrency == quoteCurrency {
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
