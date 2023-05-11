package gemini

import (
	"queueco/config"
	"queueco/plog"
)

type ReadMessageCallback func(response *MarketDataResponse) (pErr *plog.AppError)

var Subscriptions SubscriptionsMap

func init() {
	Subscriptions = make(SubscriptionsMap)
	for _, symbol := range config.GeminiSymbols {
		g := MarketDataSubscription{symbol: symbol}
		g.Connect()
		Subscriptions[symbol] = &g
	}
}

type SubscriptionsMap map[string]*MarketDataSubscription

func (s *SubscriptionsMap) RegisterCallbackForSymbol(symbol string, callback ReadMessageCallback) (pErr *plog.AppError) {
	subscription, ok := (*s)[symbol]
	if !ok {
		pErr = &plog.AppError{Location: "gemini.subscribe.RegisterCallbackForSymbol", Code: plog.SymbolNotFound,
			Message: "could not register callback as symbol not found in subscriptionsMap", Data: plog.Fields{"symbol": symbol}}
		return
	}
	subscription.RegisterCallback(callback)
	return
}

func (s *SubscriptionsMap) RegisterCallbackForAllSymbols(callback ReadMessageCallback) (pErr *plog.AppError) {
	for symbol, _ := range *s {
		pErr = s.RegisterCallbackForSymbol(symbol, callback)
		if pErr != nil {
			return
		}
	}
	return
}
