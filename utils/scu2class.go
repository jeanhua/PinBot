package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/skip2/go-qrcode"
	"log"
	"strings"
	"time"
)

type SCU2Class struct {
	client *resty.Client
}

func NewSCU2Class(token string) *SCU2Class {
	client := resty.New()
	client.SetHeader("Token", token)
	return &SCU2Class{
		client: client,
	}
}

type scu2ClassSearchResponse struct {
	List []activity `json:"list"`
}

type activity struct {
	ActivityLibraryID string `json:"activityLibraryId"`
	Name              string `json:"name"`
	Describe          string `json:"describe"`
	Doing             bool   `json:"doing"`
}

func (s *SCU2Class) Search(key string) string {
	var result scu2ClassSearchResponse
	_, err := s.client.R().SetBody(map[string]any{
		"pn":   1,
		"ps":   10,
		"name": key,
	}).SetResult(&result).Post("https://zjczs.scu.edu.cn/ccyl-api/app/activity/list-activity-library")
	if err != nil {
		log.Println(err)
		return "请求失败"
	}
	var back strings.Builder
	now := time.Now().Local()
	back.WriteString(fmt.Sprintf("当前时间 %d年%d月%d日 %d时%d分 %s\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Weekday().String()))
	back.WriteString("检索到以下系列活动：\n")
	for _, item := range result.List {
		if !item.Doing {
			continue
		}
		back.WriteString(fmt.Sprintf("活动系列id：%s\n", item.ActivityLibraryID))
		back.WriteString(fmt.Sprintf("活动系列名称：%s\n", item.Name))
		back.WriteString(fmt.Sprintf("活动系列描述：%s\n------\n", item.Describe))
	}
	return back.String()
}

type scu2ClassListResponse struct {
	Activities []struct {
		ActivityID string `json:"activityId"`
		Name       string `json:"activityName"`
		Address    string `json:"activityAddress"`
		StartDate  string `json:"startTime"`
		EndDate    string `json:"endTime"`
	} `json:"activities"`
}

func (s *SCU2Class) List(id string) string {
	var result scu2ClassListResponse
	_, err := s.client.R().SetResult(&result).Post("https://zjczs.scu.edu.cn/ccyl-api/app/activity/get-lib-detail/" + id)
	if err != nil {
		log.Println(err)
		return "查询失败"
	}
	var back strings.Builder
	now := time.Now().Local()
	back.WriteString(fmt.Sprintf("当前时间 %d年%d月%d日 %d时%d分 %s\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Weekday().String()))
	back.WriteString("系列活动包含如下：\n")
	for _, item := range result.Activities {
		back.WriteString(fmt.Sprintf("活动ID: %s\n", item.ActivityID))
		back.WriteString(fmt.Sprintf("活动名称: %s\n", item.Name))
		back.WriteString(fmt.Sprintf("活动地址: %s\n", item.Address))
		back.WriteString(fmt.Sprintf("活动开始时间: %s\n", item.StartDate))
		back.WriteString(fmt.Sprintf("活动结束时间: %s\n------\n", item.EndDate))
	}
	return back.String()
}

func (s *SCU2Class) GenQRCode(id string) (bs64_in string, bs64_out string, err error) {
	var png []byte
	png, err = qrcode.Encode("https://zjczs.scu.edu.cn/ccylmp/pages/main/index/signing?type=in&state=1&id="+id, qrcode.Medium, 256)
	if err != nil {
		log.Println(err)
		return
	}
	bs64_in = base64.StdEncoding.EncodeToString(png)
	png, err = qrcode.Encode("https://zjczs.scu.edu.cn/ccylmp/pages/main/index/signing?type=out&state=1&id="+id, qrcode.Medium, 256)
	if err != nil {
		log.Println(err)
		return
	}
	bs64_out = base64.StdEncoding.EncodeToString(png)
	return
}
