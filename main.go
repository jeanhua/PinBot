package main

import (
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/plugins/defaultplugin"
	"github.com/jeanhua/PinBot/plugins/exampleplugin"
)

func main() {
	config.LoadConfig()
	bot := botcontext.NewBot()
	registerPlugin(bot)
	bot.Run()
}

/**------------------------------**/
/**
* 插件注册
**/
func registerPlugin(bot *botcontext.BotContext) {
	// 示例插件：打印消息
	bot.Plugins.AddPlugin(exampleplugin.NewPlugin(), "示例插件", "打印日志消息", false)

	/* -----------在这里注册插件----------- */
	//

	//
	/* -----------在上面注册插件----------- */

	// 系统默认插件，包含AI聊天
	bot.Plugins.AddPlugin(defaultplugin.NewPlugin(), "系统默认插件", "系统默认插件, AI智能体, 可以聊天，逛校园集市，检索和浏览网页, 群语音聊天, 发表情包, 搜索歌曲等", true)
}

/**------------------------------**/
