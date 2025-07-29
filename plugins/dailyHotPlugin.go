package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	text, mention := utils.ExtractMessageContent(message)
	if !mention {
		return true
	}
	trimText := strings.TrimSpace(text)
	if strings.HasPrefix(trimText, "/dailyhot") {
		sp := strings.Split(trimText, " ")
		if len(sp) != 2 {
			utils.SendShortReply(message, message.UserId, "参数数量错误 示例 /dailyhot bilibili")
			return false
		}
		target := sp[1]
		hot := getDailyHot(target)
		if hot == nil {
			utils.SendShortReply(message, message.UserId, "热搜获取失败，请检查参数")
			return false
		} else if hot.Code != 200 {
			utils.SendShortReply(message, message.UserId, "热搜获取失败，请检查热搜服务")
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
			if index > 9 {
				break
			}
		}
		chain := messagechain.Group(message.GroupId)
		chain.Text(responseText)
		chain.Send()
		return false
	}
	return true
}

type dailyHotMeta struct {
	Code int `json:"code"`
	Data []struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	} `json:"data"`
}

func getDailyHot(target string) *dailyHotMeta {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, "http://localhost:6688/"+target, nil)
	if err != nil {
		log.Println("error when create http request: plugin: DailyHotPlugin: getDailyHot", err)
		return nil
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("error when send http request: plugin: DailyHotPlugin: getDailyHot", err)
		return nil
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error when get http body: plugin: DailyHotPlugin: getDailyHot", err)
		return nil
	}
	hot := &dailyHotMeta{}
	err = json.Unmarshal(bytes, hot)
	if err != nil {
		log.Println("error when get json unmarshal: plugin: DailyHotPlugin: getDailyHot", err)
		return nil
	}
	return hot
}
