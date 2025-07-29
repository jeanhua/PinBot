package aicommunicate

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
	return &functionCallTool{
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

type functionCall []*functionCallTool

func (funcs *functionCall) AddFunction(tool *functionCallTool) {
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
