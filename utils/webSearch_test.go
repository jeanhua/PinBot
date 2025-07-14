package utils

import (
	"fmt"
	"testing"
)

func TestWebSearch(t *testing.T) {
	token := "xxx"
	resp := WebSearch(token, "jeanhua", "noLimit", true, "", "", 10)
	fmt.Println(resp)
}
