package messagechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeanhua/PinBot/config"
)

type GroupForwardChain struct {
	urlpath  string
	userId   string
	nickname string
	GroupId  uint `json:"group_id"`
	Messages []struct {
		Type string `json:"type"`
		Data struct {
			UserId   string      `json:"user_id"`
			NickName string      `json:"nickname"`
			Content  MessageData `json:"content"`
		} `json:"data"`
	} `json:"messages"`
	News    []map[string]interface{} `json:"news"`
	Prompt  string                   `json:"prompt"`
	Summary string                   `json:"summary"`
	Source  string                   `json:"source"`
}

func GroupForward(groupUin uint, source string, userId string, nickname string) *GroupForwardChain {
	return &GroupForwardChain{
		urlpath:  config.GetConfig().NapCatServerUrl + "/send_group_forward_msg",
		userId:   userId,
		nickname: nickname,
		GroupId:  groupUin,
		Prompt:   "我喜欢你很久了，能不能做我女朋友",
		Summary:  "思考结果",
		Source:   source,
	}
}

func (mc *GroupForwardChain) Text(text string) MessageChain {
	mc.Messages = append(mc.Messages, struct {
		Type string "json:\"type\""
		Data struct {
			UserId   string      "json:\"user_id\""
			NickName string      "json:\"nickname\""
			Content  MessageData "json:\"content\""
		} "json:\"data\""
	}{
		Type: "node",
		Data: struct {
			UserId   string      "json:\"user_id\""
			NickName string      "json:\"nickname\""
			Content  MessageData "json:\"content\""
		}{
			UserId:   mc.userId,
			NickName: mc.nickname,
			Content: MessageData{
				Type: "text",
				Data: map[string]interface{}{
					"text": text,
				},
			},
		},
	})
	mc.News = append(mc.News, map[string]interface{}{
		"text": mc.nickname + ":" + text,
	})
	return mc
}

func (mc *GroupForwardChain) Send() {
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
