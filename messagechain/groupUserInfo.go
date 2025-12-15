package messagechain

import (
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
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
	body := map[string]any{}
	body["group_id"] = groupid
	body["user_id"] = userid
	body["no_cache"] = true
	result := groupUserInfoModel{}
	httpUtil := utils.HttpUtil{}
	err := httpUtil.Request(http.MethodPost, config.GetConfig().GetString("bot_config.napcatServerUrl")+"/get_group_member_info", httpUtil.WithJsonBody(&body), &result)
	if err != nil {
		log.Println("error when httpUtil.Request: GetUserInfo")
		return nil, err
	}
	return &GroupUserInfo{
		Nickname: result.Data.Nickname,
		Card:     result.Data.Card,
	}, nil
}
