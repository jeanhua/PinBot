package aicommunicate

import (
	"github.com/jeanhua/PinBot/ai/functioncall"
)

type AiModel interface {
	Ask(question string) *AiAnswer
}

type AiAnswer struct {
	Response       string                      `json:"response"`
	IsFunctionCall bool                        `json:"isFunction_call"`
	FunctionCall   []functioncall.FunctionCall `json:"function_call"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
