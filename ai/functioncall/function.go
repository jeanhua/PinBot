package functioncall

import (
	"fmt"
	"log"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
)

type FunctionCall struct {
	Action string         `json:"action"`
	Param  map[string]any `json:"parameters"`
}

const paramError = "变量类型错误"

func CallFunction(name string, param map[string]any, sendVoice func(text string)) (string, error) {
	log.Println("call function: name:", name, "param", param)
	switch name {
	case "browseHomepage":
		fromTime, ok := param["fromTime"].(string)
		if ok {
			return browseHomepage(fromTime), nil
		} else {
			return "", fmt.Errorf(paramError)
		}

	case "searchPost":
		keywords, ok := param["keywords"].(string)
		if ok {
			return search(keywords), nil
		} else {
			return "", fmt.Errorf(paramError)
		}

	case "viewPost":
		postId, ok := param["postId"].(string)
		if ok {
			return viewPost(postId), nil
		} else {
			return "", fmt.Errorf(paramError)
		}
	case "browseHot":
		return browseHot(), nil
	case "viewComments":
		postId, ok := param["postId"].(string)
		if ok {
			return viewComments(postId), nil
		} else {
			return "", fmt.Errorf(paramError)
		}
	case "speak":
		text, ok := param["text"].(string)
		if ok {
			sendVoice(text)
			return "已成功给用户发送语音，你可以继续回复用户，或者输出一个空格结束", nil
		} else {
			return "", fmt.Errorf(paramError)
		}
	case "webSearch":
		query, ok := param["query"].(string)
		if !ok {
			return "", fmt.Errorf(paramError)
		}

		var timeRange *string
		if tr, ok := param["timeRange"].(string); ok {
			timeRange = &tr
		} else {
			timeRange = nil
		}

		include, _ := param["include"].([]string)
		exclude, _ := param["exclude"].([]string)

		count, ok := param["count"].(int)
		if !ok {
			count = 10
		}
		result := utils.WebSearch(config.ConfigInstance.TavilyToken, query, timeRange, include, exclude, count)
		return result, nil
	case "webExplore":
		linksInterface, ok := param["links"]
		if !ok {
			return "", fmt.Errorf("缺少参数: links")
		}
		linksSlice, ok := linksInterface.([]interface{})
		if !ok {
			return "", fmt.Errorf("参数 links 格式错误，应为字符串数组")
		}
		var links []string
		for _, v := range linksSlice {
			str, ok := v.(string)
			if !ok {
				return "", fmt.Errorf("links 数组中包含非字符串元素")
			}
			links = append(links, str)
		}
		result := utils.WebExplore(links, config.ConfigInstance.TavilyToken)
		return result, nil
	default:
		return "function call不匹配，请检查后重试", nil
	}
}

var zanao *utils.Zanao = &utils.Zanao{}

func browseHomepage(fromTime string) string {
	log.Println("调用 browseHomepage")
	return zanao.GetNewest(fromTime)
}

func search(keywords string) string {
	log.Println("调用 search")
	return zanao.Search(keywords)
}

func viewPost(postId string) string {
	log.Println("调用 viewPost")
	return zanao.GetDetail(postId)
}

func browseHot() string {
	log.Println("调用 browseHot")
	return zanao.GetHot()
}

func viewComments(postId string) string {
	log.Println("调用 viewComments")
	return zanao.GetComments(postId)
}
