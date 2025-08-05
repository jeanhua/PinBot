package messagechain

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
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
	httpUtil := utils.HttpUtil{}
	err := httpUtil.RequestWithNoResponse(http.MethodPost, mc.urlpath, httpUtil.WithJsonBody(mc))
	if err != nil {
		log.Println("error when send group chain message")
	}
}
