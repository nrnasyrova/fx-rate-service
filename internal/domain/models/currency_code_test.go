package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCurrencyCode_Success(t *testing.T) {
	got, err := NewCurrencyCode("EUR")
	require.NoError(t, err)
	require.Equal(t, "EUR", got.String())
}

func TestNewCurrencyCode_InvalidCurrencyCode(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		_, err := NewCurrencyCode("EU")
		require.ErrorIs(t, err, ErrInvalidCurrencyCode)
	})

	t.Run("unsupported", func(t *testing.T) {
		_, err := NewCurrencyCode("RSD")
		require.ErrorIs(t, err, ErrInvalidCurrencyCode)
	})

	t.Run("empty", func(t *testing.T) {
		_, err := NewCurrencyCode("")
		require.ErrorIs(t, err, ErrInvalidCurrencyCode)
	})
}

func TestNewCurrencyCode_TrimSpace(t *testing.T) {
	got, err := NewCurrencyCode("  MXN ")
	require.NoError(t, err)
	require.Equal(t, "MXN", got.String())
}

func TestNewCurrencyCode_ToUpper(t *testing.T) {
	got, err := NewCurrencyCode("mxn")
	require.NoError(t, err)
	require.Equal(t, "MXN", got.String())
}
