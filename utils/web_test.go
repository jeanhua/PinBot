package utils

import (
	"fmt"
	"testing"
)

func TestWebSearch(t *testing.T) {
	token := "xxx"
	resp := WebSearch(token, "jeanhua", nil, []string{}, []string{}, 8)
	fmt.Println(resp)
}

func TestWebExplore(t *testing.T) {
	token := "xxx"
	resp := WebExplore([]string{"https://www.blog.jeanhua.cn"}, token)
	fmt.Println(resp)
}
