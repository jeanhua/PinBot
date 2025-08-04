package aicommunicate

import "strings"

type AiModel interface {
	Ask(question string) []*AiAnswer
}

type AiAnswer struct {
	Response string `json:"response"`
}

type commonRequestBody struct {
	Model           string              `json:"model"`
	Messages        []*message          `json:"messages"`
	Stream          bool                `json:"stream"`
	Enable_thinking bool                `json:"enable_thinking"`
	Tools           []*functionCallTool `json:"tools"`
}

type message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
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
		Required             []string `json:"required"`
		Strict               bool     `json:"strict"`
		AdditionalProperties bool     `json:"additionalProperties"`
	} `json:"function"`
}

type commonResponseBody struct {
	Id      string   `json:"id"`
	Choices []choice `json:"choices"`
}

type choice struct {
	Message struct {
		Role             string     `json:"role"`
		Content          string     `json:"content"`
		ReasoningContent string     `json:"reasoning_content"`
		ToolCalls        []toolCall `json:"tool_calls"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}
type toolCall struct {
	Id       string `json:"id"`
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
			Required             []string "json:\"required\""
			Strict               bool     "json:\"strict\""
			AdditionalProperties bool     `json:"additionalProperties"`
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
