package rate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCurrencyPair_Success(t *testing.T) {
	from := "EUR"
	to := "USD"
	cp, err := NewCurrencyPair(from, to)

	require.NoError(t, err)
	require.Equal(t, from, cp.from.String())
	require.Equal(t, to, cp.to.String())
}

func TestNewCurrencyPair_FromEqualTo_ErrSameCurrencies(t *testing.T) {
	from, to := "EUR", "EUR"

	_, err := NewCurrencyPair(from, to)

	require.ErrorIs(t, err, ErrSameCurrencies)
}
