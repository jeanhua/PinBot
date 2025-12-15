package messagechain

import (
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/utils"
)

type AIMessageData struct {
	urlPath   string
	GroupId   uint   `json:"group_id"`
	Character string `json:"character"`
	Text      string `json:"text"`
}

func AIMessage(groupUin uint, character string, text string) MessageChain {
	return &AIMessageData{
		urlPath:   config.GetConfig().GetString("bot_config.napcatServerUrl") + "/send_group_ai_record",
		GroupId:   groupUin,
		Character: character,
		Text:      text,
	}
}

func (mc *AIMessageData) Send() {
	httpUtil := utils.HttpUtil{}
	err := httpUtil.RequestWithNoResponse(http.MethodPost, mc.urlPath, httpUtil.WithJsonBody(mc))
	if err != nil {
		log.Println("error when send aiMessage chain message")
	}
}
