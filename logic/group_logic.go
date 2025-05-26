package logic

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	botcommand "github.com/jeanhua/PinBot/botCommand"
	"github.com/jeanhua/PinBot/config"
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

type repeaatTurple struct {
	count int32
	text  string
}

var repeatlock sync.Mutex
var repeat = &repeaatTurple{}

func onGroupMessage(msg model.GroupMessage) {
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
	utils.LogErr(fmt.Sprintf("[%s]:%s", msg.Sender.Nickname, text))

	trimText := strings.TrimSpace(text)

	if strings.TrimSpace(text) == "" {
		return
	}

	// 处理指令
	if mention {
		ret := botcommand.DealGroupCommand(trimText, &msg)
		if ret {
			return
		}
	}

	if !mention {

		// 特性

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

		config.ConfigInstance_mu.RLock()
		reply, err := zhipu.RequestReply(uint(groupId), text, config.ConfigInstance.AI_Prompt)
		config.ConfigInstance_mu.RUnlock()

		if err != nil {
			utils.LogErr(err.Error())
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

		rreply := []rune(reply)
		reply_length := len(rreply)

		botcommand.CommandMu.RLock()
		defer botcommand.CommandMu.RUnlock()

		if reply_length <= 450 && botcommand.EnableAIAudio {
			aimsg := messageChain.AIMessage(groupId, "lucy-voice-suxinjiejie", reply)
			aimsg.Send()
		} else if reply_length <= 500 {
			chain := messageChain.Group(groupId)
			chain.Reply(messageId)
			chain.Mention(int(uid))
			chain.Text(" " + reply)
			messageChain.SendMessage(chain)
		} else {
			forward := messageChain.GroupForward(msg.GroupId, "江颦的思考结果")
			chain := messageChain.Group(msg.GroupId)
			chain.Mention(msg.UserId)
			messageChain.SendMessage(chain)
			for i := 0; i <= reply_length/500; i++ {
				if (i+1)*500 < reply_length {
					forward.Text(string(rreply[i*500:(i+1)*500]), msg.SelfId, "江颦")
				} else if i*500 < reply_length {
					forward.Text(string(rreply[i*500:]), msg.SelfId, "江颦")
				}
			}
			time.Sleep(500 * time.Millisecond)
			forward.Send()
		}

		llmLock.Lock()
		ready = true
		llmLock.Unlock()
	}(uid, text, msg.GroupId, msg.MessageId)
}
