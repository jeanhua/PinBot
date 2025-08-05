package aicommunicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
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
	messageChain           []*message
	tools                  []*functionCallTool
	theLastCommunicateTime time.Time
	target                 int
	uid                    uint
	mutex                  sync.Mutex
}

// 创建新的DeepSeek AI V3实例
func NewDeepSeekV3(prompt, token string, target int, uid uint) *DeepSeekAIBot_v3 {
	return &DeepSeekAIBot_v3{
		messageChain: []*message{
			{
				Role:    "system",
				Content: prompt,
			},
		},
		tools:                  initFunctionTools(),
		token:                  token,
		SystemPrompt:           prompt,
		theLastCommunicateTime: time.Now(),
		target:                 target,
		uid:                    uid,
		mutex:                  sync.Mutex{},
	}
}

// 初始化所有可用的功能工具
func initFunctionTools() []*functionCallTool {
	tools := functionCall{}
	const strArray = "array:string"
	// 定义所有工具函数
	tools.addFunction(makeFunctionCallTools(
		"webSearch",
		"执行网络搜索，用于获取互联网相关信息",
		withParams("query", "搜索查询内容", "string", true),
		withParams("timeRange", "限制搜索结果的时间范围(可选)(day,week,month,year)", "string", false),
		withParams("include", "限定搜索结果必须包含的域名列表(可选)", strArray, false),
		withParams("exclude", "排除特定域名的搜索结果(可选)", strArray, false),
		withParams("count", "返回的最大搜索结果数量(可选)", "int", false),
	))

	tools.addFunction(makeFunctionCallTools(
		"webExplore",
		"打开某些网页链接进行网页浏览",
		withParams("links", "要打开的网页链接数组", strArray, true),
	))

	tools.addFunction(makeFunctionCallTools(
		"browseHomepage",
		"浏览校园集市论坛主页帖子",
		withParams("fromTime", "时间戳,该时间戳前的10条帖子,输入0则表示最新的10条帖子", "string", true),
	))

	tools.addFunction(makeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))

	tools.addFunction(makeFunctionCallTools(
		"searchPost",
		"搜索校园集市论坛帖子",
		withParams("keywords", "搜索关键词", "string", true),
	))

	tools.addFunction(makeFunctionCallTools(
		"viewComments",
		"浏览校园集市论坛指定帖子的评论",
		withParams("postId", "帖子ID", "string", true),
	))

	tools.addFunction(makeFunctionCallTools(
		"viewPost",
		"查看校园集市论坛某个帖子详情",
		withParams("postId", "帖子ID", "string", true),
	))

	tools.addFunction(makeFunctionCallTools(
		"speak",
		"向用户发送一段不超过60s的语音",
		withParams("text", "要发送的语音内容", "string", true),
	))

	tools.addFunction(makeFunctionCallTools("getCurrentTime", "获取当前时间"))

	tools.addFunction(makeFunctionCallTools("hateImage", "发送讨厌表情包(表达生气)", withParams("userid", "用户的Id", "string", true)))

	// 歌曲相关
	tools.addFunction(makeFunctionCallTools("searchMusic", "搜索音乐", withParams("keyword", "关键词,歌曲名称或者歌手名字", "string", true)))
	tools.addFunction(makeFunctionCallTools("shareMusic", "分享音乐", withParams("id", "音乐id,搜索到的音乐id", "string", true)))

	return tools
}

// Ask 处理用户提问并返回AI的回答
func (deepseek *DeepSeekAIBot_v3) Ask(question string) []*AiAnswer {

	ok := deepseek.mutex.TryLock()
	if !ok {
		return []*AiAnswer{
			{
				Response: "等待上一个请求完成，不要着急哦!",
			},
		}
	}
	defer deepseek.mutex.Unlock()

	deepseek.checkNeedReset(question)
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

// 检查是否需要重置对话
func (deepseek *DeepSeekAIBot_v3) checkNeedReset(question string) {
	// 三小时自动重置对话
	if deepseek.theLastCommunicateTime.Add(time.Hour * 3).Before(time.Now()) {
		deepseek.resetConversation()
	}
	deepseek.theLastCommunicateTime = time.Now()
	if strings.Contains(question, "#新对话") {
		deepseek.resetConversation()
	} else {
		deepseek.autoNewCommunication()
	}
}

// 自动新对话
func (deepseek *DeepSeekAIBot_v3) autoNewCommunication() {
	if len(deepseek.messageChain) >= 120 {
		deepseek.resetConversation()
	}
}

// 重置对话历史
func (deepseek *DeepSeekAIBot_v3) resetConversation() {
	deepseek.messageChain = []*message{
		{
			Role:    "system",
			Content: deepseek.SystemPrompt,
		},
	}
}

// 添加用户消息到对话链
func (deepseek *DeepSeekAIBot_v3) appendUserMessage(content string) {
	deepseek.messageChain = append(deepseek.messageChain, &message{
		Role:    "user",
		Content: content,
	})
}

// 添加助手消息到对话链
func (deepseek *DeepSeekAIBot_v3) appendAssistantMessage(content string) {
	deepseek.messageChain = append(deepseek.messageChain, &message{
		Role:    "assistant",
		Content: content,
	})
}

// 处理工具调用
func (deepseek *DeepSeekAIBot_v3) handleToolCalls(choice *choice, responses []*AiAnswer) []*AiAnswer {
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
func (deepseek *DeepSeekAIBot_v3) executeToolCalls(toolCalls []toolCall) error {
	for _, toolCall := range toolCalls {
		var paramMap map[string]any
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &paramMap); err != nil {
			log.Println("Unmarshal tool call args failed:", err)
			log.Println(toolCall.Function.Arguments)
			return err
		}
		callResult, err := functioncall.CallFunction(toolCall.Function.Name, paramMap, deepseek.uid, deepseek.target)
		if err != nil {
			log.Println("CallFunction failed:", err)
			callResult = "function call 调用失败"
		}

		deepseek.appendToolResponse(toolCall.Id, callResult)
	}
	return nil
}

// 添加工具响应到对话链
func (deepseek *DeepSeekAIBot_v3) appendToolResponse(toolCallId, content string) {
	deepseek.messageChain = append(deepseek.messageChain, &message{
		Role:       "tool",
		Content:    content,
		ToolCallId: toolCallId,
	})
}

// 发送请求到AI接口
func request(msg []*message, model, token string, tools []*functionCallTool) (*commonResponseBody, error) {
	debug := false
	body := &commonRequestBody{
		Model:       model,
		Messages:    msg,
		Stream:      false,
		Tools:       tools,
		Temperature: 0.9,
		TopK:        66,
		TopP:        0.8,
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

	result := &commonResponseBody{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return result, nil
}
