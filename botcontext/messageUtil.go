package botcontext

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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
func ExtractGroupMessageContent(msg *model.GroupMessage) (string, bool) {
	mention := false
	for _, t := range msg.Message {
		switch t.Type {
		case "at":
			mentionUser, ok := t.Data["qq"].(string)
			if !ok {
				log.Println("error when get mentionUser: ExtractMessageContent")
				break
			}
			if mentionUser == fmt.Sprintf("%d", msg.SelfId) {
				mention = true
			}
		}
	}
	return extractOB11SegmentMessage(msg.Message, msg.GroupId, 3), mention
}

func ExtractGroupRawMessage(msg *model.GroupMessage) (string, bool) {
	mention := false
	text := strings.Builder{}
	for _, t := range msg.Message {
		switch t.Type {
		case "at":
			mentionUser, ok := t.Data["qq"].(string)
			if !ok {
				log.Println("error when get mentionUser: ExtractMessageContent")
				break
			}
			if mentionUser == fmt.Sprintf("%d", msg.SelfId) {
				mention = true
			}
		case "text":
			text.WriteString(t.Data["text"].(string))
		}
	}
	return text.String(), mention
}

// 消息链转文本
func extractOB11SegmentMessage(segment []model.OB11Segment, groupid uint, limit int) string {
	if limit <= 0 {
		return ""
	}
	result := strings.Builder{}
	groupUserInfo := messagechain.GroupUserInfo{}
	for _, s := range segment {
		switch s.Type {
		case "text":
			result.WriteString(s.Data["text"].(string))
		case "at":
			mentionUser, ok := s.Data["qq"].(string)
			if !ok {
				log.Println("error when get mentionUser: ExtractMessageContent")
				break
			}
			mentionUserId, err := strconv.Atoi(mentionUser)
			if err != nil {
				log.Println("error when get mentionUserId: ExtractMessageContent")
				break
			}
			card, err := groupUserInfo.GetUserInfo(uint(mentionUserId), groupid)
			if err == nil {
				showName := card.Card
				if card.Card == "" {
					showName = card.Nickname
				}
				result.WriteString(fmt.Sprintf("[@%s]", showName))
			} else {
				result.WriteString(mentionUser)
			}
		case "reply":
			msgDetail := messagechain.ReplyMessageInfo{}
			messageIdStr, ok := s.Data["id"].(string)
			if !ok {
				log.Println("error when get messageIdStr: ExtraOB11SegmentMessage")
				continue
			}
			messageId, err := strconv.Atoi(messageIdStr)
			if err != nil {
				log.Println("error when convert messageIdStr: ExtraOB11SegmentMessage")
				continue
			}
			seg, err := msgDetail.GetMessageDetail(uint(messageId))
			if err != nil {
				log.Println("error when getMessageDetail: ExtraOB11SegmentMessage")
				continue
			}
			next := extractOB11SegmentMessage(seg, groupid, limit-1)
			result.WriteString("\n---↓以下为回复的消息↓---\n")
			result.WriteString(next)
			result.WriteString("\n---↑以上为回复的消息↑---\n")
		case "json":
			jsonMap := model.JsonMessage{}
			err := json.Unmarshal([]byte(s.Data["data"].(string)), &jsonMap)
			if err != nil {
				log.Println("error when json unmarsharl: json message: ExtraOB11SegmentMessage")
				continue
			}
			result.WriteString(fmt.Sprintf("[分享卡片,标题: %s,描述: %s,链接: (%s)]", jsonMap.Meta.News.Title, jsonMap.Meta.News.Desc, jsonMap.Meta.News.JumpUrl))
		}
	}
	return result.String()
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
func SendLongReply(msg *model.GroupMessage, reply string) {
	forward := messagechain.GroupForward(msg.GroupId, "聊天记录", fmt.Sprintf("%d", msg.SelfId), "江颦")
	chain := messagechain.Group(msg.GroupId)
	chain.Mention(msg.UserId)
	chain.Send()
	rreply := []rune(reply)
	for i := 0; i < len(rreply); i += 500 {
		end := i + 500
		if end > len(rreply) {
			end = len(rreply)
		}
		segment := string(rreply[i:end])
		forward.Text(segment)
	}
	forward.Send()
}
