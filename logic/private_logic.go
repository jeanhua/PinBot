package logic

import (
	"fmt"
	"strings"
	"time"

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

		config.ConfigInstance_mu.RLock()
		reply, err := zhipu.RequestReply(uint(uid), text, config.ConfigInstance.AI_Prompt)
		config.ConfigInstance_mu.RUnlock()

		if err != nil {
			utils.LogErr(err.Error())
			chain := messageChain.Friend(uid)
			chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
			messageChain.SendMessage(chain)
			llmLock.Lock()
			ready = true
			llmLock.Unlock()
			return
		}
		rreply := []rune(reply)
		reply_length := len(rreply)
		if reply_length <= 500 {
			chain := messageChain.Friend(uid)
			chain.Text(reply)
			messageChain.SendMessage(chain)
		} else {

			for i := 0; i <= reply_length/500; i++ {
				chain := messageChain.Friend(uid)
				if (i+1)*500 < reply_length {
					chain.Text(string(rreply[i*500 : (i+1)*500]))
					messageChain.SendMessage(chain)
				} else if i*500 < reply_length {
					chain.Text(string(rreply[i*500:]))
					messageChain.SendMessage(chain)
				}
				time.Sleep(time.Millisecond * 500)
			}
		}

		llmLock.Lock()
		ready = true
		llmLock.Unlock()
	}(uid, text)
}
