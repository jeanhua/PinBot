package botplugin

import (
	"log"

	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

func ExampleFriendPlugin(message *model.FriendMessage) bool {
	log.Printf("[私聊消息](%s):%s\n", message.Sender.Nickname, utils.ExtractPrivateMessageText(message))
	return true
}
func ExampleGroupPlugin(message *model.GroupMessage) bool {
	text, _ := utils.ExtractMessageContent(message)
	log.Printf("[群聊消息](%s):%s\n", message.Sender.Nickname, text)
	return true
}
