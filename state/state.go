package state

type TimestampInfo struct {
	Timestamp string `json:"timestamp"`
	Readable  string `json:"readable"`
}

type PriceInfo struct {
	Timestamp string  `json:"timestamp"`
	Readable  string  `json:"readable"`
	Price     float64 `json:"price"`
}

type CorrelationInfo struct {
	Timestamp   string  `json:"timestamp"`
	Readable    string  `json:"readable"`
	SymbolPrice float64 `json:"symbol_price"`
	HedgePrice  float64 `json:"hedge_price"`
	Correlation float64 `json:"correlation"`
}

type Day struct {
	Month int
	Day   int
	Year  int
}

type PolygonOpenClose struct {
	Status    string  `json:"status"`
	From      string  `json:"from"`
	Symbol    string  `json:"symbol"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int     `json:"volume"`
	PreMarket float64 `json:"preMarket"`
}

type Stock struct {
	Day    Day
	Symbol string
}
