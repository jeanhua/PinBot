package aicommunicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/ai/functioncall"
)

const (
	requestUrl    = "https://api.siliconflow.cn/v1/chat/completions"
	deepSeekModel = "deepseek-ai/DeepSeek-V3"
)

type DeepSeekAIBot_v3 struct {
	token                  string
	SystemPrompt           string
	messageChain           []*Message
	tools                  []*FunctionCallTool
	sendVoidce             func(text string)
	theLastCommunicateTime time.Time
}

// 创建新的DeepSeek AI V3实例
func NewDeepSeekV3(prompt, token string, sendVoidce func(text string)) *DeepSeekAIBot_v3 {
	return &DeepSeekAIBot_v3{
		messageChain: []*Message{
			{
				Role:    "system",
				Content: prompt,
			},
		},
		tools:                  initFunctionTools(),
		token:                  token,
		SystemPrompt:           prompt,
		sendVoidce:             sendVoidce,
		theLastCommunicateTime: time.Now(),
	}
}

// 初始化所有可用的功能工具
func initFunctionTools() []*FunctionCallTool {
	tools := FunctionCall{}
	const strArry = "array<string>"

	// 定义所有工具函数
	tools.AddFunction(MakeFunctionCallTools(
		"webSearch",
		"执行网络搜索，用于获取互联网相关信息",
		WithParams("query", "搜索查询内容", "string", true),
		WithParams("timeRange", "限制搜索结果的时间范围(可选)(day,week,month,year)", "string", false),
		WithParams("include", "限定搜索结果必须包含的域名列表(可选)", strArry, false),
		WithParams("exclude", "排除特定域名的搜索结果(可选)", strArry, false),
		WithParams("count", "返回的最大搜索结果数量(可选)", "int", false),
	))

	tools.AddFunction(MakeFunctionCallTools(
		"webExplore",
		"打开某些网页链接进行网页浏览",
		WithParams("links", "要打开的网页链接数组", strArry, true),
	))

	tools.AddFunction(MakeFunctionCallTools(
		"browseHomepage",
		"浏览校园集市论坛主页帖子",
		WithParams("fromTime", "时间戳,该时间戳前的10条帖子,输入0则表示最新的10条帖子", "string", true),
	))

	tools.AddFunction(MakeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))

	tools.AddFunction(MakeFunctionCallTools(
		"searchPost",
		"搜索校园集市论坛帖子",
		WithParams("keywords", "搜索关键词", "string", true),
	))

	tools.AddFunction(MakeFunctionCallTools(
		"viewComments",
		"浏览校园集市论坛指定帖子的评论",
		WithParams("postId", "帖子ID", "string", true),
	))

	tools.AddFunction(MakeFunctionCallTools(
		"viewPost",
		"查看校园集市论坛某个帖子详情",
		WithParams("postId", "帖子ID", "string", true),
	))

	tools.AddFunction(MakeFunctionCallTools(
		"speak",
		"向用户发送一段不超过60s的语音",
		WithParams("text", "要发送的语音内容", "string", true),
	))

	tools.AddFunction(MakeFunctionCallTools("getCurrentTime", "获取当前时间"))

	return tools
}

// Ask 处理用户提问并返回AI的回答
func (deepseek *DeepSeekAIBot_v3) Ask(question string) []*AiAnswer {
	if deepseek.theLastCommunicateTime.Add(time.Hour * 3).Before(time.Now()) {
		deepseek.resetConversation()
	}
	deepseek.theLastCommunicateTime = time.Now()
	// 检查是否需要重置对话
	if strings.Contains(question, "#新对话") {
		deepseek.resetConversation()
	} else {
		deepseek.autoNewCommunication()
	}

	// 添加用户消息到对话链
	deepseek.appendUserMessage(question)

	var responses []*AiAnswer

	for {
		// 发送请求获取AI响应
		answer, err := request(
			deepseek.messageChain,
			deepSeekModel,
			deepseek.token,
			deepseek.tools,
		)

		if err != nil || len(answer.Choices) == 0 {
			log.Println("Request failed:", err)
			return nil
		}

		choice := answer.Choices[0]

		// 处理工具调用
		if len(choice.Message.ToolCalls) > 0 {
			responses = deepseek.handleToolCalls(&choice, responses)
			continue
		}

		// 处理普通响应
		responses = append(responses, &AiAnswer{
			Response: choice.Message.Content,
		})
		deepseek.appendAssistantMessage(choice.Message.Content)
		break
	}

	return responses
}

// 自动新对话
func (deepseek *DeepSeekAIBot_v3) autoNewCommunication() {
	if len(deepseek.messageChain) >= 120 {
		deepseek.resetConversation()
	}
}

// 重置对话历史
func (deepseek *DeepSeekAIBot_v3) resetConversation() {
	deepseek.messageChain = []*Message{
		{
			Role:    "system",
			Content: deepseek.SystemPrompt,
		},
	}
}

// 添加用户消息到对话链
func (deepseek *DeepSeekAIBot_v3) appendUserMessage(content string) {
	deepseek.messageChain = append(deepseek.messageChain, &Message{
		Role:    "user",
		Content: content,
	})
}

// 添加助手消息到对话链
func (deepseek *DeepSeekAIBot_v3) appendAssistantMessage(content string) {
	deepseek.messageChain = append(deepseek.messageChain, &Message{
		Role:    "assistant",
		Content: content,
	})
}

// 处理工具调用
func (deepseek *DeepSeekAIBot_v3) handleToolCalls(choice *Choice, responses []*AiAnswer) []*AiAnswer {
	// 如果有内容先添加内容响应
	if strings.TrimSpace(choice.Message.Content) != "" {
		responses = append(responses, &AiAnswer{
			Response: choice.Message.Content,
		})
		deepseek.appendAssistantMessage(choice.Message.Content)
	}

	// 处理工具调用
	if err := deepseek.executeToolCalls(choice.Message.ToolCalls); err != nil {
		log.Println("Error executing tool calls:", err)
		return nil
	}

	return responses
}

// 执行工具调用
func (deepseek *DeepSeekAIBot_v3) executeToolCalls(toolCalls []ToolCall) error {
	for _, toolCall := range toolCalls {
		var paramMap map[string]any
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &paramMap); err != nil {
			log.Println("Unmarshal tool call args failed:", err)
			return err
		}

		callResult, err := functioncall.CallFunction(toolCall.Function.Name, paramMap, deepseek.sendVoidce)
		if err != nil {
			log.Println("CallFunction failed:", err)
			return err
		}

		deepseek.appendToolResponse(toolCall.Id, callResult)
	}
	return nil
}

// 添加工具响应到对话链
func (deepseek *DeepSeekAIBot_v3) appendToolResponse(toolCallId, content string) {
	deepseek.messageChain = append(deepseek.messageChain, &Message{
		Role:       "tool",
		Content:    content,
		ToolCallId: toolCallId,
	})
}

// 发送请求到AI接口
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
		log.Println("Request body:", string(bodyBytes))
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
		log.Println("Response body:", string(respBytes))
	}

	result := &CommonResponseBody{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return result, nil
}
