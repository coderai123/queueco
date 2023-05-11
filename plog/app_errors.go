package plog

import (
	"encoding/json"
	"fmt"
)

// AppError json mapping are added so that error data can be
// sent as json to external services such as elastic-search, sentry,
// newrelic. This leads to improved usage of functionalities provided by
// external services such as indexing, alerting on certain fields.
type AppError struct {
	Location string                 `json:"location"`
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Data     map[string]interface{} `json:"data"`
}

func (a *AppError) Error() string {
	return a.String()
}

func (a *AppError) String() string {
	dataBytes, _ := json.Marshal(a.Data)
	return fmt.Sprintf("{Location: %s, Code: %d, Message: %s, Data: %s}", a.Location, a.Code, a.Message, string(dataBytes))
}
