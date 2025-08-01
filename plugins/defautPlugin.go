package plugins

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/datastructure/concurrent"
	"github.com/jeanhua/PinBot/datastructure/tuple"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

var (
	llmLock    sync.Mutex
	aiModelMap = concurrent.NewConcurrentMap[uint, aicommunicate.AiModel]()
	repeatMap  = concurrent.NewConcurrentMap[uint, tuple.Tuple[int, string]]()
)

var DefaultPlugin = botcontext.NewPluginContext("default plugin", defaultPluginOnFriend, defaultPluginOnGroup, "系统默认插件, AI智能体, 可以聊天，逛校园集市，检索和浏览网页, 群语音聊天等")

func defaultPluginOnFriend(message *model.FriendMessage) bool {
	text := utils.ExtractPrivateMessageText(message)
	handlePrivateAIChat(message, text)
	return false
}

func defaultPluginOnGroup(message *model.GroupMessage) bool {
	text, mention := utils.ExtractMessageContent(message)
	// 复读机
	repeat, ok := repeatMap.Get(message.GroupId)
	if ok {
		if repeat.First >= 2 && repeat.Second == text {
			msg := messagechain.Group(message.GroupId).Text(text)
			msg.Send()
			repeatMap.Set(message.GroupId, tuple.Of(1, text))
			return false
		} else if repeat.Second != text {
			repeatMap.Set(message.GroupId, tuple.Of(1, text))
		} else {
			repeatMap.Set(message.GroupId, tuple.Of(repeat.First+1, text))
		}
	} else {
		repeatMap.Set(message.GroupId, tuple.Of(1, text))
	}
	// AI聊天
	if !mention {
		return true
	}
	handleGroupAIChat(message, text)
	return false
}

// 处理AI聊天
func handleGroupAIChat(msg *model.GroupMessage, text string) {
	if !llmLock.TryLock() {
		sendBusyResponse(msg)
		return
	}
	defer llmLock.Unlock()
	processGroupAIResponse(msg, text)
}

// 发送忙碌响应
func sendBusyResponse(msg *model.GroupMessage) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(msg.UserId)
	chain.Text(" 正在思考中，不要着急哦")
	chain.Send()
}

// 处理群AI响应
func processGroupAIResponse(msg *model.GroupMessage, text string) {
	uid := msg.UserId
	deepseek := getOrCreateGroupAIModel(msg.GroupId)
	replies := deepseek.Ask(fmt.Sprintf("[nickname: %s]: %s", msg.Sender.Nickname, text))
	if replies == nil {
		sendErrorResponse(msg, msg.Sender.UserId)
		return
	}
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
		utils.SendShortReply(msg, uid, response)
	} else {
		utils.SendLongReply(msg, rreply, replyLength)
	}
}

// 处理私聊AI聊天
func handlePrivateAIChat(msg *model.FriendMessage, text string) {
	if !llmLock.TryLock() {
		sendPrivateBusyResponse(msg.UserId)
		return
	}
	defer llmLock.Unlock()

	processPrivateAIResponse(msg, text)
}

// 发送忙碌响应
func sendPrivateBusyResponse(uid uint) {
	chain := messagechain.Friend(uid)
	chain.Text("正在思考中，不要着急哦")
	chain.Send()
}

// 处理私聊AI响应
func processPrivateAIResponse(msg *model.FriendMessage, text string) {
	uid := msg.UserId
	deepseek := getOrCreatePrivateAIModel(uid)
	replies := deepseek.Ask(text)
	if replies == nil {
		chain := messagechain.Friend(msg.Sender.UserId)
		chain.Text("遇到了一点小问题，请稍后再试")
		chain.Send()
		return
	}

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
func getOrCreatePrivateAIModel(uid uint) aicommunicate.AiModel {
	deepseek, ok := aiModelMap.Get(uid)
	if !ok {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AiPrompt,
			config.GetConfig().SiliconflowToken,
			func(text string) {
				sendPrivateMessage(uid, text)
			},
		)
		aiModelMap.Set(uid, deepseek)
	}
	return deepseek
}

// 获取或创建群AI模型
func getOrCreateGroupAIModel(uid uint) aicommunicate.AiModel {
	deepseek, ok := aiModelMap.Get(uid)
	if !ok {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AiPrompt,
			config.GetConfig().SiliconflowToken,
			func(text string) {
				chain := messagechain.AIMessage(uid, "lucy-voice-suxinjiejie", text)
				chain.Send()
			},
		)
		aiModelMap.Set(uid, deepseek)
	}
	return deepseek
}

// 发送私聊错误响应
func sendPrivateErrorResponse(uid uint) {
	chain := messagechain.Friend(uid)
	chain.Text("抱歉，我遇到了一些问题，请稍后再试。")
	chain.Send()
}

// 发送私聊回复消息
func sendPrivateReply(uid uint, response string) {
	rreply := []rune(response)
	replyLength := len(rreply)

	if replyLength <= 500 {
		sendPrivateMessage(uid, response)
	} else {
		sendLongPrivateMessage(uid, rreply, replyLength)
	}
}

// 发送短私聊消息
func sendPrivateMessage(uid uint, text string) {
	chain := messagechain.Friend(uid)
	chain.Text(text)
	chain.Send()
}

// 发送长私聊消息（分段）
func sendLongPrivateMessage(uid uint, rreply []rune, replyLength int) {
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
