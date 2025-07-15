package aicommunicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jeanhua/PinBot/ai/functioncall"
)

type DeepSeekAIBot_v3 struct {
	token        string
	SystemPrompt string
	messageChain []*Message
	tools        []*FunctionCallTool
	sendVoidce   func(text string)
}

const requestUrl = "https://api.siliconflow.cn/v1/chat/completions"

func NewDeepSeekV3(prompt, token string, sendVoidce func(text string)) *DeepSeekAIBot_v3 {
	tools := FunctionCall{}
	tools.AddFunction(MakeFunctionCallTools(
		"webSearch",
		"执行网络搜索，用于获取互联网相关信息",
		WithParams("query", "搜索查询内容", "string", true),
		WithParams("timeRange", "限制搜索结果的时间范围(可选)(day,week,month,year)", "string", false),
		WithParams("include", "限定搜索结果必须包含的域名列表(可选)", "array<string>", false),
		WithParams("exclude", "排除特定域名的搜索结果(可选)", "array<string>", false),
		WithParams("count", "返回的最大搜索结果数量(可选)", "int", false),
	))
	tools.AddFunction(MakeFunctionCallTools(
		"webExplore",
		"打开某些网页链接进行网页浏览",
		WithParams("links", "要打开的网页链接数组", "array<string>", true),
	))
	tools.AddFunction(MakeFunctionCallTools("browseHomepage", "浏览校园集市论坛主页帖子", WithParams("fromTime", "时间戳,该时间戳前的10条帖子,输入0则表示最新的10条帖子,通过获取帖子后再输入最后帖子的时间戳来实现翻页", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))
	tools.AddFunction(MakeFunctionCallTools("searchPost", "搜索校园集市论坛帖子", WithParams("keywords", "搜索关键词", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("viewComments", "浏览校园集市论坛指定帖子的评论", WithParams("postId", "帖子ID", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("viewPost", "查看校园集市论坛某个帖子详情", WithParams("postId", "帖子ID", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("speak", "向用户发送一段不超过60s的语音", WithParams("text", "要发送的语音内容", "string", true)))

	return &DeepSeekAIBot_v3{
		messageChain: []*Message{
			{
				Role:    "system",
				Content: prompt,
			},
		},
		tools:        tools,
		token:        token,
		SystemPrompt: prompt,
		sendVoidce:   sendVoidce,
	}
}

func request(msg []*Message, model, token string, tools []*FunctionCallTool) (*CommonResponseBody, error) {
	debug := false
	body := &CommonRequestBody{
		Model:    model,
		Messages: msg,
		Stream:   false,
		Tools:    tools,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	if debug {
		log.Println(string(bodyBytes))
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, bodyBytes)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	if debug {
		log.Println(string(respBytes))
	}
	result := &CommonResponseBody{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return result, nil
}

func (deepseek *DeepSeekAIBot_v3) Ask(question string) []*AiAnswer {
	if strings.Contains(question, "#新对话") {
		deepseek.messageChain = []*Message{
			{
				Role:    "system",
				Content: deepseek.SystemPrompt,
			},
		}
	}

	deepseek.messageChain = append(deepseek.messageChain, &Message{
		Role:    "user",
		Content: question,
	})

	var back []*AiAnswer

	for {
		answer, err := request(
			deepseek.messageChain,
			"deepseek-ai/DeepSeek-V3",
			deepseek.token,
			deepseek.tools,
		)
		if err != nil || len(answer.Choices) == 0 {
			log.Println("Request failed:", err)
			return nil
		}

		choice := answer.Choices[0]
		toolCalls := choice.Message.ToolCalls
		if len(toolCalls) != 0 {
			if choice.Message.Content != "" {
				deepseek.messageChain = append(deepseek.messageChain,
					&Message{
						Role:    "assistant",
						Content: choice.Message.Content,
					},
				)
				back = append(back, &AiAnswer{
					Response: choice.Message.Content,
				})
			}
			for _, toolCall := range toolCalls {
				var paramMap map[string]any
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &paramMap); err != nil {
					log.Println("Unmarshal tool call args failed:", err)
					return nil
				}
				callResult, err := functioncall.CallFunction(toolCall.Function.Name, paramMap, deepseek.sendVoidce)
				if err != nil {
					log.Println("CallFunction failed:", err)
					return nil
				}
				deepseek.messageChain = append(deepseek.messageChain,
					&Message{
						Role:       "tool",
						Content:    callResult,
						ToolCallId: toolCall.Id,
					},
				)
			}
			continue
		}
		deepseek.messageChain = append(deepseek.messageChain,
			&Message{
				Role:    "assistant",
				Content: choice.Message.Content,
			},
		)
		back = append(back, &AiAnswer{
			Response: choice.Message.Content,
		})
		break
	}
	return back
}
