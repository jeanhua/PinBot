package messagechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeanhua/PinBot/config"
)

type FriendChain struct {
	urlpath string
	Userid  string        `json:"user_id"`
	Message []MessageData `json:"message"`
}

func Friend(friendUin uint) *FriendChain {
	return &FriendChain{
		urlpath: config.GetConfig().NapCatServerUrl + "/send_private_msg",
		Userid:  fmt.Sprintf("%d", friendUin),
		Message: make([]MessageData, 0),
	}
}

func (mc *FriendChain) Text(text string) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	})
	return mc
}

func (mc *FriendChain) Reply(id int) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	return mc
}

func (mc *FriendChain) Send() {
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
