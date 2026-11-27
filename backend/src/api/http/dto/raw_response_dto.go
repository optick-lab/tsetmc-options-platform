package dto

import "time"

// RawResponseListItemDto represents a row in the stored raw responses list (excluding the body).
type RawResponseListItemDto struct {
	ID         int64     `json:"id"`
	Endpoint   string    `json:"endpoint"`
	InsCode    *string   `json:"ins_code"`
	HttpStatus int       `json:"http_status"`
	IsBlocked  bool      `json:"is_blocked"`
	FetchedAt  time.Time `json:"fetched_at"`
	BodySize   int64     `json:"body_size"`
	TradeDate  *string   `json:"trade_date"`
}

// TradeItemDto represents a single trade parsed from TSETMC GetTradeHistory.
type TradeItemDto struct {
	NTran    int     `json:"n_tran"`
	Time     string  `json:"time"`
	Volume   int64   `json:"volume"`
	Price    float64 `json:"price"`
	Canceled bool    `json:"canceled"`
}

// TradeSummaryDto represents the aggregated metrics for a trade history response (excluding canceled trades).
type TradeSummaryDto struct {
	TradeCount  int     `json:"trade_count"`
	TotalVolume int64   `json:"total_volume"`
	TotalValue  float64 `json:"total_value"`
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
	FirstPrice  float64 `json:"first_price"`
	LastPrice   float64 `json:"last_price"`
}

// TradesResponseDto wraps the summary metrics and individual trade items.
type TradesResponseDto struct {
	Summary TradeSummaryDto `json:"summary"`
	Trades  []TradeItemDto  `json:"trades"`
}
