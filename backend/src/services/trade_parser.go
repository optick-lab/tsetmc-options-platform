package services

import (
	"encoding/json"
	"fmt"
	"math"

	"backend/src/api/http/dto"
)

// tsetmcTradeRaw represents a single item in the TSETMC GetTradeHistory JSON response.
type tsetmcTradeRaw struct {
	NTran    int     `json:"nTran"`
	HEven    int     `json:"hEven"`
	QTitTran int64   `json:"qTitTran"`
	PTran    float64 `json:"pTran"`
	Canceled int     `json:"canceled"`
}

// tsetmcTradeHistoryPayload represents the root envelope of TSETMC GetTradeHistory.
type tsetmcTradeHistoryPayload struct {
	TradeHistory []tsetmcTradeRaw `json:"tradeHistory"`
}

// FormatHEven formats an integer time representation (e.g. 122940 or 94033) into "HH:MM:SS".
func FormatHEven(hEven int) string {
	if hEven < 0 {
		return "00:00:00"
	}
	hour := hEven / 10000
	minute := (hEven / 100) % 100
	second := hEven % 100
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

// ParseTradeHistoryJSON parses the raw JSON body of a TSETMC GetTradeHistory response,
// extracts individual trade items, and computes the summary metrics excluding canceled trades.
func ParseTradeHistoryJSON(body []byte) (*dto.TradesResponseDto, error) {
	var payload tsetmcTradeHistoryPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse trade history json: %w", err)
	}

	rawTrades := payload.TradeHistory
	trades := make([]dto.TradeItemDto, len(rawTrades))

	var totalVolume int64
	var totalValue float64
	minPrice := math.MaxFloat64
	maxPrice := -math.MaxFloat64
	tradeCount := 0

	var earliestNTran int = math.MaxInt
	var latestNTran int = math.MinInt
	var firstPrice float64
	var lastPrice float64

	for i, t := range rawTrades {
		isCanceled := t.Canceled != 0

		trades[i] = dto.TradeItemDto{
			NTran:    t.NTran,
			Time:     FormatHEven(t.HEven),
			Volume:   t.QTitTran,
			Price:    t.PTran,
			Canceled: isCanceled,
		}

		// Exclude canceled trades from summary
		if !isCanceled {
			tradeCount++
			totalVolume += t.QTitTran
			totalValue += float64(t.QTitTran) * t.PTran

			if t.PTran < minPrice {
				minPrice = t.PTran
			}
			if t.PTran > maxPrice {
				maxPrice = t.PTran
			}

			// Chronological tracking via nTran (transaction number) or tie-breaker
			if t.NTran < earliestNTran {
				earliestNTran = t.NTran
				firstPrice = t.PTran
			}
			if t.NTran > latestNTran {
				latestNTran = t.NTran
				lastPrice = t.PTran
			}
		}
	}

	if tradeCount == 0 {
		minPrice = 0
		maxPrice = 0
		firstPrice = 0
		lastPrice = 0
	}

	summary := dto.TradeSummaryDto{
		TradeCount:  tradeCount,
		TotalVolume: totalVolume,
		TotalValue:  totalValue,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		FirstPrice:  firstPrice,
		LastPrice:   lastPrice,
	}

	return &dto.TradesResponseDto{
		Summary: summary,
		Trades:  trades,
	}, nil
}
