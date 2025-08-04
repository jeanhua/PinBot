package plugins

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

var DailyHotPlugin = botcontext.NewPluginContext("daily hot", dailyHotOnFriend, dailyHotOnGroup, "每日热搜，发送 /dailyhot help 查看支持列表")

func dailyHotOnFriend(message *model.FriendMessage) bool {
	return true
}

func dailyHotOnGroup(message *model.GroupMessage) bool {
	text, mention := botcontext.ExtractMessageContent(message)
	if !mention {
		return true
	}
	trimText := strings.TrimSpace(text)
	if !strings.HasPrefix(trimText, "/dailyhot") {
		return true
	}
	target, err := getDailyHotParam(trimText)
	if err != nil {
		botcontext.SendShortReply(message, message.UserId, err.Error())
		return false
	}
	if target == "help" {
		sendDailyHotHelp(message.GroupId)
		return false
	}
	hot := getDailyHot(target)
	if hot == nil {
		botcontext.SendShortReply(message, message.UserId, "热搜获取失败，请检查参数")
		return false
	} else if hot.Code != 200 {
		botcontext.SendShortReply(message, message.UserId, "热搜获取失败，请检查热搜服务")
		return false
	}
	responseText := ""
	dataLen := len(hot.Data)
	for index, v := range hot.Data {
		if index != dataLen-1 && index < 9 {
			responseText += fmt.Sprintf("%d: %s\n%s\n\n", index+1, v.Title, v.Url)
		} else {
			responseText += fmt.Sprintf("%d: %s\n%s", index+1, v.Title, v.Url)
		}
		if index >= 9 {
			break
		}
	}
	chain := messagechain.Group(message.GroupId)
	chain.Text(responseText)
	chain.Send()
	return false
}

func sendDailyHotHelp(uid uint) {
	chain := messagechain.Group(uid)
	chain.LocalImage("./Pluginres/dailyhot/help.png")
	chain.Send()
}

func getDailyHotParam(text string) (string, error) {
	sp := strings.Split(text, " ")
	if len(sp) != 2 {
		return "", fmt.Errorf("参数数量错误 示例 /dailyhot bilibili")
	}
	return sp[1], nil
}

type dailyHotMeta struct {
	Code int `json:"code"`
	Data []struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	} `json:"data"`
}

func getDailyHot(target string) *dailyHotMeta {
	hot := dailyHotMeta{}
	httpUtil := utils.HttpUtil{}
	err := httpUtil.Request(http.MethodGet, "http://localhost:6688/"+target, nil, &hot)
	if err != nil {
		return nil
	}
	return &hot
}
