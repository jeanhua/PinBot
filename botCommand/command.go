package botcommand

import (
	"strings"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var Plugins = []PluginMeta{}

type PluginMeta struct {
	Name        string
	Description string
}

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(config.GetConfig().HelpWords.Group)
		chain.Send()
		return false
	}
	return true
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Friend(msg.UserId)
		chain.Text(config.GetConfig().HelpWords.Friend)
		chain.Send()
		return false
	case "/plugin", "/plugins", "/插件":
		pluginLen := len(Plugins)
		if msg.UserId == config.GetConfig().Admin_id && pluginLen != 0 {
			chain := messagechain.Friend(msg.UserId)
			text := ""
			for index, p := range Plugins {
				if index != pluginLen-1 {
					text += p.Name + ":\n" + p.Description + "\n\n"
				} else {
					text += p.Name + ":\n" + p.Description
				}
			}
			chain.Text(text)
			chain.Send()
			return false
		}
		return true
	default:
		return true
	}
}
