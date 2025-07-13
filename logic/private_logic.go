package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

func onPrivateMessage(msg model.FriendMessage) {
	text := ""
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		}
	}
	utils.LogErr(fmt.Sprintf("[%s]:%s", msg.Sender.Nickname, text))
	trimText := strings.TrimSpace(text)
	uid := msg.UserId
	if strings.TrimSpace(text) == "" {
		return
	}

	// 处理指令
	ret := botcommand.DealFriendCommand(trimText, &msg)
	if ret {
		return
	}

	if llmLock.TryLock() == false {
		chain := messagechain.Friend(uid)
		chain.Text("正在思考中，不要着急哦")
		messagechain.SendMessage(chain)
		return
	}
	defer llmLock.Unlock()
	deepseek := aiModelMap[uint(msg.UserId)]
	if deepseek == nil {
		deepseek = aicommunicate.NewDeepSeekV3(config.ConfigInstance.AI_Prompt, config.ConfigInstance.SiliconflowToken, func(text string) {
			chain := messagechain.Friend(uid)
			chain.Text(text)
			messagechain.SendMessage(chain)
		})
		aiModelMap[uint(msg.UserId)] = deepseek
	}
	reply := deepseek.Ask(text)

	if reply == nil {
		chain := messagechain.Friend(uid)
		chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
		messagechain.SendMessage(chain)
		return
	}
	rreply := []rune(reply.Response)
	replyLength := len(rreply)
	if replyLength <= 500 {
		chain := messagechain.Friend(uid)
		chain.Text(reply.Response)
		messagechain.SendMessage(chain)
	} else {
		for i := 0; i <= replyLength/500; i++ {
			chain := messagechain.Friend(uid)
			if (i+1)*500 < replyLength {
				chain.Text(string(rreply[i*500 : (i+1)*500]))
				messagechain.SendMessage(chain)
			} else if i*500 < replyLength {
				chain.Text(string(rreply[i*500:]))
				messagechain.SendMessage(chain)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}
