package botcommand

import (
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var (
	EnableAIAudio = true
)

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(" @我发送 清除记录 ->清除聊天记录\n@我发送 /enable(disable) AI语音 ->开关AI语音")
		messageChain.SendMessage(chain)
		return true
	case "/enable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		EnableAIAudio = true
		chain.Text(" 已开启AI语音")
		messageChain.SendMessage(chain)
		return true
	case "/disable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		EnableAIAudio = false
		chain.Text(" 已关闭AI语音")
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Friend(msg.UserId)
		chain.Text(" @我发送 清除记录 可以清除聊天记录哦")
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}
