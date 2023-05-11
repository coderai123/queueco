package gemini

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"queueco/config"
	"queueco/plog"
	"time"
)

const reconnectTime = time.Second

type MarketDataSubscription struct {
	*websocket.Conn
	symbol string
}

func (g *MarketDataSubscription) Connect() {
	url := config.GeminiMarketDataUrl + g.symbol
	for {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			plog.LogError("gemini.startWebsocketConn", fmt.Sprintf("failed to dial %s, got %v", url, err),
				plog.Fields{"url": url, "symbol": g.symbol})
			time.Sleep(time.Second)
			continue
		}
		g.Conn = conn
		break
	}
}

func (g *MarketDataSubscription) RegisterCallback(callback ReadMessageCallback) {
	go func() {
		for {
			_, message, err := g.ReadMessage()
			if err != nil {
				plog.LogError("gemini.RegisterCallback",
					fmt.Sprintf("Error while reading message from subscription, got %v", err),
					plog.Fields{"symbol": g.symbol})
				plog.LogInfo("gemini.RegisterCallback",
					"Reconnecting and subscribing", plog.Fields{"symbol": g.symbol})
				g.Close()
				g.Connect()
				g.RegisterCallback(callback)
				return
			}
			var resp MarketDataResponse
			err = json.Unmarshal(message, &resp)
			if err != nil {
				panic(err)
			}
			pErr := callback(&resp)
			if pErr != nil {
				panic(pErr)
			}
		}
	}()
}
