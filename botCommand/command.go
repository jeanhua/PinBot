package botcommand

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jeanhua/PinBot/config"
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
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
		config.ConfigInstance_mu.RLock()
		chain.Text(config.ConfigInstance.HelpWords.Group)
		config.ConfigInstance_mu.RUnlock()
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
	case "/zanao post":
		zanao := &utils.Zanao{}
		resp, err := zanao.GetNewest()
		if err != nil {
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text("我遇到了一点错误，请稍后再试")
			messageChain.SendMessage(chain)
			return true
		}
		groupForward := messageChain.GroupForward(msg.GroupId, "集市最新帖子")
		for _, v := range resp.Data.List {
			groupForward.Text(fmt.Sprintf("%s\n%s", v.Title, v.Content), msg.SelfId, "江颦")
		}
		groupForward.Send()
		return true
	case "/zanao hot":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		zanao := &utils.Zanao{}
		resp, err := zanao.GetHot()
		if err != nil {
			chain.Text("我遇到了一点错误，请稍后再试")
			messageChain.SendMessage(chain)
			return true
		}
		text := "实时热帖：\n"
		for i, v := range resp.Data.List {
			text += fmt.Sprintf("[%d]%s\n", i+1, v.Title)
		}
		text = strings.TrimSpace(text)
		chain.Text(text)
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Friend(msg.UserId)
		config.ConfigInstance_mu.RLock()
		chain.Text(config.ConfigInstance.HelpWords.Friend)
		config.ConfigInstance_mu.RUnlock()
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}
