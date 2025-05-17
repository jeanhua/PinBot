package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	llm "github.com/jeanhua/PinBot/LLM"
	botcommand "github.com/jeanhua/PinBot/botCommand"
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var DEBUG bool
var zhipu *llm.ZhiPu

var llmLock sync.RWMutex
var ready bool

func main() {
	DEBUG = false
	ready = true
	zhipu = llm.NewZhiPu()
	http.HandleFunc("/Pinbot", Handler)
	log.Println("Server starting on http://localhost:7823...")
	log.Fatal(http.ListenAndServe(":7823", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error: ", err)
	}
	handleMessage(body)
	r.Body.Close()
}

func handleMessage(message []byte) {
	basic := &model.Message{}
	err := json.Unmarshal(message, &basic)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if basic.PostType != "message" {
		return
	}
	friendmsg := model.FriendMessage{}
	err = json.Unmarshal(message, &friendmsg)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if friendmsg.MessageType == "private" {
		handelPrivate(friendmsg)
	} else if friendmsg.MessageType == "group" {
		groupmsg := model.GroupMessage{}
		err := json.Unmarshal(message, &groupmsg)
		if err != nil {
			log.Println("error:", err)
			return
		}
		handleGroup(groupmsg)
	}
}

func handelPrivate(msg model.FriendMessage) {
	text := ""
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		}
	}
	if DEBUG {
		log.Println(text)
	}
	trimText := strings.TrimSpace(text)
	uid := msg.UserId
	if strings.TrimSpace(text) == "清除记录" {
		zhipu.Clear(uint(uid))
		chain := messageChain.Friend(uid)
		chain.Text("清除成功")
		messageChain.SendMessage(chain)
		return
	} else if strings.TrimSpace(text) == "" {
		return
	}

	// 处理指令
	ret := botcommand.DealFriendCommand(trimText, &msg)
	if ret {
		return
	}

	llmLock.RLock()
	if ready == false {
		chain := messageChain.Friend(uid)
		chain.Text("正在思考中，不要着急哦")
		messageChain.SendMessage(chain)
		llmLock.RUnlock()
		return
	}
	llmLock.RUnlock()

	go func(uid int, text string) {
		llmLock.Lock()
		if ready == false {
			chain := messageChain.Friend(uid)
			chain.Text("正在思考中，不要着急哦")
			messageChain.SendMessage(chain)
			llmLock.Unlock()
			return
		}
		ready = false
		llmLock.Unlock()
		reply, err := zhipu.RequestReply(uint(uid), text)
		if err != nil {
			log.Println("zhipu error: ", err)
			chain := messageChain.Friend(uid)
			chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
			messageChain.SendMessage(chain)
			llmLock.Lock()
			ready = true
			llmLock.Unlock()
			return
		}

		chain := messageChain.Friend(uid)
		chain.Text(strings.TrimSpace(reply))
		messageChain.SendMessage(chain)
		llmLock.Lock()
		ready = true
		llmLock.Unlock()
	}(uid, text)
}

type repeaatTurple struct {
	count int32
	text  string
}

var repeatlock sync.Mutex
var repeat = &repeaatTurple{}

func handleGroup(msg model.GroupMessage) {
	text := ""
	mention := false
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		} else if t.Type == "at" {
			if t.Data["qq"].(string) == strconv.Itoa(msg.SelfId) {
				mention = true
			}
		}
	}
	if DEBUG {
		log.Println(text)
	}

	trimText := strings.TrimSpace(text)

	if strings.TrimSpace(text) == "" {
		return
	}

	// 处理指令
	ret := botcommand.DealGroupCommand(trimText, &msg)
	if ret {
		return
	}

	if !mention {

		// 特性
		if trimText == "?" || trimText == "？" {
			chain := messageChain.Group(msg.GroupId)
			chain.Text("¿")
			messageChain.SendMessage(chain)
			return
		} else if strings.Contains(trimText, "我是") {
			chain := messageChain.Group(msg.GroupId)
			chain.Text("你是?")
			messageChain.SendMessage(chain)
			return
		}

		// 复读机
		repeatlock.Lock()
		if repeat.count >= 3 && repeat.text == strings.TrimSpace(text) {
			chain := messageChain.Group(msg.GroupId)
			chain.Text(repeat.text)
			messageChain.SendMessage(chain)
			repeat.count = -100
		} else if repeat.text == strings.TrimSpace(text) {
			repeat.count += 1
		} else {
			repeat.count = 1
			repeat.text = strings.TrimSpace(text)
		}
		repeatlock.Unlock()
		return
	}

	llmLock.RLock()
	if ready == false {
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(" 正在思考中，不要着急哦")
		messageChain.SendMessage(chain)
		llmLock.RUnlock()
		return
	}
	llmLock.RUnlock()
	uid := msg.UserId

	if strings.TrimSpace(text) == "清除记录" {
		zhipu.Clear(uint(msg.GroupId))
		chain := messageChain.Group(msg.GroupId)
		chain.Mention(int(uid))
		chain.Text(" 清除成功")
		messageChain.SendMessage(chain)
		return
	}

	go func(uid int, text string, groupId int, messageId int) {
		llmLock.Lock()
		if ready == false {
			chain := messageChain.Group(groupId)
			chain.Reply(messageId)
			chain.Mention(uid)
			chain.Text(" 正在思考中，不要着急哦")
			messageChain.SendMessage(chain)
			llmLock.Unlock()
			return
		}
		ready = false
		llmLock.Unlock()
		reply, err := zhipu.RequestReply(uint(groupId), text)
		if err != nil {
			log.Println("zhipu error: ", err)
			chain := messageChain.Group(groupId)
			chain.Reply(messageId)
			chain.Mention(int(uid))
			chain.Text(" 抱歉，我遇到了一些问题，请稍后再试。")
			messageChain.SendMessage(chain)
			llmLock.Lock()
			ready = true
			llmLock.Unlock()
			return
		}

		if len(reply) >= 450 {
			chain := messageChain.Group(groupId)
			chain.Reply(messageId)
			chain.Mention(int(uid))
			chain.Text(" " + strings.TrimSpace(reply))
			messageChain.SendMessage(chain)
		} else {
			aimsg := messageChain.AIMessage(groupId, "lucy-voice-suxinjiejie", reply)
			aimsg.Send()
		}

		llmLock.Lock()
		ready = true
		llmLock.Unlock()
	}(uid, text, msg.GroupId, msg.MessageId)
}
