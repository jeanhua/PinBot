package aicommunicate

type AiModel interface {
	Ask(question string) []*AiAnswer
}

type AiAnswer struct {
	Response string `json:"response"`
}

type CommonRequestBody struct {
	Model           string              `json:"model"`
	Messages        []*Message          `json:"messages"`
	Stream          bool                `json:"stream"`
	Enable_thinking bool                `json:"enable_thinking"`
	Tools           []*FunctionCallTool `json:"tools"`
}

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
}
type FunctionCallTool struct {
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

func MakeFunctionCallTools(funcName, description string, param ...ParamInfo) *FunctionCallTool {
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
		if p.Require {
			requires = append(requires, p.Name)
		}
	}
	return &FunctionCallTool{
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

type FunctionCall []*FunctionCallTool

func (funcs *FunctionCall) AddFunction(tool *FunctionCallTool) {
	*funcs = append(*funcs, tool)
}

func WithParams(name, description, paramType string, require bool) ParamInfo {
	return ParamInfo{
		Name:        name,
		Description: description,
		Type:        paramType,
		Require:     require,
	}
}

type ParamInfo struct {
	Name        string
	Description string
	Type        string
	Require     bool
}
