package aicommunicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jeanhua/PinBot/config"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jeanhua/PinBot/ai/functioncall"
	"github.com/jeanhua/PinBot/botcontext"
	"github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
)

type AiBot struct {
	token                  string
	SystemPrompt           string
	messageChain           []*message
	tools                  []*functionCallTool
	theLastCommunicateTime time.Time
	target                 int
	mutex                  sync.Mutex
}

// NewAiBot 创建新的AI bot实例
func NewAiBot(prompt, token string, target int) *AiBot {
	return &AiBot{
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
		withParams("count", "返回的最大搜索结果数量(可选)", "number", false),
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

	tools.addFunction(makeFunctionCallTools("hateImage", "发送讨厌表情包(表达生气，表情包内容为动漫卡通指着对方头像说：爬！)", withParams("userid", "用户的Id", "string", true)))

	// 歌曲相关
	tools.addFunction(makeFunctionCallTools("searchMusic", "搜索音乐", withParams("query", "关键词,歌曲名称或者歌手名字", "string", true)))
	tools.addFunction(makeFunctionCallTools("shareMusic", "分享音乐", withParams("id", "音乐id,搜索到的音乐id", "string", true)))

	// 第二课堂相关
	tools.addFunction(makeFunctionCallTools("scu2ClassSearch", "检索第二课堂系列活动", withParams("keyword", "关键词,活动名称的关键词", "string", true)))
	tools.addFunction(makeFunctionCallTools("scu2ClassList", "通过系列活动ID获取具体活动", withParams("activityLibId", "系列活动ID", "string", true)))
	tools.addFunction(makeFunctionCallTools("scu2ClassShare", "发送具体活动的签到签退二维码", withParams("activityId", "活动ID", "string", true)))

	return tools
}

func (aiBot *AiBot) SendMsg(msg string, group_msg *model.GroupMessage, friend_msg *model.FriendMessage) {
	mutMsg := []rune(msg)
	if aiBot.target == functioncall.TargetFriend {
		if len(mutMsg) <= 500 {
			chain := messagechain.Friend(friend_msg.Sender.UserId)
			chain.Reply(friend_msg.MessageId)
			chain.Text(msg)
			chain.Send()
		} else {
			for i := 0; i < len(mutMsg); i += 500 {
				end := i + 500
				if end > len(mutMsg) {
					end = len(mutMsg)
				}
				segment := string(mutMsg[i:end])

				chain := messagechain.Friend(friend_msg.Sender.UserId)
				if i == 0 {
					chain.Reply(friend_msg.MessageId)
				}
				chain.Text(segment)
				chain.Send()
			}
		}
	} else {
		if len(mutMsg) <= 500 {
			chain := messagechain.Group(group_msg.GroupId)
			chain.Reply(group_msg.MessageId)
			chain.Mention(group_msg.Sender.UserId)
			chain.Text(" " + msg)
			chain.Send()
		} else {
			botcontext.SendLongReply(group_msg, mutMsg)
		}
	}
}

// Ask 处理用户提问并返回AI的回答
func (aiBot *AiBot) Ask(question string, group_msg *model.GroupMessage, friend_msg *model.FriendMessage) {

	ok := aiBot.mutex.TryLock()
	if !ok {
		aiBot.SendMsg("等待上一个请求完成，不要着急哦!", group_msg, friend_msg)
		return
	}
	defer aiBot.mutex.Unlock()

	aiBot.checkNeedReset(question)
	// 添加用户消息到对话链
	aiBot.appendUserMessage(question)

	for {
		// 发送请求获取AI响应
		answer, err := request(
			aiBot.messageChain,
			config.GetConfig().AiModelName,
			aiBot.token,
			aiBot.tools,
		)

		if err != nil || len(answer.Choices) == 0 {
			log.Println("Request failed:", err)
			aiBot.SendMsg("抱歉，我遇到了一些问题，请稍后再试。", group_msg, friend_msg)
			return
		}

		choice := answer.Choices[0]

		// 处理工具调用
		if len(choice.Message.ToolCalls) > 0 {
			aiBot.handleToolCalls(&choice, group_msg, friend_msg)
			continue
		}

		// 处理普通响应
		aiBot.SendMsg(choice.Message.Content, group_msg, friend_msg)
		aiBot.appendMessage(&choice.Message)
		break
	}
}

// checkNeedReset 检查是否需要重置对话
func (aiBot *AiBot) checkNeedReset(question string) {
	// 三小时自动重置对话
	if aiBot.theLastCommunicateTime.Add(time.Hour * 3).Before(time.Now()) {
		aiBot.resetConversation()
	}
	aiBot.theLastCommunicateTime = time.Now()
	if strings.Contains(question, "#新对话") {
		aiBot.resetConversation()
	} else {
		aiBot.autoNewCommunication()
	}
}

// autoNewCommunication 自动新对话
func (aiBot *AiBot) autoNewCommunication() {
	if len(aiBot.messageChain) >= 120 {
		aiBot.resetConversation()
	}
}

// resetConversation 重置对话历史
func (aiBot *AiBot) resetConversation() {
	aiBot.messageChain = []*message{
		{
			Role:    "system",
			Content: aiBot.SystemPrompt,
		},
	}
}

// appendUserMessage 添加用户消息到对话链
func (aiBot *AiBot) appendUserMessage(content string) {
	aiBot.messageChain = append(aiBot.messageChain, &message{
		Role:    "user",
		Content: content,
	})
}

// appendMessage 添加消息
func (aiBot *AiBot) appendMessage(msg *message) {
	aiBot.messageChain = append(aiBot.messageChain, msg)
}

// handleToolCalls 处理工具调用
func (aiBot *AiBot) handleToolCalls(choice *choice, group_msg *model.GroupMessage, friend_msg *model.FriendMessage) {
	aiBot.appendMessage(&choice.Message)
	if choice.Message.Content != "" {
		aiBot.SendMsg(choice.Message.Content, group_msg, friend_msg)
	}
	// 处理工具调用
	if err := aiBot.executeToolCalls(choice.Message.ToolCalls, group_msg, friend_msg); err != nil {
		log.Println("Error executing tool calls:", err)
		aiBot.SendMsg("function call 调用出现问题", group_msg, friend_msg)
		return
	}
}

// executeToolCalls 执行工具调用
func (aiBot *AiBot) executeToolCalls(toolCalls []toolCall, group_msg *model.GroupMessage, friend_msg *model.FriendMessage) error {
	for _, toolCall := range toolCalls {
		var paramMap map[string]any
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &paramMap); err != nil {
			log.Println("Unmarshal tool call args failed:", err)
			log.Println(toolCall.Function.Arguments)
			return err
		}
		var uid uint
		if aiBot.target == functioncall.TargetFriend {
			uid = friend_msg.Sender.UserId
		} else {
			uid = group_msg.GroupId
		}
		callResult, err := functioncall.CallFunction(toolCall.Function.Name, paramMap, uid, aiBot.target)
		if err != nil {
			log.Println("CallFunction failed:", err)
			callResult = "function call 调用失败"
		}
		aiBot.appendToolResponse(toolCall.Id, callResult)
	}
	return nil
}

// appendToolResponse 添加工具响应到对话链
func (aiBot *AiBot) appendToolResponse(toolCallId, content string) {
	aiBot.messageChain = append(aiBot.messageChain, &message{
		Role:       "tool",
		Content:    content,
		ToolCallId: toolCallId,
	})
}

// request 发送请求到AI接口
func request(msg []*message, model, token string, tools []*functionCallTool) (*commonResponseBody, error) {
	body := &commonRequestBody{
		Model:       model,
		Messages:    msg,
		Stream:      false,
		Tools:       tools,
		Temperature: config.GetConfig().AiTemperature,
		TopK:        config.GetConfig().AiTopK,
		TopP:        config.GetConfig().AiTopP,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, config.GetConfig().AiRequestUrl, bytes.NewBuffer(bodyBytes))
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

	result := &commonResponseBody{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return result, nil
}
