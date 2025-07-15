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

// 处理私聊消息
func onPrivateMessage(msg model.FriendMessage) {
	text := extractPrivateMessageText(msg)
	if strings.TrimSpace(text) == "" {
		return
	}

	utils.LogErr(fmt.Sprintf("[%s]:%s", msg.Sender.Nickname, text))
	trimText := strings.TrimSpace(text)

	// 处理指令
	if botcommand.DealFriendCommand(trimText, &msg) {
		return
	}

	// 处理AI聊天
	handlePrivateAIChat(msg, text)
}

// 从私聊消息中提取文本内容
func extractPrivateMessageText(msg model.FriendMessage) string {
	text := ""
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		}
	}
	return text
}

// 处理私聊AI聊天
func handlePrivateAIChat(msg model.FriendMessage, text string) {
	if !llmLock.TryLock() {
		sendPrivateBusyResponse(msg.UserId)
		return
	}
	defer llmLock.Unlock()

	processPrivateAIResponse(msg, text)
}

// 发送忙碌响应
func sendPrivateBusyResponse(uid int) {
	chain := messagechain.Friend(uid)
	chain.Text("正在思考中，不要着急哦")
	messagechain.SendMessage(chain)
}

// 处理私聊AI响应
func processPrivateAIResponse(msg model.FriendMessage, text string) {
	uid := msg.UserId
	deepseek := getOrCreatePrivateAIModel(uid)
	replies := deepseek.Ask(text)

	for _, reply := range replies {
		if strings.TrimSpace(reply.Response) == "" {
			continue
		}

		if reply == nil {
			sendPrivateErrorResponse(uid)
			return
		}

		sendPrivateReply(uid, reply.Response)
	}
}

// 获取或创建私聊AI模型
func getOrCreatePrivateAIModel(uid int) aicommunicate.AiModel {
	deepseek := aiModelMap[uint(uid)]
	if deepseek == nil {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.ConfigInstance.AI_Prompt,
			config.ConfigInstance.SiliconflowToken,
			func(text string) {
				sendPrivateMessage(uid, text)
			},
		)
		aiModelMap[uint(uid)] = deepseek
	}
	return deepseek
}

// 发送私聊错误响应
func sendPrivateErrorResponse(uid int) {
	chain := messagechain.Friend(uid)
	chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
	messagechain.SendMessage(chain)
}

// 发送私聊回复消息
func sendPrivateReply(uid int, response string) {
	rreply := []rune(response)
	replyLength := len(rreply)

	if replyLength <= 500 {
		sendPrivateMessage(uid, response)
	} else {
		sendLongPrivateMessage(uid, rreply, replyLength)
	}
}

// 发送短私聊消息
func sendPrivateMessage(uid int, text string) {
	chain := messagechain.Friend(uid)
	chain.Text(text)
	messagechain.SendMessage(chain)
}

// 发送长私聊消息（分段）
func sendLongPrivateMessage(uid int, rreply []rune, replyLength int) {
	for i := 0; i <= replyLength/500; i++ {
		start := i * 500
		end := (i + 1) * 500

		var segment string
		if end < replyLength {
			segment = string(rreply[start:end])
		} else if start < replyLength {
			segment = string(rreply[start:])
		} else {
			continue
		}

		sendPrivateMessage(uid, segment)
		time.Sleep(500 * time.Millisecond)
	}
}
