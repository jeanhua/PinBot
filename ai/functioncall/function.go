package functioncall

import (
	"fmt"
	"log"

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
		return browseHomepage(), nil
	case "search":
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
	default:
		return "", fmt.Errorf("调用了无匹配的function call: %s", name)
	}
}

var zanao *utils.Zanao = &utils.Zanao{}

func browseHomepage() string {
	log.Println("调用 browseHomepage")
	return zanao.GetNewest()
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
