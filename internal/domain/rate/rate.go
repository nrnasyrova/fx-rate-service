package rate

type Rate struct {
	Pair        CurrencyPair
	Quote       QuoteE6
	UpdatedAtMs int64
}
