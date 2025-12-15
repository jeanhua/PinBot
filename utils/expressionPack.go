package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jeanhua/PinBot/config"
	"net/url"
	"strings"
	"time"
)

type expressPackResult struct {
	Code int      `json:"code"`
	Data []string `json:"res"`
}

func SearchExpressPack(keyword string) string {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	result := new(expressPackResult)
	encodeKey := url.QueryEscape(keyword)
	_, err := client.R().SetResult(result).Get(fmt.Sprintf("http://101.35.2.25/api/img/apihzbqbbaidu.php?id=%s&key=%s&limit=10&words=%s", config.GetConfig().GetString("function_call_config.express_pack.id"), config.GetConfig().GetString("function_call_config.express_pack.key"), encodeKey))
	if err != nil {
		return "查找失败，请稍后再试"
	}
	searchString := strings.Builder{}
	searchString.WriteString("查询到如下表情包[" + keyword + "]\n")
	for _, word := range result.Data {
		searchString.WriteString(word + "\n")
	}
	searchString.WriteString("\n")
	return searchString.String()
}
