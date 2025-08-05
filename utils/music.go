package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MusicSearch struct{}

type musicSearchResult struct {
	Result struct {
		Songs []struct {
			Id      uint   `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"songs"`
	} `json:"result"`
}

func (*MusicSearch) Search(keyword string) string {
	httpUtil := HttpUtil{}
	result := musicSearchResult{}
	encodeStr := url.QueryEscape(keyword)
	err := httpUtil.Request(http.MethodGet, fmt.Sprintf("https://163api.qijieya.cn/search?keywords=%s&limit=10", encodeStr), nil, &result)
	if err != nil {
		log.Println("error when search music")
		return "搜索失败"
	}
	str := strings.Builder{}
	for _, song := range result.Result.Songs {
		artistNames := ""
		for _, artist := range song.Artists {
			artistNames += artist.Name + " "
		}
		str.WriteString(fmt.Sprintf("[id:%d] %s - %s\n", song.Id, song.Name, artistNames))
	}
	return str.String()
}
