package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	botcommand "github.com/jeanhua/PinBot/botCommand"
	"github.com/jeanhua/PinBot/config"
	messageChain "github.com/jeanhua/PinBot/messagechain"
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
		chain := messageChain.Friend(uid)
		chain.Text("正在思考中，不要着急哦")
		messageChain.SendMessage(chain)
		return
	}
	defer llmLock.Unlock()
	zhipu := zhipuMap[uint(msg.UserId)]
	if zhipu == nil {
		zhipu = aicommunicate.NewZhipu(config.ConfigInstance.ZhipuToken, config.ConfigInstance.AI_Prompt)
		zhipuMap[uint(msg.UserId)] = zhipu
	}
	reply := zhipu.Ask(text)

	if reply == nil {
		chain := messageChain.Friend(uid)
		chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
		messageChain.SendMessage(chain)
		return
	}
	rreply := []rune(reply.Response)
	replyLength := len(rreply)
	if replyLength <= 500 {
		chain := messageChain.Friend(uid)
		chain.Text(reply.Response)
		messageChain.SendMessage(chain)
	} else {
		for i := 0; i <= replyLength/500; i++ {
			chain := messageChain.Friend(uid)
			if (i+1)*500 < replyLength {
				chain.Text(string(rreply[i*500 : (i+1)*500]))
				messageChain.SendMessage(chain)
			} else if i*500 < replyLength {
				chain.Text(string(rreply[i*500:]))
				messageChain.SendMessage(chain)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}
