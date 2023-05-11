package orderbook

import (
	"github.com/stretchr/testify/assert"
	"queueco/clients/gemini"
	"queueco/gui"
	"queueco/plog"
	"testing"
)

var mockMDResp1 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4002.00",
	Remaining: "3.00",
	Side:      "ask",
	Type:      "change"}}}
var mockMDResp2 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4002.00",
	Remaining: "1.00",
	Side:      "ask",
	Type:      "change"}}}
var mockMDResp3 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4000.01",
	Remaining: "3.00",
	Side:      "bid",
	Type:      "change"}}}
var mockMDResp4 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4000.02",
	Remaining: "4.00",
	Side:      "bid",
	Type:      "change"}}}
var mockMDResp5 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4001.95",
	Remaining: "0.50",
	Side:      "ask",
	Type:      "change"}}}
var mockMDResp6 = gemini.MarketDataResponse{Events: []gemini.MarketDataEvent{{Delta: "",
	Price:     "4002.00",
	Remaining: "3.00",
	Side:      "ask",
	Type:      "trade"}}}

func TestOrderBook_getBestBidsAndAsks(t *testing.T) {
	tmpChan := make(chan gui.SymbolViewUpdate, 100)
	symbol := "BTCUSD"
	ob := NewOrderBook(symbol, tmpChan)
	var pErr *plog.AppError
	var b, bq, a, aq float64

	pErr = ob.UpdateOrderBook(&mockMDResp1)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{0, 0, 4002, 3}, []float64{b, bq, a, aq})

	pErr = ob.UpdateOrderBook(&mockMDResp2)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{0, 0, 4002, 1}, []float64{b, bq, a, aq})

	pErr = ob.UpdateOrderBook(&mockMDResp3)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{4000.01, 3, 4002, 1}, []float64{b, bq, a, aq})

	pErr = ob.UpdateOrderBook(&mockMDResp4)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{4000.02, 4, 4002, 1}, []float64{b, bq, a, aq})

	pErr = ob.UpdateOrderBook(&mockMDResp5)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{4000.02, 4, 4001.95, 0.5}, []float64{b, bq, a, aq})

	// no change on trade event
	pErr = ob.UpdateOrderBook(&mockMDResp6)
	assert.Nil(t, pErr)
	b, bq, a, aq = ob.getBestBidsAndAsks()
	assert.Equal(t, []float64{4000.02, 4, 4001.95, 0.5}, []float64{b, bq, a, aq})
}
