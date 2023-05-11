package orderbook

import (
	"fmt"
	"github.com/emirpasic/gods/utils"
	"queueco/clients/gemini"
	"queueco/config"
	"queueco/gui"
	"queueco/plog"
	"strconv"
)

type Side int

const (
	Buy Side = iota
	Sell
)

type Order struct {
	orderId   int
	price     float64
	qty       float64
	qtyFilled float64
	side      Side
}

type OrderBook struct {
	symbol        string
	bestBid       float64
	bestBidQty    float64
	bestAsk       float64
	bestAskQty    float64
	Bids          *Heap
	Asks          *Heap
	guiUpdateChan chan<- gui.SymbolViewUpdate
}

func ReverseFloat64Comparator(a, b interface{}) int {
	return -utils.Float64Comparator(a, b)
}

func NewOrderBook(symbol string, guiUpdateChan chan<- gui.SymbolViewUpdate) *OrderBook {
	return &OrderBook{symbol: symbol, Bids: NewHeap(), Asks: NewHeap(), guiUpdateChan: guiUpdateChan}
}

func (ob *OrderBook) UpdateOrderBook(md *gemini.MarketDataResponse) (pErr *plog.AppError) {
	for _, event := range md.Events {
		if event.Type == "change" {
			price, err := strconv.ParseFloat(event.Price, 64)
			if err != nil {
				plog.LogError("orderbook.UpdateOrderBook", fmt.Sprintf("Error parsing price %s, got %v",
					event.Price, err), plog.Fields{"event": event})
				continue
			}

			size, err := strconv.ParseFloat(event.Remaining, 64)
			if err != nil {
				plog.LogError("orderbook.UpdateOrderBook", fmt.Sprintf("Error parsing Remaining %s, got %v",
					event.Remaining, err), plog.Fields{"event": event})
				continue
			}

			if event.Side == "bid" {
				if size == 0 {
					ob.Bids.Delete(-price)
				} else {
					ob.Bids.CreateOrUpdate(-price, size)
				}
			} else if event.Side == "ask" {
				if size == 0 {
					ob.Asks.Delete(price)
				} else {
					ob.Asks.CreateOrUpdate(price, size)
				}
			}

			bestBid, bidQty, bestAsk, askQty := ob.getBestBidsAndAsks()
			symbolViewUpdate := gui.GenerateSymbolViewUpdate(ob.symbol, ob.bestBid, ob.bestBidQty, ob.bestAsk,
				ob.bestAskQty, bestBid, bidQty, bestAsk, askQty)
			ob.bestBid = bestBid
			ob.bestBidQty = bidQty
			ob.bestAsk = bestAsk
			ob.bestAskQty = askQty

			if symbolViewUpdate.HasChanged() {
				if config.GUIEnabled {
					ob.UpdateGUI(symbolViewUpdate)
				} else {
					fmt.Println(fmt.Sprintf("%.2f %f  -  %.2f %f", bestBid, bidQty, bestAsk, askQty))
				}
			}

		} else if event.Type == "trade" {

		}
	}
	return
}

func (ob *OrderBook) UpdateGUI(s gui.SymbolViewUpdate) {
	ob.guiUpdateChan <- s
}

func (ob *OrderBook) getBestBidsAndAsks() (float64, float64, float64, float64) {
	var node *HeapNode
	var ok bool
	var bestBid, bidQty, bestAsk, askQty float64

	node, ok = ob.Bids.Peek()
	if ok {
		bestBid = -node.key
		bidQty = node.value.(float64)
	}

	node, ok = ob.Asks.Peek()
	if ok {
		bestAsk = node.key
		askQty = node.value.(float64)
	}

	return bestBid, bidQty, bestAsk, askQty
}
