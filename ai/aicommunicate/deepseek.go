package aicommunicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	tools.AddFunction(MakeFunctionCallTools("browseHomepage", "浏览校园集市论坛主页", WithParams("fromTime", "时间戳,该时间戳前的10条帖子,输入0则表示最新的10条帖子,通过获取帖子后再输入最后帖子的时间戳来实现翻页", "string")))
	tools.AddFunction(MakeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))
	tools.AddFunction(MakeFunctionCallTools("searchPost", "搜索校园集市论坛帖子", WithParams("keywords", "搜索关键词", "string")))
	tools.AddFunction(MakeFunctionCallTools("viewComments", "浏览校园集市论坛指定帖子的评论", WithParams("postId", "帖子ID", "string")))
	tools.AddFunction(MakeFunctionCallTools("viewPost", "调用这个工具可以向用户发送一段不超过60s的语音，偶尔可以调用玩一下", WithParams("postId", "帖子ID", "string")))
	tools.AddFunction(MakeFunctionCallTools("speak", "调用这个工具可以向用户发送一段不超过60s的语音，偶尔可以调用玩一下", WithParams("text", "要发送的文本内容", "string")))
	tools.AddFunction(MakeFunctionCallTools("webSearch", "执行网页搜索获取结果",
		WithParams("query", "搜索关键词", "string"),
		WithParams("token", "身份验证Token（可选）", "string"),
		WithParams("freshness", "限制搜索结果的新鲜程度（可选），例如 - noLimit，不限（默认）oneDay，一天内 oneWeek，一周内 oneMonth，一个月内 oneYear，一年内 YYYY-MM-DD..YYYY-MM-DD，搜索日期范围，例如：\"2025-01-01..2025-04-06\" YYYY-MM-DD，搜索指定日期，例如：\"2025-04-06\"", "string"),
		WithParams("summary", "是否生成摘要信息（可选）", "bool"),
		WithParams("include", "指定包含的网站或域名（可选），多个域名使用|或,分隔，最多不能超过20个", "string"),
		WithParams("exclude", "指定排除的网站或域名（可选），排除搜索的网站范围。多个域名使用|或,分隔，最多不能超过20个", "string"),
		WithParams("count", "搜索结果数量，默认为5（可选）", "int")))
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
		fmt.Println(string(bodyBytes))
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
		fmt.Println(string(respBytes))
	}
	result := &CommonResponseBody{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return result, nil
}

func (deepseek *DeepSeekAIBot_v3) Ask(question string) *AiAnswer {
	if strings.Contains(question, "#新对话") {
		deepseek.messageChain = []*Message{
			{
				Role:    "system",
				Content: deepseek.SystemPrompt,
			},
		}
	}

	var finalAnswer *AiAnswer

	deepseek.messageChain = append(deepseek.messageChain, &Message{
		Role:    "user",
		Content: question,
	})

	for {
		answer, err := request(
			deepseek.messageChain,
			"deepseek-ai/DeepSeek-V3",
			deepseek.token,
			deepseek.tools,
		)
		if err != nil || len(answer.Choices) == 0 {
			fmt.Println("Request failed:", err)
			return nil
		}

		choice := answer.Choices[0]
		toolCalls := choice.Message.ToolCalls
		if choice.FinishReason == "tool_calls" {
			toolCall := toolCalls[0]
			var paramMap map[string]any
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &paramMap); err != nil {
				fmt.Println("Unmarshal tool call args failed:", err)
				return nil
			}
			callResult, err := functioncall.CallFunction(toolCall.Function.Name, paramMap, deepseek.sendVoidce)
			if err != nil {
				fmt.Println("CallFunction failed:", err)
				return nil
			}
			deepseek.messageChain = append(deepseek.messageChain,
				&Message{
					Role:       "tool",
					Content:    callResult,
					ToolCallId: toolCall.Id,
				},
			)
			continue
		}
		finalAnswer = &AiAnswer{
			Response: choice.Message.Content,
		}
		deepseek.messageChain = append(deepseek.messageChain,
			&Message{Role: "assistant", Content: finalAnswer.Response},
		)
		break
	}
	return finalAnswer
}
