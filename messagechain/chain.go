package messageChain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jeanhua/PinBot/model"
)

const (
	DEBUG      = false
	ServerHost = "http://localhost:7824"
)

type MessageChain interface {
	Text(text string)
	Reply(userid int)
	Mention(userid int)
	build() []byte
}

type FriendChain struct {
	Userid  string        `json:"user_id"`
	Message []MessageData `json:"message"`
}

type GroupChain struct {
	Groupid string        `json:"group_id"`
	Message []MessageData `json:"message"`
}

type MessageData struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type AIMessageData struct {
	GroupId   int    `json:"group_id"`
	Character string `json:"character"`
	Text      string `json:"text"`
}

type GroupForwardChain struct {
	GroupId  int `json:"group_id"`
	Messages []struct {
		Type string `json:"type"`
		Data struct {
			UserId   string      `json:"user_id"`
			NickName string      `json:"nickname"`
			Content  MessageData `json:"content"`
		} `json:"data"`
	} `json:"messages"`
	News    map[string]interface{} `json:"news"`
	Prompt  string                 `json:"prompt"`
	Summary string                 `json:"summary"`
	Source  string                 `json:"source"`
}

func Group(groupUin int) *GroupChain {
	return &GroupChain{
		Groupid: strconv.Itoa(groupUin),
		Message: make([]MessageData, 0),
	}
}

func Friend(friendUin int) *FriendChain {
	return &FriendChain{
		Userid:  strconv.Itoa(friendUin),
		Message: make([]MessageData, 0),
	}
}

func GroupForward(groupUin int, source string) *GroupForwardChain {
	return &GroupForwardChain{
		GroupId: groupUin,
		Prompt:  "我喜欢你很久了，能不能做我女朋友",
		Summary: "思考结果",
		News: map[string]interface{}{
			"text": "文本消息",
		},
		Source: source,
	}
}

func AIMessage(groupUin int, charactor string, text string) *AIMessageData {
	return &AIMessageData{
		GroupId:   groupUin,
		Character: charactor,
		Text:      text,
	}
}

func (mc *FriendChain) Text(text string) {
	mc.Message = append(mc.Message, MessageData{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	})
}
func (mc *GroupChain) Text(text string) {
	mc.Message = append(mc.Message, MessageData{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	})
}
func (mc *GroupForwardChain) Text(text string, userId int, nickname string) {
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
			UserId:   strconv.Itoa(userId),
			NickName: nickname,
			Content: MessageData{
				Type: "text",
				Data: map[string]interface{}{
					"text": text,
				},
			},
		},
	})
}

func (mc *FriendChain) Reply(id int) {
	mc.Message = append(mc.Message, MessageData{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	})
}
func (mc *GroupChain) Reply(id int) {
	mc.Message = append(mc.Message, MessageData{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	})
}
func (mc *FriendChain) Mention(userid int) {
	mc.Message = append(mc.Message, MessageData{
		Type: "at",
		Data: map[string]interface{}{
			"qq": strconv.Itoa(userid),
		},
	})
}
func (mc *GroupChain) Mention(userid int) {
	mc.Message = append(mc.Message, MessageData{
		Type: "at",
		Data: map[string]interface{}{
			"qq": strconv.Itoa(userid),
		},
	})
}
func (mc *FriendChain) build() []byte {
	result, err := json.Marshal(&mc)
	if err != nil {
		log.Println(err)
		return nil
	}
	if DEBUG {
		fmt.Println(string(result))
	}
	return result
}
func (mc *GroupChain) build() []byte {
	result, err := json.Marshal(&mc)
	if err != nil {
		log.Println(err)
		return nil
	}
	if DEBUG {
		fmt.Println(string(result))
	}
	return result
}

func SendMessage(chain MessageChain) (*model.Response, error) {
	data := chain.build()
	if data == nil {
		return nil, fmt.Errorf("failed to build message")
	}

	url := ServerHost + "/send_private_msg"
	if _, ok := chain.(*GroupChain); ok {
		url = ServerHost + "/send_group_msg"
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	back := &model.Response{}
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyData, back)
	if err != nil {
		return nil, err
	}
	return back, nil
}

func (msg *AIMessageData) Send() error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	url := ServerHost + "/send_group_ai_record"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (msg *GroupForwardChain) Send() error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	url := ServerHost + "/send_group_forward_msg"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
