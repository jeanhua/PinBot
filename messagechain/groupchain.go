package messagechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jeanhua/PinBot/config"
)

type GroupChain struct {
	urlpath string
	Groupid string        `json:"group_id"`
	Message []MessageData `json:"message"`
}

func Group(groupUin int) *GroupChain {
	return &GroupChain{
		urlpath: config.GetConfig().NapCatServerUrl + "/send_group_msg",
		Groupid: strconv.Itoa(groupUin),
		Message: make([]MessageData, 0),
	}
}

func (mc *GroupChain) Text(text string) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	})
	return mc
}

func (mc *GroupChain) Reply(id int) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	return mc
}

func (mc *GroupChain) Mention(userid int) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "at",
		Data: map[string]interface{}{
			"qq": strconv.Itoa(userid),
		},
	})
	return mc
}

func (mc *GroupChain) Send() {
	body, err := json.Marshal(mc)
	if err != nil {
		fmt.Println("error when json marshal: Send GroupChain")
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, mc.urlpath, bytes.NewReader(body))
	if err != nil {
		fmt.Println("error when create http request: Send GroupChain")
		return
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("error when send http request: Send GroupChain")
		return
	}
	defer resp.Body.Close()
}
