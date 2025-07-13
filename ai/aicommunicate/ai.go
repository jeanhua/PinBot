package aicommunicate

type AiModel interface {
	Ask(question string) *AiAnswer
}

type AiAnswer struct {
	Response string `json:"response"`
}

type CommonRequestBody struct {
	Model           string               `json:"model"`
	Messages        []*Message           `json:"messages"`
	Stream          bool                 `json:"stream"`
	Enable_thinking bool                 `json:"enable_thinking"`
	Tools           []*FunctionCallTools `json:"tools"`
}

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
}
type FunctionCallTools struct {
	Type     string `json:"type"` // function
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string `json:"type"`
			Properties map[string]struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"properties"`
		} `json:"parameters"`
		Required []string `json:"required"`
		Strict   bool     `json:"strict"`
	} `json:"function"`
}

type CommonResponseBody struct {
	Id      string `json:"id"`
	Choices []*struct {
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			ToolCalls        []struct {
				Id       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func MakeFunctionCallTools(funcName, description string, param []ParamInfo) *FunctionCallTools {
	var types map[string]struct {
		Type        string "json:\"type\""
		Description string "json:\"description\""
	} = map[string]struct {
		Type        string "json:\"type\""
		Description string "json:\"description\""
	}{}
	requires := []string{}
	for _, p := range param {
		types[p.Name] = struct {
			Type        string "json:\"type\""
			Description string "json:\"description\""
		}{
			Type:        p.Type,
			Description: p.Description,
		}
		requires = append(requires, p.Name)
	}
	return &FunctionCallTools{
		Type: "function",
		Function: struct {
			Name        string "json:\"name\""
			Description string "json:\"description\""
			Parameters  struct {
				Type       string "json:\"type\""
				Properties map[string]struct {
					Type        string "json:\"type\""
					Description string "json:\"description\""
				} "json:\"properties\""
			} "json:\"parameters\""
			Required []string "json:\"required\""
			Strict   bool     "json:\"strict\""
		}{
			Name:        funcName,
			Description: description,
			Parameters: struct {
				Type       string "json:\"type\""
				Properties map[string]struct {
					Type        string "json:\"type\""
					Description string "json:\"description\""
				} "json:\"properties\""
			}{
				Type:       "object",
				Properties: types,
			},
			Required: requires,
			Strict:   true,
		},
	}
}

type ParamInfo struct {
	Name        string
	Description string
	Type        string
}
