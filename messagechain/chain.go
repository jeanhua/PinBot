package messagechain

type MessageChain interface {
	Send()
}

type MessageData struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
