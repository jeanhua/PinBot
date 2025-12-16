package aibot

import (
	"encoding/json"
	"testing"

	"github.com/jeanhua/PinBot/testoutput"
)

func TestStruct(t *testing.T) {
	funcs := initFunctionTools()
	jsonResult, err := json.Marshal(funcs)
	if err != nil {
		panic(err)
	}
	testoutput.Output(string(jsonResult))
}
