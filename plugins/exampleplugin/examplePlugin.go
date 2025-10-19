package exampleplugin

import (
	"log"

	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/model"
)

type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) OnFriendMsg(message *model.FriendMessage) bool {
	log.Printf("[私聊消息](%s %d):%s\n", message.Sender.Nickname, message.Sender.UserId, botcontext.ExtractPrivateMessageText(message))
	return true
}
func (p *Plugin) OnGroupMsg(message *model.GroupMessage) bool {
	text, _ := botcontext.ExtractGroupMessageContent(message)
	log.Printf("[群聊消息(%d)](%s):%s\n", message.GroupId, message.Sender.Nickname, text)
	return true
}
