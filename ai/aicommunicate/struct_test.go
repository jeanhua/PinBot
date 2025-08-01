package aicommunicate

import (
	"encoding/json"
	"testing"

	"github.com/jeanhua/PinBot/testoutput"
)

func TestStruct(t *testing.T) {
	funcs := initFunctionTools()
	jsonResult, _ := json.Marshal(funcs)
	testoutput.Output(string(jsonResult))
}
