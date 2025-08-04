package botcontext

import (
	"strings"

	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/model"
)

// 处理私聊消息
func (bot *BotContext) onPrivateMessage(msg *model.FriendMessage) {
	text := ExtractPrivateMessageText(msg)
	if strings.TrimSpace(text) == "" {
		return
	}
	trimText := strings.TrimSpace(text)

	// 处理指令
	if !botcommand.DealFriendCommand(trimText, msg) {
		return
	}

	bot.Plugins.excuteFriend(msg)
}
