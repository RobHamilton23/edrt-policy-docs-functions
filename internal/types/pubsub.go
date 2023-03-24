package types

type PubSubMessage struct {
	Data       []byte                 `json:"data"`
	Attributes map[string]interface{} `json:"attributes"`
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PolicyDocsMessage struct {
	Site        string `json:"site"`
	Environment string `json:"env"`
	Hostname    string `json:"hostname"`
}
