package gui

type ViewUpdate struct {
	Value      float64
	HasChanged bool
}

type SymbolViewUpdate struct {
	Symbol  string
	BestBid ViewUpdate
	BidQty  ViewUpdate
	BestAsk ViewUpdate
	AskQty  ViewUpdate
}

func GenerateSymbolViewUpdate(symbol string, oldBid, oldBidQty, oldAsk, oldAskQty, newBid, newBidQty, newAsk, newAskQty float64) SymbolViewUpdate {
	s := SymbolViewUpdate{Symbol: symbol, BestBid: ViewUpdate{}, BidQty: ViewUpdate{}, BestAsk: ViewUpdate{}, AskQty: ViewUpdate{}}
	if oldBid != newBid {
		s.BestBid.HasChanged = true
		s.BestBid.Value = newBid
	}

	if oldBidQty != newBidQty {
		s.BidQty.HasChanged = true
		s.BidQty.Value = newBidQty
	}

	if oldAsk != newAsk {
		s.BestAsk.HasChanged = true
		s.BestAsk.Value = newAsk
	}

	if oldAskQty != newAskQty {
		s.AskQty.HasChanged = true
		s.AskQty.Value = newAskQty
	}
	return s
}

func (s *SymbolViewUpdate) HasChanged() bool {
	return s.BestBid.HasChanged || s.BidQty.HasChanged || s.BestAsk.HasChanged || s.AskQty.HasChanged
}
