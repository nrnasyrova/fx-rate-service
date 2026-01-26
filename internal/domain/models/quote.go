package models

import (
	"math"
)

type QuoteE6 int64

const quoteScaleE6 = 1_000_000

func NewQuoteE6FromFloat(v float64) QuoteE6 {
	return QuoteE6(math.Round(v * quoteScaleE6))
}
