package utils

import (
	"fmt"

	"github.com/jeanhua/PinBot/model"
)

func ExtractPrivateMessageText(msg *model.FriendMessage) string {
	text := ""
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		}
	}
	return text
}

func ExtractMessageContent(msg *model.GroupMessage) (string, bool) {
	text := ""
	mention := false

	for _, t := range msg.Message {
		switch t.Type {
		case "text":
			text += t.Data["text"].(string)
		case "at":
			if t.Data["qq"].(string) == fmt.Sprintf("%d", msg.SelfId) {
				mention = true
			}
		}
	}
	return text, mention
}
