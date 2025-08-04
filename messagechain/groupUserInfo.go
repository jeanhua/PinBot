package messagechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/config"
)

type GroupUserInfo struct {
	Card     string
	Nickname string
}

type groupUserInfoModel struct {
	Data struct {
		Card     string `json:"card"`
		Nickname string `json:"nickname"`
	} `json:"data"`
}

func (g *GroupUserInfo) GetUserInfo(userid, groupid uint) (*GroupUserInfo, error) {
	client := http.Client{}
	body := map[string]any{}
	body["group_id"] = groupid
	body["user_id"] = userid
	body["no_cache"] = true
	bs, err := json.Marshal(&body)
	if err != nil {
		log.Println("error when json marsharl: GetUserInfo")
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, config.GetConfig().NapCatServerUrl+"/get_group_member_info", bytes.NewReader(bs))
	if err != nil {
		log.Println("error when create http request: GetUserInfo")
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Println("error when send http request: GetUserInfo")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error status code: GetUserInfo")
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error when read http response: GetUserInfo")
		return nil, err
	}
	result := groupUserInfoModel{}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		log.Println("error when json unmarsharl: GetUserInfo")
		return nil, err
	}
	return &GroupUserInfo{
		Nickname: result.Data.Nickname,
		Card:     result.Data.Card,
	}, nil
}
