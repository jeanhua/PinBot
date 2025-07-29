package messagechain

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jeanhua/PinBot/config"
)

type GroupChain struct {
	urlpath string
	Groupid string        `json:"group_id"`
	Message []MessageData `json:"message"`
}

func Group(groupUin uint) *GroupChain {
	return &GroupChain{
		urlpath: config.GetConfig().NapCatServerUrl + "/send_group_msg",
		Groupid: fmt.Sprintf("%d", groupUin),
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

func (mc *GroupChain) Reply(id uint) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	return mc
}

func (mc *GroupChain) Mention(userid uint) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "at",
		Data: map[string]interface{}{
			"qq": fmt.Sprintf("%d", userid),
		},
	})
	return mc
}

func (mc *GroupChain) UrlImage(url string) MessageChain {
	mc.Message = append(mc.Message, MessageData{
		Type: "image",
		Data: map[string]interface{}{
			"file":    url,
			"summary": "[图片]",
		},
	})
	return mc
}

func (mc *GroupChain) LocalImage(path string) MessageChain {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("error when open file: MessageChain: LocalImage", err)
		return mc
	}
	base64 := base64.StdEncoding.EncodeToString(file)
	mc.Message = append(mc.Message, MessageData{
		Type: "image",
		Data: map[string]interface{}{
			"file":    "base64://" + base64,
			"summary": "[图片]",
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
