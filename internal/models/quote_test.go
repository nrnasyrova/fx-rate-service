package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQuoteE6FromFloat_ExactSixDecimals(t *testing.T) {
	q := NewValueE6FromFloat(5.109727)
	require.Equal(t, ValueE6(5109727), q)
}

func TestNewQuoteE6FromFloat_RoundsMoreThanSix(t *testing.T) {
	q := NewValueE6FromFloat(1.2345675)
	require.Equal(t, ValueE6(1234568), q)
}
