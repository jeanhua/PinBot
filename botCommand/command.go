package botcommand

import (
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(" @我发送 清除记录 可以清除聊天记录哦")
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Friend(msg.UserId)
		chain.Text(" @我发送 清除记录 可以清除聊天记录哦")
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}
