package botcontext

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

// 从私聊消息链中提取文本
func ExtractPrivateMessageText(msg *model.FriendMessage) string {
	text := ""
	for _, t := range msg.Message {
		if t.Type == "text" {
			text += t.Data["text"].(string)
		}
	}
	return text
}

// 从群聊消息链中提取文本和是否AT机器人
func ExtractMessageContent(msg *model.GroupMessage) (string, bool) {
	text := strings.Builder{}
	mention := false
	groupUserInfo := messagechain.GroupUserInfo{}

	for _, t := range msg.Message {
		switch t.Type {
		case "text":
			text.WriteString(t.Data["text"].(string))
		case "at":
			mentionUser, ok := t.Data["qq"].(string)
			if !ok {
				log.Println("error when get mentionUser: ExtractMessageContent")
				break
			}
			mentionUserId, err := strconv.Atoi(mentionUser)
			if err != nil {
				log.Println("error when get mentionUserId: ExtractMessageContent")
				break
			}
			if mentionUser == fmt.Sprintf("%d", msg.SelfId) {
				mention = true
			} else {
				card, err := groupUserInfo.GetUserInfo(uint(mentionUserId), msg.GroupId)
				if err == nil {
					showName := card.Card
					if card.Card == "" {
						showName = card.Nickname
					}
					text.WriteString(fmt.Sprintf("[@%s]", showName))
				} else {
					text.WriteString(mentionUser)
				}
			}
		}
	}
	return text.String(), mention
}

// 发送短回复
func SendShortReply(msg *model.GroupMessage, uid uint, response string) {
	chain := messagechain.Group(msg.GroupId)
	chain.Reply(msg.MessageId)
	chain.Mention(uid)
	chain.Text(" " + response)
	chain.Send()
}

// 发送长回复
func SendLongReply(msg *model.GroupMessage, rreply []rune, replyLength int) {
	forward := messagechain.GroupForward(msg.GroupId, "聊天记录", fmt.Sprintf("%d", msg.SelfId), "江颦")
	chain := messagechain.Group(msg.GroupId)
	chain.Mention(msg.UserId)
	chain.Send()

	for i := 0; i <= replyLength/500; i++ {
		start := i * 500
		end := (i + 1) * 500

		if end < replyLength {
			forward.Text(string(rreply[start:end]))
		} else if start < replyLength {
			forward.Text(string(rreply[start:]))
		}
	}

	time.Sleep(500 * time.Millisecond)
	forward.Send()
}
