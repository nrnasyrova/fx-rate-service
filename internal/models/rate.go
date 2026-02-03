package models

import "time"

type Rate struct {
	Pair      CurrencyPair
	Quote     ValueE6
	UpdatedAt time.Time
}
