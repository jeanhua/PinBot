package botcommand

import (
	"strings"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var Plugins []*PluginMeta

type PluginMeta struct {
	Name        string
	Description string
	IsPublic    bool
}

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messagechain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text("\n" + config.GetConfig().GetString("bot_config.help_words.group"))
		chain.Send()
		return false
	case "/plugin", "/plugins", "/插件":
		text := getPluginHelpText(false)
		chain := messagechain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text("\n" + text)
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
		chain.Text(config.GetConfig().GetString("bot_config.help_words.friend"))
		chain.Send()
		return false
	case "/plugin", "/plugins", "/插件":
		text := getPluginHelpText(false)
		chain := messagechain.Friend(msg.UserId)
		chain.Text(text)
		chain.Send()
		return false
	default:
		return true
	}
}

func getPluginHelpText(isAll bool) string {
	pluginLen := len(Plugins)
	if pluginLen != 0 {
		text := ""
		for index, p := range Plugins {
			if !isAll && !p.IsPublic {
				continue
			}
			if index != pluginLen-1 {
				text += p.Name + ":\n" + p.Description + "\n\n"
			} else {
				text += p.Name + ":\n" + p.Description
			}
		}
		return text
	}
	return "\n未加载任何插件"
}
