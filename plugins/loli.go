package plugins

import (
	"log"
	"strings"
	"sync"

	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

var LoliPlugin = botcontext.NewPluginContext("loli", loliPluginOnFriend, loliPluginOnGroup, "二次元萝莉插件(群聊)，发送 /loli 获取随机图片")

var (
	loliLock = sync.Mutex{}
)

func loliPluginOnFriend(message *model.FriendMessage) bool {
	return true
}
func loliPluginOnGroup(message *model.GroupMessage) bool {
	text, mention := botcontext.ExtractMessageContent(message)
	if !mention {
		return true
	}
	trimText := strings.TrimSpace(text)
	if trimText == "/loli" {
		if !loliLock.TryLock() {
			botcontext.SendShortReply(message, message.UserId, "反应不过来了，待会再试😘")
			return false
		}
		defer loliLock.Unlock()
		sendGroupLoliImage(message.GroupId)
		return false
	}
	return true
}

func sendGroupLoliImage(uid uint) {
	chain := messagechain.Group(uid)
	log.Println("发送图片: loli")
	chain.UrlImage("https://www.loliapi.com/acg/pe/")
	chain.Send()
}
