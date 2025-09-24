package plugins

import (
	"fmt"
	"sync"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/ai/functioncall"
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/datastructure/concurrent"
	"github.com/jeanhua/PinBot/datastructure/tuple"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var (
	llmLock    sync.Mutex
	currentRun = 0
	aiModelMap = concurrent.NewConcurrentMap[uint, aicommunicate.AiModel]()
	repeatMap  = concurrent.NewConcurrentMap[uint, tuple.Tuple[int, string]]()
)

var DefaultPlugin = botcontext.NewPluginContext("default plugin", defaultPluginOnFriend, defaultPluginOnGroup, "系统默认插件, AI智能体, 可以聊天，逛校园集市，检索和浏览网页, 群语音聊天, 发表情包, 搜索歌曲等")

func defaultPluginOnFriend(message *model.FriendMessage) bool {
	text := botcontext.ExtractPrivateMessageText(message)
	handlePrivateAIChat(message, text)
	return false
}

func defaultPluginOnGroup(message *model.GroupMessage) bool {
	text, mention := botcontext.ExtractGroupMessageContent(message)
	// 复读机
	repeat, ok := repeatMap.Get(message.GroupId)
	if ok {
		if repeat.First >= 2 && repeat.Second == text {
			msg := messagechain.Group(message.GroupId).Text(text)
			msg.Send()
			repeatMap.Set(message.GroupId, tuple.Of(-100, text))
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
	llmLock.Lock()
	if currentRun > config.GetConfig().MaxRun {
		llmLock.Unlock()
		sendBusyResponse(msg)
		return
	}
	currentRun += 1
	llmLock.Unlock()
	defer func() {
		llmLock.Lock()
		currentRun -= 1
		llmLock.Unlock()
	}()
	processGroupAIResponse(msg, text)
}

// 发送忙碌响应
func sendBusyResponse(msg *model.GroupMessage) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(msg.UserId)
	chain.Text(" 人太多了忙不过来了，待会再来问我哦!")
	chain.Send()
}

// 处理群AI响应
func processGroupAIResponse(msg *model.GroupMessage, text string) {
	deepseek := getOrCreateGroupAIModel(msg.GroupId)
	showName := msg.Sender.Card
	if showName == "" {
		showName = msg.Sender.Nickname
	}
	deepseek.Ask(fmt.Sprintf("[nickname: %s(%d)]: %s", showName, msg.Sender.UserId, text), msg, nil)
}

// 处理私聊AI聊天
func handlePrivateAIChat(msg *model.FriendMessage, text string) {
	llmLock.Lock()
	if currentRun > config.GetConfig().MaxRun {
		llmLock.Unlock()
		sendPrivateBusyResponse(msg.UserId)
		return
	}
	currentRun += 1
	llmLock.Unlock()
	defer func() {
		llmLock.Lock()
		currentRun--
		llmLock.Unlock()
	}()
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
	deepseek.Ask(fmt.Sprintf("[nickname: %s(%d)]: %s", msg.Sender.Nickname, msg.Sender.UserId, text), nil, msg)
}

// 获取或创建私聊AI模型
func getOrCreatePrivateAIModel(uid uint) aicommunicate.AiModel {
	deepseek, ok := aiModelMap.Get(uid)
	if !ok {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AiPrompt,
			config.GetConfig().AIToken,
			functioncall.TargetFriend,
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
			config.GetConfig().AIToken,
			functioncall.TargetGroup,
		)
		aiModelMap.Set(uid, deepseek)
	}
	return deepseek
}
