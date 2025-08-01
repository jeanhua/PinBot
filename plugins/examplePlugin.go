package plugins

import (
	"log"

	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

var ExamplePlugin = botcontext.NewPluginContext("example plugin", examplePluginOnFriend, examplePluginOnGroup, "示例插件")

func examplePluginOnFriend(message *model.FriendMessage) bool {
	log.Printf("[私聊消息](%s %d):%s\n", message.Sender.Nickname, message.Sender.UserId, utils.ExtractPrivateMessageText(message))
	return true
}
func examplePluginOnGroup(message *model.GroupMessage) bool {
	text, _ := utils.ExtractMessageContent(message)
	log.Printf("[群聊消息(%d)](%s):%s\n", message.GroupId, message.Sender.Nickname, text)
	return true
}
