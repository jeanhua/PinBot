package defaultplugin

import (
	"fmt"
	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/ai/functioncall"
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/datastructure/concurrent"
	"github.com/jeanhua/PinBot/datastructure/tuple"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"log"
)

type Plugin struct {
	currentRun chan struct{}
	aiModelMap *concurrent.ConcurrentMap[uint, aicommunicate.AiModel]
	repeatMap  *concurrent.ConcurrentMap[uint, tuple.Tuple[int, string]]
}

func NewPlugin() *Plugin {
	return &Plugin{
		currentRun: make(chan struct{}, config.GetConfig().MaxRun),
		aiModelMap: concurrent.NewConcurrentMap[uint, aicommunicate.AiModel](),
		repeatMap:  concurrent.NewConcurrentMap[uint, tuple.Tuple[int, string]](),
	}
}

func (p *Plugin) OnFriendMsg(message *model.FriendMessage) bool {
	text := botcontext.ExtractPrivateMessageText(message)
	p.handlePrivateAIChat(message, text)
	return false
}

func (p *Plugin) OnGroupMsg(message *model.GroupMessage) bool {
	text, mention := botcontext.ExtractGroupMessageContent(message)
	rawMsg, _ := botcontext.ExtractGroupRawMessage(message)
	// 复读机
	repeat, ok := p.repeatMap.Get(message.GroupId)
	if ok {
		if repeat.First >= 2 && repeat.Second == rawMsg {
			msg := messagechain.Group(message.GroupId).Text(rawMsg)
			p.repeatMap.Set(message.GroupId, tuple.Of(-100, rawMsg))
			msg.Send()
			return false
		} else if repeat.Second != rawMsg {
			p.repeatMap.Set(message.GroupId, tuple.Of(1, rawMsg))
		} else {
			p.repeatMap.Set(message.GroupId, tuple.Of(repeat.First+1, rawMsg))
		}
	} else {
		p.repeatMap.Set(message.GroupId, tuple.Of(1, rawMsg))
	}
	// AI聊天
	if !mention {
		return true
	}
	p.handleGroupAIChat(message, text)
	return false
}

// 处理群AI聊天
func (p *Plugin) handleGroupAIChat(msg *model.GroupMessage, text string) {
	select {
	case p.currentRun <- struct{}{}:
		log.Println("new group handler run...")
		p.processGroupAIResponse(msg, text)
		<-p.currentRun
		log.Println("finish a group handler run...")
	default:
		sendGroupBusyResponse(msg)
	}

}

// 发送忙碌响应
func sendGroupBusyResponse(msg *model.GroupMessage) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(msg.UserId)
	chain.Text(" 人太多了忙不过来了，待会再来问我哦!")
	chain.Send()
}

// 处理群AI响应
func (p *Plugin) processGroupAIResponse(msg *model.GroupMessage, text string) {
	deepseek := p.getOrCreateGroupAIModel(msg.GroupId)
	showName := msg.Sender.Card
	if showName == "" {
		showName = msg.Sender.Nickname
	}
	deepseek.Ask(fmt.Sprintf("[nickname: %s(%d)]: %s", showName, msg.Sender.UserId, text), msg, nil)
}

// 处理私聊AI聊天
func (p *Plugin) handlePrivateAIChat(msg *model.FriendMessage, text string) {
	select {
	case p.currentRun <- struct{}{}:
		log.Println("new private handler run...")
		p.processPrivateAIResponse(msg, text)
		<-p.currentRun
		log.Println("finish a private handler run...")
	default:
		sendPrivateBusyResponse(msg.Sender.UserId)
	}

}

// 发送忙碌响应
func sendPrivateBusyResponse(uid uint) {
	chain := messagechain.Friend(uid)
	chain.Text("正在思考中，不要着急哦")
	chain.Send()
}

// 处理私聊AI响应
func (p *Plugin) processPrivateAIResponse(msg *model.FriendMessage, text string) {
	uid := msg.UserId
	deepseek := p.getOrCreatePrivateAIModel(uid)
	deepseek.Ask(fmt.Sprintf("[nickname: %s(%d)]: %s", msg.Sender.Nickname, msg.Sender.UserId, text), nil, msg)
}

// 获取或创建私聊AI模型
func (p *Plugin) getOrCreatePrivateAIModel(uid uint) aicommunicate.AiModel {
	deepseek, ok := p.aiModelMap.Get(uid)
	if !ok {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AiPrompt,
			config.GetConfig().AIToken,
			functioncall.TargetFriend,
		)
		p.aiModelMap.Set(uid, deepseek)
	}
	return deepseek
}

// 获取或创建群AI模型
func (p *Plugin) getOrCreateGroupAIModel(uid uint) aicommunicate.AiModel {
	deepseek, ok := p.aiModelMap.Get(uid)
	if !ok {
		deepseek = aicommunicate.NewDeepSeekV3(
			config.GetConfig().AiPrompt,
			config.GetConfig().AIToken,
			functioncall.TargetGroup,
		)
		p.aiModelMap.Set(uid, deepseek)
	}
	return deepseek
}
