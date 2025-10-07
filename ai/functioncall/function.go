package functioncall

import (
	"fmt"
	"log"
	"time"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/utils"
)

type FunctionHandler interface {
	Handle(params map[string]any, uid uint, target int) (string, error)
}

type FunctionCall struct {
	Action string         `json:"action"`
	Param  map[string]any `json:"parameters"`
}

var zanao = &utils.Zanao{}

// 函数注册表
var functionRegistry = map[string]FunctionHandler{
	"browseHomepage": &browseHomepageHandler{},
	"searchPost":     &searchPostHandler{},
	"viewPost":       &viewPostHandler{},
	"browseHot":      &browseHotHandler{},
	"viewComments":   &viewCommentsHandler{},
	"speak":          &speakHandler{},
	"webSearch":      &webSearchHandler{},
	"webExplore":     &webExploreHandler{},
	"getCurrentTime": &getCurrentTimeHandler{},
	"hateImage":      &hateImageHandler{},   // 讨厌表情包
	"searchMusic":    &searchMusicHandler{}, // 搜索网易云音乐
	"shareMusic":     &shareMusicHandler{},  // 分享网易云音乐
	// 第二课堂相关
	"scu2ClassSearch": &scu2ClassSearchHandler{},
	"scu2ClassList":   &scu2ClassListHandler{},
	"scu2ClassShare":  &scu2ClassShareHandler{},
}

const (
	TargetGroup = iota
	TargetFriend
)

func CallFunction(name string, params map[string]any, uid uint, target int) (string, error) {
	log.Println("call function: name:", name, "params", params)
	handler, ok := functionRegistry[name]
	if !ok {
		return "function call不匹配，请检查后重试", nil
	}
	return handler.Handle(params, uid, target)
}

type browseHomepageHandler struct{}

func (h *browseHomepageHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	fromTime, err := getStringParam(params, "fromTime")
	if err != nil {
		return "", err
	}
	return zanao.GetNewest(fromTime), nil
}

type searchPostHandler struct{}

func (h *searchPostHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	keywords, err := getStringParam(params, "keywords")
	if err != nil {
		return "", err
	}
	return zanao.Search(keywords), nil
}

type viewPostHandler struct{}

func (h *viewPostHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	postId, err := getStringParam(params, "postId")
	if err != nil {
		return "", err
	}
	return zanao.GetDetail(postId), nil
}

type browseHotHandler struct{}

func (h *browseHotHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	return zanao.GetHot(), nil
}

type viewCommentsHandler struct{}

func (h *viewCommentsHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	postId, err := getStringParam(params, "postId")
	if err != nil {
		return "", err
	}
	return zanao.GetComments(postId), nil
}

type speakHandler struct{}

func (h *speakHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	text, err := getStringParam(params, "text")
	if err != nil {
		return "", err
	}
	switch target {
	case TargetGroup:
		chain := messagechain.AIMessage(uid, "lucy-voice-suxinjiejie", text)
		chain.Send()
	case TargetFriend:
		chain := messagechain.Friend(uid)
		chain.Text("[语音消息]" + text)
		chain.Send()
	}
	return "已成功给用户发送语音，你可以继续回复用户，或者输出一个空格结束", nil
}

type webSearchHandler struct{}

func (h *webSearchHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	query, err := getStringParam(params, "query")
	if err != nil {
		return "", err
	}

	timeRange := getOptionalStringParam(params, "timeRange")
	include, _ := params["include"].([]string)
	exclude, _ := params["exclude"].([]string)
	count := getIntParam(params, "count", 10)

	result := utils.WebSearch(config.GetConfig().TavilyToken, query, timeRange, include, exclude, count)
	return result, nil
}

type webExploreHandler struct{}

func (h *webExploreHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	linksInterface, ok := params["links"]
	if !ok {
		return "", fmt.Errorf("缺少参数: links")
	}

	links, err := convertToStringSlice(linksInterface)
	if err != nil {
		return "", err
	}

	result := utils.WebExplore(links, config.GetConfig().TavilyToken)
	return result, nil
}

type getCurrentTimeHandler struct{}

func (h *getCurrentTimeHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	now := time.Now().Local()
	return fmt.Sprintf("当前时间是 %d年%d月%d日 %d时%d分 %s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Weekday().String()), nil
}

type hateImageHandler struct{}

func (h *hateImageHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	userId, err := getStringParam(params, "userid")
	if err != nil {
		return "", err
	}
	switch target {
	case TargetGroup:
		chain := messagechain.Group(uid)
		chain.UrlImage("https://api.mhimg.cn/api/biaoqingbao_pa?qq=" + userId)
		chain.Send()
	case TargetFriend:
		chain := messagechain.Friend(uid)
		chain.UrlImage("https://api.mhimg.cn/api/biaoqingbao_pa?qq=" + userId)
		chain.Send()
	}
	return "发送成功", nil
}

type searchMusicHandler struct{}

func (*searchMusicHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	keyword, err := getStringParam(params, "query")
	if err != nil {
		return "", err
	}
	searchMusicUtil := utils.MusicSearch{}
	result := searchMusicUtil.Search(keyword)
	return result, nil
}

type shareMusicHandler struct{}

func (*shareMusicHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	id, err := getStringParam(params, "id")
	if err != nil {
		return "", err
	}
	if target == TargetGroup {
		chain := messagechain.Group(uid)
		chain.Music(id)
		chain.Send()
	} else {
		chain := messagechain.Friend(uid)
		chain.Music(id)
		chain.Send()
	}
	return "分享成功", nil
}

// 四川大学第二课堂
type scu2ClassSearchHandler struct{}

func (*scu2ClassSearchHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	scu2class := utils.NewSCU2Class(config.GetConfig().SCU2ClassToken)
	keyword, err := getStringParam(params, "keyword")
	if err != nil {
		return "", err
	}
	result := scu2class.Search(keyword)
	return result, nil
}

type scu2ClassListHandler struct{}

func (*scu2ClassListHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	scu2class := utils.NewSCU2Class(config.GetConfig().SCU2ClassToken)
	activityLibId, err := getStringParam(params, "activityLibId")
	if err != nil {
		return "", err
	}
	result := scu2class.List(activityLibId)
	return result, nil
}

type scu2ClassShareHandler struct{}

func (*scu2ClassShareHandler) Handle(params map[string]any, uid uint, target int) (string, error) {
	scu2class := utils.NewSCU2Class(config.GetConfig().SCU2ClassToken)
	activityId, err := getStringParam(params, "activityId")
	if err != nil {
		return "", err
	}
	qrIn, qrOut, err := scu2class.GenQRCode(activityId)
	if err != nil {
		return "", err
	}
	if target == TargetGroup {
		chain := messagechain.Group(uid)
		chain.Text("签到码:\n")
		chain.Base64Image(qrIn)
		chain.Text("签退码:\n")
		chain.Base64Image(qrOut)
		chain.Send()
	} else {
		chain := messagechain.Friend(uid)
		chain.Text("签到码:\n")
		chain.Base64Image(qrIn)
		chain.Text("签退码:\n")
		chain.Base64Image(qrOut)
		chain.Send()
	}
	return "发送成功", nil
}

// tools

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
