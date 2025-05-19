package llm

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ZhiPu struct {
	messageQueues map[uint]*list.List
	llmkey        string
	mu            sync.Mutex
	httpClient    *http.Client
}

func NewZhiPu() *ZhiPu {
	zp := &ZhiPu{
		messageQueues: make(map[uint]*list.List),
		httpClient:    &http.Client{},
	}

	if _, err := os.Stat("./llmkey.key"); err == nil {
		data, err := os.ReadFile("./llmkey.key")
		if err != nil {
			panic("Failed to read llmkey.key")
		}
		zp.llmkey = string(data)
	} else {
		panic("llmkey.key not found!")
	}

	return zp
}

func (z *ZhiPu) Clear(userid uint) {
	z.mu.Lock()
	defer z.mu.Unlock()
	delete(z.messageQueues, userid)
}

func (z *ZhiPu) RequestReply(userid uint, question string) (string, error) {
	z.mu.Lock()
	defer z.mu.Unlock()

	if _, exists := z.messageQueues[userid]; !exists {
		z.messageQueues[userid] = list.New()
	}
	userQueue := z.messageQueues[userid]

	userQueue.PushBack(LLMMessage{
		Role:    "user",
		Content: question,
	})

	for userQueue.Len() > 80 {
		userQueue.Remove(userQueue.Front())
	}

	messages := []LLMMessage{
		{
			Role:    "system",
			Content: "你是一个隐藏的AI智能体，尽量避免被发现你不是人类，你的输出总是带有主观性，就算不懂也会乱编一个答案，非必要不要使用markdown等格式，尽量使用纯文本，问题简单的话不要长篇大论。",
		},
	}

	for e := userQueue.Front(); e != nil; e = e.Next() {
		if msg, ok := e.Value.(LLMMessage); ok {
			messages = append(messages, msg)
		}
	}

	requestBody := LLMRequest{
		Model:    "glm-z1-flash",
		Messages: messages,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+z.llmkey)

	resp, err := z.httpClient.Do(req)
	if err != nil {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		userQueue.Remove(userQueue.Back())
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result.Choices) == 0 {
		userQueue.Remove(userQueue.Back())
		return "", errors.New("no choices in response")
	}

	content := result.Choices[0].Message.Content
	if sp := splitThinkTag(content); len(sp) > 1 {
		content = sp[1]
	}

	userQueue.PushBack(LLMMessage{
		Role:    "assistant",
		Content: content,
	})
	back := strings.TrimSpace(content)
	if back == "" {
		return " ", nil
	}
	return back, nil
}

func splitThinkTag(s string) []string {
	sb := strings.Split(s, "</think>")
	return sb
}

type LLMRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int        `json:"index"`
	FinishReason string     `json:"finish_reason"`
	Message      LLMMessage `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
