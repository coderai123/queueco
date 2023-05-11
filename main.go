package main

import (
	"os"
	"os/signal"
	"queueco/clients/gemini"
	"queueco/config"
	"queueco/gui"
	"queueco/orderbook"
	"queueco/plog"
	"time"
)

func NewOrderBooks(symbolViewUpdates chan<- gui.SymbolViewUpdate) map[string]*orderbook.OrderBook {
	orderbooks := make(map[string]*orderbook.OrderBook)
	for _, symbol := range config.GeminiSymbols {
		orderbooks[symbol] = orderbook.NewOrderBook(symbol, symbolViewUpdates)
	}
	return orderbooks
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	symbolViewUpdates := make(chan gui.SymbolViewUpdate, 500000)

	orderbooks := NewOrderBooks(symbolViewUpdates)
	for symbol, ob := range orderbooks {
		gemini.Subscriptions.RegisterCallbackForSymbol(symbol, ob.UpdateOrderBook)
	}

	if config.GUIEnabled {
		terminal := gui.NewGUI(config.GeminiSymbols)
		pErr := terminal.Start(symbolViewUpdates)
		if pErr != nil {
			panic(pErr)
		}
	} else {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				//plog.LogInfo("main", "ticker", plog.Fields{"Heartbeat": t.String()})
				continue
			case i := <-interrupt:
				plog.LogInfo("main", "terminating program", plog.Fields{"os signal": i.String()})
				return
			}
		}
	}

}
