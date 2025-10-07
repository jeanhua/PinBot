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

func TestMusic(t *testing.T) {
	musicSearch := MusicSearch{}
	result := musicSearch.Search("鹿晗")
	log.Println(result)
}

func TestSCU2Class_Search(t *testing.T) {
	scu2class := NewSCU2Class("xxx")
	c_in, c_out, err := scu2class.GenQRCode("231")
	if err != nil {
		log.Println(err)
	}
	log.Println(c_in)
	log.Println(c_out)
}
