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
			return "已成功给用户发送语音", nil
		} else {
			return "", fmt.Errorf(paramError)
		}
	case "webSearch":
		query, ok := param["query"].(string)
		if !ok {
			return "", fmt.Errorf("缺少或无效的搜索关键词 'query'")
		}
		freshness, _ := param["freshness"].(string)
		include, _ := param["include"].(string)
		exclude, _ := param["exclude"].(string)

		summary := false
		if s, ok := param["summary"].(bool); ok {
			summary = s
		}
		count := 10 // 默认值
		if c, ok := param["count"].(int); ok {
			count = c
		}
		result := utils.WebSearch(config.ConfigInstance.BochaToken, query, freshness, summary, include, exclude, count)
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
