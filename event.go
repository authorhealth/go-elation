package elation

import "encoding/json"

type Event struct {
	Data          json.RawMessage `json:"data"`
	Action        string          `json:"action"`
	EventID       int64           `json:"event_id"`
	ApplicationID int64           `json:"application_id"`
	Resource      string          `json:"resource"`
}
