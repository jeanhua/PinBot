package aicommunicate

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/jeanhua/PinBot/ai/functioncall"
	"github.com/jeanhua/PinBot/config"
)

type ZhipuAIBot_z1_flash struct {
	token        string
	SystemPrompt string
	messageChain []*Message
}

const toolPrompt = `正常情况下正常聊天就好，但是当你想要使用工具时，或者用户提到的问题和下面的工具相关时，你可以使用工具
你有以下可用工具：

## 校园集市论坛相关：
- 'browse_homepage()'：访问主页，获取最新的帖子列表。
- 'browse_hot()'：访问24小时内热度最高的帖子列表。
- 'search(keyword:string)'：使用关键词搜索相关帖子。
- 'view_post(post_id:string)'：查看某篇帖子的详细内容。
- 'view_comments(post_id:string)'：查看某篇帖子的评论区内容。

当你要执行某个操作时，请以'#Call' + 严格的 JSON 格式输出你的动作和参数，例如：

#Call
{
  "action": "search",
  "parameters": {
    "keyword": "图书馆"
  }
}

注意事项：
1.请确保每次响应只包含一个动作，并且不要添加任何额外解释。我会根据你的指令执行操作并将结果反馈给你。
2.不要连续多次调用#Call，调用几次后就回答问题
3.有些#Call调用在短时间内是不变的，比如热帖，评论，帖子详情，请求过了就不要重复请求了
`

func NewZhipu(token string, prompt string) *ZhipuAIBot_z1_flash {
	zp := &ZhipuAIBot_z1_flash{
		token:        token,
		SystemPrompt: prompt + "\n" + toolPrompt,
	}
	zp.messageChain = []*Message{
		{
			Role:    "system",
			Content: zp.SystemPrompt,
		},
	}
	return zp
}

const requestUrl string = "https://open.bigmodel.cn/api/paas/v4/chat/completions"

type RequestBody struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}

type ResponseBody struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type thinkingResponse struct {
	thinking string
	text     string
}

func request(msg []*Message, token string) *thinkingResponse {
	client := &http.Client{}
	body := &RequestBody{
		Model:    "glm-z1-flash",
		Messages: msg,
		Stream:   false,
	}
	postBytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	request, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(postBytes))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		log.Println(err)
		return nil
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	var response ResponseBody
	json.Unmarshal(respBytes, &response)
	sort.Slice(response.Choices, func(i, j int) bool {
		return response.Choices[i].Index < response.Choices[j].Index
	})
	var respText string
	for _, v := range response.Choices {
		respText += v.Message.Content
	}
	splitText := strings.Split(respText, "</think>")
	if len(splitText) == 2 {
		return &thinkingResponse{
			thinking: strings.TrimSpace(splitText[0]),
			text:     strings.TrimSpace(splitText[1]),
		}
	} else {
		return nil
	}
}

func (zp *ZhipuAIBot_z1_flash) Ask(question string) *AiAnswer {
	question = strings.TrimSpace(question)
	if strings.HasPrefix(question, "#新对话") {
		zp.messageChain = []*Message{
			{
				Role:    "system",
				Content: zp.SystemPrompt,
			},
		}
	}
	zp.messageChain = append(zp.messageChain, &Message{
		Role:    "user",
		Content: question,
	})
	isFunccall := false
	funcCallNums := 0
	funcs := []*functioncall.FunctionCall{}
	for {
		resp := request(zp.messageChain, zp.token)
		if resp == nil {
			return nil
		}
		//log.Println(resp.text)
		zp.messageChain = append(zp.messageChain, &Message{
			Role:    "assistant",
			Content: resp.text,
		})
		if strings.HasPrefix(resp.text, "#Call") {
			if funcCallNums >= config.ConfigInstance.FunctionCallMaxC {
				zp.messageChain = append(zp.messageChain, &Message{
					Role:    "user",
					Content: "你的调用次数已达限制，请先回答用户问题",
				})
				continue
			}
			funcCallNums += 1
			funccallText := strings.TrimPrefix(resp.text, "#Call")
			var funccall functioncall.FunctionCall
			err := json.Unmarshal([]byte(funccallText), &funccall)
			if err != nil {
				zp.messageChain = append(zp.messageChain, &Message{
					Role:    "user",
					Content: "FunctionCall格式错误，请重新输出",
				})
				continue
			}
			funcs = append(funcs, &funccall)
			log.Println("正在调用" + funccall.Action)
			callResult := functioncall.CallFunction(&funccall)
			if callResult == "" {
				callResult = "调用" + funccall.Action + "失败"
			}
			zp.messageChain = append(zp.messageChain, &Message{
				Role:    "user",
				Content: callResult,
			})
			continue
		} else {
			return &AiAnswer{
				Response:       resp.text,
				IsFunctionCall: isFunccall,
				FunctionCall:   funcs,
			}
		}
	}
}
