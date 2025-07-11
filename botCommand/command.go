package botcommand

import (
	"strings"
	"sync"

	"github.com/jeanhua/PinBot/config"
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var (
	EnableAIAudio = false
	CommandMu     sync.RWMutex
)

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		chain.Text(config.ConfigInstance.HelpWords.Group)
		messageChain.SendMessage(chain)
		return true
	case "/enable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		CommandMu.Lock()
		EnableAIAudio = true
		CommandMu.Unlock()
		chain.Text(" 已开启AI语音")
		messageChain.SendMessage(chain)
		return true
	case "/disable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		CommandMu.Lock()
		EnableAIAudio = false
		CommandMu.Unlock()
		chain.Text(" 已关闭AI语音")
		messageChain.SendMessage(chain)
		return true
	}
	return false
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Friend(msg.UserId)
		chain.Text(config.ConfigInstance.HelpWords.Friend)
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}
