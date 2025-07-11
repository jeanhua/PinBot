package logic

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
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

	if llmLock.TryLock() == false {
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(" 正在思考中，不要着急哦")
		messageChain.SendMessage(chain)
		return
	}
	defer llmLock.Unlock()
	uid := msg.UserId
	zhipu := zhipuMap[uint(msg.GroupId)]
	if zhipu == nil {
		zhipu = aicommunicate.NewZhipu(config.ConfigInstance.ZhipuToken, config.ConfigInstance.AI_Prompt)
		zhipuMap[uint(msg.GroupId)] = zhipu
	}
	reply := zhipu.Ask(text)

	if reply == nil {
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(int(uid))
		chain.Text(" 抱歉，我遇到了一些问题，请稍后再试。")
		messageChain.SendMessage(chain)
		return
	}

	rreply := []rune(reply.Response)
	reply_length := len(rreply)

	botcommand.CommandMu.RLock()
	defer botcommand.CommandMu.RUnlock()

	if reply_length <= 450 && botcommand.EnableAIAudio {
		aimsg := messageChain.AIMessage(msg.GroupId, "lucy-voice-suxinjiejie", reply.Response)
		aimsg.Send()
	} else if reply_length <= 500 {
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(int(uid))
		chain.Text(" " + reply.Response)
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
}
