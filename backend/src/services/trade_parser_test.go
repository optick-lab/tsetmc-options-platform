package services

import (
	"testing"
)

// Fixture with canceled trade and custom values
const fixtureWithCanceled = `{
  "tradeHistory": [
    {"insCode":null,"dEven":0,"nTran":3,"hEven":122940,"qTitTran":10,"pTran":299.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},
    {"insCode":null,"dEven":0,"nTran":2,"hEven":94332,"qTitTran":100,"pTran":500.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":1},
    {"insCode":null,"dEven":0,"nTran":1,"hEven":94033,"qTitTran":381,"pTran":290.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0}
  ]
}`

// Actual 14-record body from DB
const fixtureActualStoredBody = `{"tradeHistory":[{"insCode":null,"dEven":0,"nTran":14,"hEven":122940,"qTitTran":10,"pTran":299.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":13,"hEven":121521,"qTitTran":6,"pTran":299.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":11,"hEven":115231,"qTitTran":10,"pTran":300.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":12,"hEven":115231,"qTitTran":4,"pTran":299.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":10,"hEven":114556,"qTitTran":300,"pTran":313.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":9,"hEven":114441,"qTitTran":200,"pTran":313.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":8,"hEven":114327,"qTitTran":100,"pTran":307.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":7,"hEven":113700,"qTitTran":20,"pTran":292.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":6,"hEven":101140,"qTitTran":89,"pTran":291.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":5,"hEven":101104,"qTitTran":10,"pTran":291.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":4,"hEven":95744,"qTitTran":1,"pTran":291.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":3,"hEven":94630,"qTitTran":19,"pTran":293.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":2,"hEven":94332,"qTitTran":101,"pTran":296.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0},{"insCode":null,"dEven":0,"nTran":1,"hEven":94033,"qTitTran":381,"pTran":290.00,"qTitNgJ":0,"iSensVarP":"\u0000","pPhSeaCotJ":0.0,"pPbSeaCotJ":0.0,"iAnuTran":0,"xqVarPJDrPRf":0.0,"canceled":0}]}`

func TestFormatHEven(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{122940, "12:29:40"},
		{94033, "09:40:33"},
		{95744, "09:57:44"},
		{0, "00:00:00"},
		{100, "00:01:00"},
	}

	for _, tc := range tests {
		actual := FormatHEven(tc.input)
		if actual != tc.expected {
			t.Errorf("FormatHEven(%d) = %s; want %s", tc.input, actual, tc.expected)
		}
	}
}

func TestParseTradeHistoryJSON_ExcludesCanceledFromSummary(t *testing.T) {
	result, err := ParseTradeHistoryJSON([]byte(fixtureWithCanceled))
	if err != nil {
		t.Fatalf("unexpected error parsing fixture: %v", err)
	}

	if len(result.Trades) != 3 {
		t.Fatalf("expected 3 trades in list, got %d", len(result.Trades))
	}

	// Canceled trade check in items
	if !result.Trades[1].Canceled {
		t.Errorf("expected trade index 1 (nTran=2) to have Canceled=true")
	}

	summary := result.Summary
	if summary.TradeCount != 2 {
		t.Errorf("expected TradeCount=2 (excluding canceled), got %d", summary.TradeCount)
	}

	expectedVolume := int64(391) // 381 + 10
	if summary.TotalVolume != expectedVolume {
		t.Errorf("expected TotalVolume=%d, got %d", expectedVolume, summary.TotalVolume)
	}

	expectedValue := 381.0*290.0 + 10.0*299.0 // 110490 + 2990 = 113480
	if summary.TotalValue != expectedValue {
		t.Errorf("expected TotalValue=%f, got %f", expectedValue, summary.TotalValue)
	}

	if summary.MinPrice != 290.0 {
		t.Errorf("expected MinPrice=290, got %f", summary.MinPrice)
	}
	if summary.MaxPrice != 299.0 {
		t.Errorf("expected MaxPrice=299 (ignoring canceled 500), got %f", summary.MaxPrice)
	}
	if summary.FirstPrice != 290.0 {
		t.Errorf("expected FirstPrice=290, got %f", summary.FirstPrice)
	}
	if summary.LastPrice != 299.0 {
		t.Errorf("expected LastPrice=299, got %f", summary.LastPrice)
	}
}

func TestParseTradeHistoryJSON_ActualStoredBody(t *testing.T) {
	result, err := ParseTradeHistoryJSON([]byte(fixtureActualStoredBody))
	if err != nil {
		t.Fatalf("unexpected error parsing actual fixture: %v", err)
	}

	if len(result.Trades) != 14 {
		t.Fatalf("expected 14 trades, got %d", len(result.Trades))
	}

	summary := result.Summary
	if summary.TradeCount != 14 {
		t.Errorf("expected trade_count=14, got %d", summary.TradeCount)
	}
	if summary.TotalVolume != 1251 {
		t.Errorf("expected total_volume=1251, got %d", summary.TotalVolume)
	}
	if summary.TotalValue != 377073.0 {
		t.Errorf("expected total_value=377073.0, got %f", summary.TotalValue)
	}
	if summary.MinPrice != 290.0 {
		t.Errorf("expected min_price=290.0, got %f", summary.MinPrice)
	}
	if summary.MaxPrice != 313.0 {
		t.Errorf("expected max_price=313.0, got %f", summary.MaxPrice)
	}
	if summary.FirstPrice != 290.0 {
		t.Errorf("expected first_price=290.0, got %f", summary.FirstPrice)
	}
	if summary.LastPrice != 299.0 {
		t.Errorf("expected last_price=299.0, got %f", summary.LastPrice)
	}
}
