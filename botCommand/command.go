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
		chain.Text(config.ConfigInstance.HelpWords.Group)
		messagechain.SendMessage(chain)
		return true
	}
	return false
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Friend(msg.UserId)
		chain.Text(config.ConfigInstance.HelpWords.Friend)
		messagechain.SendMessage(chain)
		return true
	default:
		return false
	}
}
