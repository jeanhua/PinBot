package utils

import (
	"log"
	"testing"
)

func TestWebSearch(t *testing.T) {
	token := "xxx"
	resp := WebSearch(token, "jeanhua", nil, []string{}, []string{}, 8)
	log.Println(resp)
}

func TestWebExplore(t *testing.T) {
	token := "xxx"
	resp := WebExplore([]string{"https://www.blog.jeanhua.cn"}, token)
	log.Println(resp)
}
