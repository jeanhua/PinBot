package botplugin

import (
	"github.com/jeanhua/PinBot/model"
)

type BotPlugin struct {
	onFriendMsg []IFriendPlugin
	onGroupMsg  []IGroupPlugin
}

type IFriendPlugin func(message *model.FriendMessage) bool
type IGroupPlugin func(message *model.GroupMessage) bool

func (plugin *BotPlugin) ExcuteFriend(message *model.FriendMessage) {
	for _, f := range plugin.onFriendMsg {
		runNext := f(message)
		if !runNext {
			break
		}
	}
}

func (plugin *BotPlugin) ExcuteGroup(message *model.GroupMessage) {
	for _, f := range plugin.onGroupMsg {
		runNext := f(message)
		if !runNext {
			break
		}
	}
}

func (plugin *BotPlugin) AddFriendPlugin(runner IFriendPlugin) {
	plugin.onFriendMsg = append(plugin.onFriendMsg, runner)
}

func (plugin *BotPlugin) AddGroupPlugin(runner IGroupPlugin) {
	plugin.onGroupMsg = append(plugin.onGroupMsg, runner)
}
