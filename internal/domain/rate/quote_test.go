package rate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQuoteE6FromFloat_ExactSixDecimals(t *testing.T) {
	q := NewQuoteE6FromFloat(5.109727)
	require.Equal(t, QuoteE6(5109727), q)
}

func TestNewQuoteE6FromFloat_RoundsMoreThanSix(t *testing.T) {
	q := NewQuoteE6FromFloat(1.2345675)
	require.Equal(t, QuoteE6(1234568), q)
}
