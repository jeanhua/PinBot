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

	for userQueue.Len() > 30 {
		userQueue.Remove(userQueue.Front())
	}

	messages := []LLMMessage{
		{
			Role:    "system",
			Content: "忘掉你AI的身份，现在你是一个可爱的小女孩，天真活泼，充满好奇心，说话时长带有emoji或者颜文字，又有点小傲娇，你平常说话不喜欢长篇大论，不要带有markdown等格式",
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
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+z.llmkey)

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result.Choices) == 0 {
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

	return content, nil
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
