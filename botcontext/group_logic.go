package botcontext

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

// 复读机
type repeatTuple struct {
	count int32
	text  string
}

// 处理群组消息
func (bot *BotContext) onGroupMessage(msg *model.GroupMessage) {
	text, mention := utils.ExtractMessageContent(msg)

	if strings.TrimSpace(text) == "" {
		return
	}

	trimText := strings.TrimSpace(text)

	// 处理指令
	if mention {
		if botcommand.DealGroupCommand(trimText, msg) {
			return
		}
	} else {
		// 非提及消息处理复读机功能
		handleRepeatFeature(msg, trimText)
		return
	}

	bot.Plugins.ExcuteGroup(msg)

	// AI聊天处理
	text = fmt.Sprintf("[%s]: %s", msg.Sender.Nickname, text)
	handleAIChat(msg, text)
}

// 处理复读机功能
func handleRepeatFeature(msg *model.GroupMessage, text string) {
	repeatLock.Lock()
	defer repeatLock.Unlock()

	if repeat.count >= 3 && repeat.text == text {
		chain := messagechain.Group(msg.GroupId)
		chain.Text(repeat.text)
		chain.Send()
		repeat.count = -100
	} else if repeat.text == text {
		repeat.count++
	} else {
		repeat.count = 1
		repeat.text = text
	}
}

// 处理AI聊天
func handleAIChat(msg *model.GroupMessage, text string) {
	if !llmLock.TryLock() {
		sendBusyResponse(msg)
		return
	}
	defer llmLock.Unlock()
	processAIResponse(msg, text)
}

// 发送忙碌响应
func sendBusyResponse(msg *model.GroupMessage) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(msg.UserId)
	chain.Text(" 正在思考中，不要着急哦")
	chain.Send()
}

// 处理AI响应
func processAIResponse(msg *model.GroupMessage, text string) {
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
func getOrCreateAIModel(groupId uint) aicommunicate.AiModel {
	deepseek := aiModelMap[uint(groupId)]
	if deepseek == nil {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AI_Prompt,
			config.GetConfig().SiliconflowToken,
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
func sendErrorResponse(msg *model.GroupMessage, uid uint) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(uid)
	chain.Text(" 抱歉，我遇到了一些问题，请稍后再试。")
	chain.Send()
}

// 发送回复消息
func sendReply(msg *model.GroupMessage, uid uint, response string) {
	rreply := []rune(response)
	replyLength := len(rreply)

	if replyLength <= 500 {
		sendShortReply(msg, uid, response)
	} else {
		sendLongReply(msg, rreply, replyLength)
	}
}

// 发送短回复
func sendShortReply(msg *model.GroupMessage, uid uint, response string) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(uid)
	chain.Text(" " + response)
	chain.Send()
}

// 发送长回复
func sendLongReply(msg *model.GroupMessage, rreply []rune, replyLength int) {
	forward := messagechain.GroupForward(msg.GroupId, "聊天记录", fmt.Sprintf("%d", msg.SelfId), "江颦")
	chain := messagechain.Group(msg.GroupId)
	chain.Mention(msg.UserId)
	chain.Send()

	for i := 0; i <= replyLength/500; i++ {
		start := i * 500
		end := (i + 1) * 500

		if end < replyLength {
			forward.Text(string(rreply[start:end]))
		} else if start < replyLength {
			forward.Text(string(rreply[start:]))
		}
	}

	time.Sleep(500 * time.Millisecond)
	forward.Send()
}
