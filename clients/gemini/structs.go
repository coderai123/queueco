package gemini

type MarketDataEvent struct {
	Delta     string `json:"delta,omitempty"`
	Price     string `json:"price,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Remaining string `json:"remaining,omitempty"`
	Side      string `json:"side,omitempty"`
	Type      string `json:"type,omitempty"`
	Amount    string `json:"amount,omitempty"`
	MakerSide string `json:"makerSide,omitempty"`
	TID       int64  `json:"tid,omitempty"`
}

type MarketDataResponse struct {
	EventID     int64             `json:"eventId"`
	Events      []MarketDataEvent `json:"events"`
	Timestamp   int64             `json:"timestamp"`
	TimestampMs int64             `json:"timestampms"`
	Type        string            `json:"type"`
}
