package models

type CurrencyPair struct {
	from CurrencyCode
	to   CurrencyCode
}

func NewCurrencyPair(from, to string) (CurrencyPair, error) {
	f, err := NewCurrencyCode(from)
	if err != nil {
		return CurrencyPair{}, err
	}

	t, err := NewCurrencyCode(to)
	if err != nil {
		return CurrencyPair{}, err
	}

	if from == to {
		return CurrencyPair{}, ErrSameCurrencies
	}

	return CurrencyPair{
		from: f,
		to:   t,
	}, nil
}

func (p CurrencyPair) From() CurrencyCode {
	return p.from
}

func (p CurrencyPair) To() CurrencyCode {
	return p.to
}
