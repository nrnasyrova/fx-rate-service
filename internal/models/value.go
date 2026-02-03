package models

import (
	"math"
)

type ValueE6 int64

const valueScaleE6 = 1_000_000

func NewValueE6FromFloat(v float64) ValueE6 {
	return ValueE6(math.Round(v * valueScaleE6))
}
