package plugin

import (
	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/model"
)

type BotPlugin struct {
	plugins []PluginContext
}

type PluginContext interface {
	OnFriendMsg(message *model.FriendMessage) bool
	OnGroupMsg(message *model.GroupMessage) bool
}

func (p *BotPlugin) ExecuteFriend(message *model.FriendMessage) {
	for _, f := range p.plugins {
		runNext := f.OnFriendMsg(message)
		if !runNext {
			break
		}
	}
}

func (p *BotPlugin) ExecuteGroup(message *model.GroupMessage) {
	for _, f := range p.plugins {
		runNext := f.OnGroupMsg(message)
		if !runNext {
			break
		}
	}
}

func (p *BotPlugin) AddPlugin(plugin PluginContext, name, description string, isPublic bool) {
	p.plugins = append(p.plugins, plugin)
	botcommand.Plugins = append(botcommand.Plugins, &botcommand.PluginMeta{
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
	})
}
