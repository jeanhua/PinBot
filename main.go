package main

import (
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/plugins/defaultplugin"
	"github.com/jeanhua/PinBot/plugins/exampleplugin"
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
	instance.Plugins.AddPlugin(exampleplugin.NewPlugin())

	/* -----------在这里注册插件----------- */
	//

	//
	/* -----------在上面注册插件----------- */

	// 系统默认插件，包含AI聊天
	instance.Plugins.AddPlugin(defaultplugin.NewPlugin())
}

/**------------------------------**/
