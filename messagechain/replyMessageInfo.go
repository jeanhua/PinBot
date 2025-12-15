package messagechain

import (
	"fmt"
	"net/http"

	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

type ReplyMessageInfo struct{}

func (*ReplyMessageInfo) GetMessageDetail(messageId uint) ([]model.OB11Segment, error) {
	httpUtil := utils.HttpUtil{}
	body := map[string]uint{
		"message_id": messageId,
	}
	result := model.MessageDetail{}
	err := httpUtil.Request(http.MethodPost, config.GetConfig().GetString("bot_config.napcatServerUrl")+"/get_msg", httpUtil.WithJsonBody(&body), &result)
	if err != nil {
		return nil, fmt.Errorf("error when GetMessageDetail")
	}
	return result.Data.Message, nil
}
