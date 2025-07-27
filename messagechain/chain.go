package messagechain

const (
	DEBUG      = false
	ServerHost = "http://localhost:7824"
)

type MessageChain interface {
	Send()
}

type MessageData struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
