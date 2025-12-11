package botcontext

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/model"
)

// 处理接收到的消息
func HandleMessage(message []byte, bot *BotContext) {
	basicMsg, err := parseBasicMessage(message)
	if err != nil || basicMsg.PostType != "message" {
		return
	}
	if isFriendMessage(message) {
		handleFriendMessage(message, bot)
	} else if isGroupMessage(message) {
		handleGroupMessage(message, bot)
	}
}

// 解析基本消息结构
func parseBasicMessage(message []byte) (*model.Message, error) {
	basic := &model.Message{}
	if err := json.Unmarshal(message, basic); err != nil {
		log.Printf("Error parsing basic message: %v", err)
		return nil, err
	}
	return basic, nil
}

// 检查是否为好友消息
func isFriendMessage(message []byte) bool {
	friendMsg := model.FriendMessage{}
	if err := json.Unmarshal(message, &friendMsg); err != nil {
		log.Printf("Error parsing friend message: %v", err)
		return false
	}
	return friendMsg.MessageType == "private"
}

// 检查是否为群组消息
func isGroupMessage(message []byte) bool {
	groupMsg := model.GroupMessage{}
	if err := json.Unmarshal(message, &groupMsg); err != nil {
		log.Printf("Error parsing group message: %v", err)
		return false
	}
	return groupMsg.MessageType == "group"
}

// 处理好友消息
func handleFriendMessage(message []byte, bot *BotContext) {
	friendMsg := model.FriendMessage{}
	if err := json.Unmarshal(message, &friendMsg); err != nil {
		log.Printf("Error parsing friend message: %v", err)
		return
	}

	if isExcludedFriend(friendMsg.UserId) {
		return
	}

	if friendMsg.Sender.UserId != friendMsg.SelfId && isIncludedFriend(friendMsg.UserId) {
		bot.onPrivateMessage(&friendMsg)
	}
}

// 处理群组消息
func handleGroupMessage(message []byte, bot *BotContext) {
	groupMsg := model.GroupMessage{}
	if err := json.Unmarshal(message, &groupMsg); err != nil {
		log.Printf("Error parsing group message: %v", err)
		return
	}

	if isExcludedGroup(groupMsg.GroupId) {
		return
	}

	if groupMsg.Sender.UserId != groupMsg.SelfId && isIncludedGroup(groupMsg.GroupId) {
		bot.onGroupMessage(&groupMsg)
	}
}

// 检查好友是否在排除列表中
func isExcludedFriend(userId uint) bool {
	for _, uin := range config.GetConfig().Friend.Exclude {
		if uin == fmt.Sprintf("%d", userId) {
			return true
		}
	}
	return false
}

// 检查好友是否在包含列表中
func isIncludedFriend(userId uint) bool {
	for _, uin := range config.GetConfig().Friend.Include {
		if uin == "all" || uin == fmt.Sprintf("%d", userId) {
			return true
		}
	}
	return false
}

// 检查群组是否在排除列表中
func isExcludedGroup(groupId uint) bool {
	for _, uin := range config.GetConfig().Group.Exclude {
		if uin == fmt.Sprintf("%d", groupId) {
			return true
		}
	}
	return false
}

// 检查群组是否在包含列表中
func isIncludedGroup(groupId uint) bool {
	for _, uin := range config.GetConfig().Group.Include {
		if uin == "all" || uin == fmt.Sprintf("%d", groupId) {
			return true
		}
	}
	return false
}
