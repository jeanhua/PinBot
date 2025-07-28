package messagechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeanhua/PinBot/config"
)

type AIMessageData struct {
	urlpath   string
	GroupId   uint   `json:"group_id"`
	Character string `json:"character"`
	Text      string `json:"text"`
}

func AIMessage(groupUin uint, charactor string, text string) *AIMessageData {
	return &AIMessageData{
		urlpath:   config.GetConfig().NapCatServerUrl + "/send_group_ai_record",
		GroupId:   groupUin,
		Character: charactor,
		Text:      text,
	}
}

func (mc *AIMessageData) Send() {
	body, err := json.Marshal(mc)
	if err != nil {
		fmt.Println("error when json marshal: Send GroupForwardChain")
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, mc.urlpath, bytes.NewReader(body))
	if err != nil {
		fmt.Println("error when create http request: Send GroupForwardChain")
		return
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("error when send http request: Send GroupForwardChain")
		return
	}
	defer resp.Body.Close()
}
