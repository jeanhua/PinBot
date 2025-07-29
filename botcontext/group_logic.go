package botcontext

import (
	"strings"

	"github.com/jeanhua/PinBot/botcommand"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

// 处理群组消息
func (bot *BotContext) onGroupMessage(msg *model.GroupMessage) {
	text, mention := utils.ExtractMessageContent(msg)
	if strings.TrimSpace(text) == "" {
		return
	}

	trimText := strings.TrimSpace(text)

	// 处理指令
	if mention {
		if !botcommand.DealGroupCommand(trimText, msg) {
			return
		}
	}

	bot.Plugins.excuteGroup(msg)
}
