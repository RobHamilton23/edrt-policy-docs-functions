package types

type PubSubMessage struct {
	Data       []byte                 `json:"data"`
	Attributes map[string]interface{} `json:"attributes"`
}

type MessagePublishedData struct {
	Message PubSubMessage
}
