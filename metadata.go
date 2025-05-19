package elation

type Metadata struct {
	Data          *map[string]string `json:"data"`
	ObjectID      *string            `json:"object_id"`
	ObjectWebLink *string            `json:"object_web_link"`
}
