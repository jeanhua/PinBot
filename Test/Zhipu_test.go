package test

import (
	"fmt"
	"testing"

	llm "github.com/jeanhua/PinBot/LLM"
)

func TestZhipu(t *testing.T) {
	zhipu := llm.NewZhiPu()
	resp, err := zhipu.RequestReply(1213, "你好", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
