package logic

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
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
			chain := messagechain.Group(msg.GroupId)
			chain.Text(repeat.text)
			messagechain.SendMessage(chain)
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
		chain := messagechain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(" 正在思考中，不要着急哦")
		messagechain.SendMessage(chain)
		return
	}
	defer llmLock.Unlock()
	uid := msg.UserId
	deepseek := aiModelMap[uint(msg.GroupId)]
	if deepseek == nil {
		deepseek = aicommunicate.NewDeepSeekV3(config.ConfigInstance.AI_Prompt, config.ConfigInstance.SiliconflowToken, func(text string) {
			aimsg := messagechain.AIMessage(msg.GroupId, "lucy-voice-suxinjiejie", text)
			log.Println("发送语音")
			aimsg.Send()
		})
		aiModelMap[uint(msg.GroupId)] = deepseek
	}
	replys := deepseek.Ask(text)
	for _, reply := range replys {
		if strings.TrimSpace(reply.Response) == "" {
			continue
		}

		if reply == nil {
			chain := messagechain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(int(uid))
			chain.Text(" 抱歉，我遇到了一些问题，请稍后再试。")
			messagechain.SendMessage(chain)
			return
		}

		rreply := []rune(reply.Response)
		replyLength := len(rreply)

		if replyLength <= 500 {
			chain := messagechain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(int(uid))
			chain.Text(" " + reply.Response)
			messagechain.SendMessage(chain)
		} else {
			forward := messagechain.GroupForward(msg.GroupId, "聊天记录")
			chain := messagechain.Group(msg.GroupId)
			chain.Mention(msg.UserId)
			messagechain.SendMessage(chain)
			for i := 0; i <= replyLength/500; i++ {
				if (i+1)*500 < replyLength {
					forward.Text(string(rreply[i*500:(i+1)*500]), msg.SelfId, "江颦")
				} else if i*500 < replyLength {
					forward.Text(string(rreply[i*500:]), msg.SelfId, "江颦")
				}
			}
			time.Sleep(500 * time.Millisecond)
			forward.Send()
		}
	}
}
