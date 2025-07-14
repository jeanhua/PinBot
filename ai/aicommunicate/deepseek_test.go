package aicommunicate

import (
	"encoding/json"
	"log"
	"testing"
)

func TestDeepseekFuncs(t *testing.T) {
	tools := &FunctionCall{}
	tools.AddFunction(MakeFunctionCallTools("browseHomepage", "浏览校园集市论坛主页"))
	tools.AddFunction(MakeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))
	tools.AddFunction(MakeFunctionCallTools("search", "搜索校园集市论坛帖子", WithParams("keywords", "搜索关键词", "string")))
	tools.AddFunction(MakeFunctionCallTools("viewComments", "浏览校园集市论坛指定帖子的评论", WithParams("postId", "帖子ID", "string")))
	tools.AddFunction(MakeFunctionCallTools("viewPost", "调用这个工具可以向用户发送一段不超过60s的语音，偶尔可以调用玩一下", WithParams("postId", "帖子ID", "string")))
	tools.AddFunction(MakeFunctionCallTools("speak", "调用这个工具可以向用户发送一段不超过60s的语音，偶尔可以调用玩一下", WithParams("text", "要发送的文本内容", "string")))
	res, err := json.Marshal(tools)
	if err != nil {
		panic("error in json")
	}
	log.Println(string(res))
	log.Println("yes")
}
