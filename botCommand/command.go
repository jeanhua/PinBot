package botcommand

import (
	"strings"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(config.GetConfig().HelpWords.Group)
		chain.Send()
		return true
	}
	return false
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Friend(msg.UserId)
		chain.Text(config.GetConfig().HelpWords.Friend)
		chain.Send()
		return true
	default:
		return false
	}
}
