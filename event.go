package elation

type Event[ResourceT any] struct {
	Data          ResourceT `json:"data"`
	Action        string    `json:"action"`
	EventID       int64     `json:"event_id"`
	ApplicationID string    `json:"application_id"`
	Resource      string    `json:"resource"`
}
