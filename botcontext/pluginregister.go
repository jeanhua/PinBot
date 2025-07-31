package botcontext

import (
	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/model"
)

type BotPlugin struct {
	plugins []*PluginContext
}

type PluginContext struct {
	onGroupMsg  GroupPluginFunc
	onFriendMsg FriendPluginFunc
	name        string
	description string
	isPublic    bool
}

func NewPluginContext(name string, onFriend FriendPluginFunc, onGroup GroupPluginFunc, description string) *PluginContext {
	return &PluginContext{
		name:        name,
		onGroupMsg:  onGroup,
		onFriendMsg: onFriend,
		description: description,
		isPublic:    true,
	}
}

func (p *PluginContext) SetPrivate() *PluginContext {
	p.isPublic = false
	return p
}

type FriendPluginFunc func(message *model.FriendMessage) bool
type GroupPluginFunc func(message *model.GroupMessage) bool

func (plugin *BotPlugin) excuteFriend(message *model.FriendMessage) {
	for _, f := range plugin.plugins {
		runNext := f.onFriendMsg(message)
		if !runNext {
			break
		}
	}
}

func (plugin *BotPlugin) excuteGroup(message *model.GroupMessage) {
	for _, f := range plugin.plugins {
		runNext := f.onGroupMsg(message)
		if !runNext {
			break
		}
	}
}

func (p *BotPlugin) AddPlugin(plugin *PluginContext) {
	p.plugins = append(p.plugins, plugin)
	botcommand.Plugins = append(botcommand.Plugins, &botcommand.PluginMeta{
		Name:        plugin.name,
		Description: plugin.description,
		IsPublic:    plugin.isPublic,
	})
}
