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
	/* -----------在这里注册插件----------- */
	//
	instance.Plugins.AddPlugin(plugins.DailyHotPlugin)
	//
	/* -----------在上面注册插件----------- */

	// 示例插件：打印消息
	instance.Plugins.AddPlugin(plugins.ExamplePlugin.SetPrivate())
	// 系统默认插件，包含AI聊天
	instance.Plugins.AddPlugin(plugins.DefaultPlugin)
}

/**------------------------------**/
