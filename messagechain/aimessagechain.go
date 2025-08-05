package messagechain

import (
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
)

type AIMessageData struct {
	urlpath   string
	GroupId   uint   `json:"group_id"`
	Character string `json:"character"`
	Text      string `json:"text"`
}

func AIMessage(groupUin uint, charactor string, text string) MessageChain {
	return &AIMessageData{
		urlpath:   config.GetConfig().NapCatServerUrl + "/send_group_ai_record",
		GroupId:   groupUin,
		Character: charactor,
		Text:      text,
	}
}

func (mc *AIMessageData) Send() {
	httpUtil := utils.HttpUtil{}
	err := httpUtil.RequestWithNoResponse(http.MethodPost, mc.urlpath, httpUtil.WithJsonBody(mc))
	if err != nil {
		log.Println("error when send aiMessage chain message")
	}
}
