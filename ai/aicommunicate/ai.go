package aicommunicate

import (
	"strings"

	"github.com/jeanhua/PinBot/model"
)

type AiModel interface {
	Ask(question string, group_msg *model.GroupMessage, friend_msg *model.FriendMessage)
}

type commonRequestBody struct {
	Model           string              `json:"model"`
	Messages        []*message          `json:"messages"`
	Stream          bool                `json:"stream"`
	Enable_thinking bool                `json:"enable_thinking"`
	Tools           []*functionCallTool `json:"tools"`
	Temperature     float32             `json:"temperature"`
	TopK            int                 `json:"top_k"`
	TopP            float32             `json:"top_p"`
}

type message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []toolCall `json:"tool_calls"`
	ToolCallId string     `json:"tool_call_id"`
}
type functionCallTool struct {
	Type     string `json:"type"` // function
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string          `json:"type"`
			Properties *map[string]any `json:"properties"`
		} `json:"parameters"`
		Required []string `json:"required"`
		Strict   bool     `json:"strict"`
	} `json:"function"`
}

type commonResponseBody struct {
	Id      string   `json:"id"`
	Choices []choice `json:"choices"`
}

type choice struct {
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type toolCall struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func makeFunctionCallTools(funcName, description string, param ...paramInfo) *functionCallTool {
	var propoties = map[string]any{}
	requires := []string{}
	const arrayPrefex = "array:"
	for _, p := range param {
		if strings.HasPrefix(p.Type, arrayPrefex) {
			propoties[p.Name] = &map[string]any{
				"type": "array",
				"items": map[string]string{
					"type": strings.TrimPrefix(p.Type, arrayPrefex),
				},
				"description": p.Description,
			}
		} else {
			propoties[p.Name] = &map[string]any{
				"type":        p.Type,
				"description": p.Description,
			}
		}
		if p.Require {
			requires = append(requires, p.Name)
		}
	}
	return &functionCallTool{
		Type: "function",
		Function: struct {
			Name        string "json:\"name\""
			Description string "json:\"description\""
			Parameters  struct {
				Type       string          "json:\"type\""
				Properties *map[string]any "json:\"properties\""
			} "json:\"parameters\""
			Required []string "json:\"required\""
			Strict   bool     "json:\"strict\""
		}{
			Name:        funcName,
			Description: description,
			Parameters: struct {
				Type       string          "json:\"type\""
				Properties *map[string]any "json:\"properties\""
			}{
				Type:       "object",
				Properties: &propoties,
			},
			Required: requires,
			Strict:   true,
		},
	}
}

type functionCall []*functionCallTool

func (funcs *functionCall) addFunction(tool *functionCallTool) {
	*funcs = append(*funcs, tool)
}

func withParams(name, description, paramType string, require bool) paramInfo {
	return paramInfo{
		Name:        name,
		Description: description,
		Type:        paramType,
		Require:     require,
	}
}

type paramInfo struct {
	Name        string
	Description string
	Type        string
	Require     bool
}
