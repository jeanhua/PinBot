package functioncall

import (
	"fmt"
	"log"
	"time"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
)

type FunctionHandler interface {
	Handle(params map[string]any, sendVoice func(text string)) (string, error)
}

type FunctionCall struct {
	Action string         `json:"action"`
	Param  map[string]any `json:"parameters"`
}

const paramError = "变量类型错误"

var zanao = &utils.Zanao{}

// 函数注册表
var functionRegistry = map[string]FunctionHandler{
	"browseHomepage": &BrowseHomepageHandler{},
	"searchPost":     &SearchPostHandler{},
	"viewPost":       &ViewPostHandler{},
	"browseHot":      &BrowseHotHandler{},
	"viewComments":   &ViewCommentsHandler{},
	"speak":          &SpeakHandler{},
	"webSearch":      &WebSearchHandler{},
	"webExplore":     &WebExploreHandler{},
	"getCurrentTime": &GetCurrentTimeHandler{},
}

func CallFunction(name string, params map[string]any, sendVoice func(text string)) (string, error) {
	log.Println("call function: name:", name, "params", params)
	handler, ok := functionRegistry[name]
	if !ok {
		return "function call不匹配，请检查后重试", nil
	}
	return handler.Handle(params, sendVoice)
}

type BrowseHomepageHandler struct{}

func (h *BrowseHomepageHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	fromTime, err := getStringParam(params, "fromTime")
	if err != nil {
		return "", err
	}
	return zanao.GetNewest(fromTime), nil
}

type SearchPostHandler struct{}

func (h *SearchPostHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	keywords, err := getStringParam(params, "keywords")
	if err != nil {
		return "", err
	}
	return zanao.Search(keywords), nil
}

type ViewPostHandler struct{}

func (h *ViewPostHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	postId, err := getStringParam(params, "postId")
	if err != nil {
		return "", err
	}
	return zanao.GetDetail(postId), nil
}

type BrowseHotHandler struct{}

func (h *BrowseHotHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	return zanao.GetHot(), nil
}

type ViewCommentsHandler struct{}

func (h *ViewCommentsHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	postId, err := getStringParam(params, "postId")
	if err != nil {
		return "", err
	}
	return zanao.GetComments(postId), nil
}

type SpeakHandler struct{}

func (h *SpeakHandler) Handle(params map[string]any, sendVoice func(text string)) (string, error) {
	text, err := getStringParam(params, "text")
	if err != nil {
		return "", err
	}
	sendVoice(text)
	return "已成功给用户发送语音，你可以继续回复用户，或者输出一个空格结束", nil
}

type WebSearchHandler struct{}

func (h *WebSearchHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	query, err := getStringParam(params, "query")
	if err != nil {
		return "", err
	}

	timeRange := getOptionalStringParam(params, "timeRange")
	include, _ := params["include"].([]string)
	exclude, _ := params["exclude"].([]string)
	count := getIntParam(params, "count", 10)

	result := utils.WebSearch(config.ConfigInstance.TavilyToken, query, timeRange, include, exclude, count)
	return result, nil
}

type WebExploreHandler struct{}

func (h *WebExploreHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	linksInterface, ok := params["links"]
	if !ok {
		return "", fmt.Errorf("缺少参数: links")
	}

	links, err := convertToStringSlice(linksInterface)
	if err != nil {
		return "", err
	}

	result := utils.WebExplore(links, config.ConfigInstance.TavilyToken)
	return result, nil
}

type GetCurrentTimeHandler struct{}

func (h *GetCurrentTimeHandler) Handle(params map[string]any, _ func(text string)) (string, error) {
	now := time.Now().Local()
	return fmt.Sprintf("当前时间是 %d年%d月%d日 %d时%d分 %s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Weekday().String()), nil
}

func getStringParam(params map[string]any, key string) (string, error) {
	val, ok := params[key].(string)
	if !ok {
		return "", fmt.Errorf("参数 %s 类型错误或缺失", key)
	}
	return val, nil
}

func getOptionalStringParam(params map[string]any, key string) *string {
	if val, ok := params[key].(string); ok {
		return &val
	}
	return nil
}

func getIntParam(params map[string]any, key string, defaultValue int) int {
	if val, ok := params[key].(int); ok {
		return val
	}
	return defaultValue
}

func convertToStringSlice(linksInterface interface{}) ([]string, error) {
	linksSlice, ok := linksInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("参数 links 格式错误，应为字符串数组")
	}

	var links []string
	for _, v := range linksSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("links 数组中包含非字符串元素")
		}
		links = append(links, str)
	}
	return links, nil
}
