package logic

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

// 复读机
type repeatTuple struct {
	count int32
	text  string
}

// 处理群组消息
func onGroupMessage(msg model.GroupMessage) {
	text, mention := extractMessageContent(msg)

	if strings.TrimSpace(text) == "" {
		return
	}

	trimText := strings.TrimSpace(text)

	// 处理指令
	if mention {
		if botcommand.DealGroupCommand(trimText, &msg) {
			return
		}
	} else {
		// 非提及消息处理复读机功能
		handleRepeatFeature(msg, trimText)
		return
	}

	// AI聊天处理
	text = fmt.Sprintf("[%s]: %s", msg.Sender.Nickname, text)
	handleAIChat(msg, text)
}

// 从消息中提取文本内容和是否提及机器人
func extractMessageContent(msg model.GroupMessage) (string, bool) {
	text := ""
	mention := false

	for _, t := range msg.Message {
		switch t.Type {
		case "text":
			text += t.Data["text"].(string)
		case "at":
			if t.Data["qq"].(string) == strconv.Itoa(msg.SelfId) {
				mention = true
			}
		}
	}

	return text, mention
}

// 处理复读机功能
func handleRepeatFeature(msg model.GroupMessage, text string) {
	repeatLock.Lock()
	defer repeatLock.Unlock()

	if repeat.count >= 3 && repeat.text == text {
		chain := messagechain.Group(msg.GroupId)
		chain.Text(repeat.text)
		messagechain.SendMessage(chain)
		repeat.count = -100
	} else if repeat.text == text {
		repeat.count++
	} else {
		repeat.count = 1
		repeat.text = text
	}
}

// 处理AI聊天
func handleAIChat(msg model.GroupMessage, text string) {
	if !llmLock.TryLock() {
		sendBusyResponse(msg)
		return
	}
	defer llmLock.Unlock()
	processAIResponse(msg, text)
}

// 发送忙碌响应
func sendBusyResponse(msg model.GroupMessage) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(msg.UserId)
	chain.Text(" 正在思考中，不要着急哦")
	messagechain.SendMessage(chain)
}

// 处理AI响应
func processAIResponse(msg model.GroupMessage, text string) {
	uid := msg.UserId
	deepseek := getOrCreateAIModel(msg.GroupId)
	replies := deepseek.Ask(text)
	for _, reply := range replies {
		if strings.TrimSpace(reply.Response) == "" {
			continue
		}
		if reply == nil {
			sendErrorResponse(msg, uid)
			return
		}
		sendReply(msg, uid, reply.Response)
	}
}

// 获取或创建AI模型
func getOrCreateAIModel(groupId int) aicommunicate.AiModel {
	deepseek := aiModelMap[uint(groupId)]
	if deepseek == nil {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.ConfigInstance.AI_Prompt,
			config.ConfigInstance.SiliconflowToken,
			func(text string) {
				aimsg := messagechain.AIMessage(groupId, "lucy-voice-suxinjiejie", text)
				aimsg.Send()
			},
		)
		aiModelMap[uint(groupId)] = deepseek
	}
	return deepseek
}

// 发送错误响应
func sendErrorResponse(msg model.GroupMessage, uid int) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(int(uid))
	chain.Text(" 抱歉，我遇到了一些问题，请稍后再试。")
	messagechain.SendMessage(chain)
}

// 发送回复消息
func sendReply(msg model.GroupMessage, uid int, response string) {
	rreply := []rune(response)
	replyLength := len(rreply)

	if replyLength <= 500 {
		sendShortReply(msg, uid, response)
	} else {
		sendLongReply(msg, rreply, replyLength)
	}
}

// 发送短回复
func sendShortReply(msg model.GroupMessage, uid int, response string) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(int(uid))
	chain.Text(" " + response)
	messagechain.SendMessage(chain)
}

// 发送长回复
func sendLongReply(msg model.GroupMessage, rreply []rune, replyLength int) {
	forward := messagechain.GroupForward(msg.GroupId, "聊天记录")
	chain := messagechain.Group(msg.GroupId)
	chain.Mention(msg.UserId)
	messagechain.SendMessage(chain)

	for i := 0; i <= replyLength/500; i++ {
		start := i * 500
		end := (i + 1) * 500

		if end < replyLength {
			forward.Text(string(rreply[start:end]), msg.SelfId, "江颦")
		} else if start < replyLength {
			forward.Text(string(rreply[start:]), msg.SelfId, "江颦")
		}
	}

	time.Sleep(500 * time.Millisecond)
	forward.Send()
}
