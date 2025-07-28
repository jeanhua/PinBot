package main

import (
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/plugins"
)

func main() {
	bot := botcontext.NewBot()
	registerPlugin(bot)
	bot.Run()
}

/**------------------------------**/
/**
* 插件注册
**/
func registerPlugin(instance *botcontext.BotContext) {
	// 示例插件：打印消息
	instance.Plugins.AddFriendPlugin(plugins.ExampleFriendPlugin)
	instance.Plugins.AddGroupPlugin(plugins.ExampleGroupPlugin)
	// 系统默认插件，包含AI聊天
	instance.Plugins.AddFriendPlugin(plugins.DefaultFriendPlugin)
	instance.Plugins.AddGroupPlugin(plugins.DefaultGroupPlugin)

}

/**------------------------------**/
