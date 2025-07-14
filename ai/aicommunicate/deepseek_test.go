package aicommunicate

import (
	"encoding/json"
	"log"
	"testing"
)

func TestDeepseekFuncs(t *testing.T) {
	tools := &FunctionCall{}
	tools.AddFunction(MakeFunctionCallTools(
		"webSearch",
		"执行网络搜索，用于获取互联网相关信息",
		WithParams("query", "搜索查询内容", "string", true),
		WithParams("timeRange", "限制搜索结果的时间范围(可选)(day,week,month,year)", "string", false),
		WithParams("include", "限定搜索结果必须包含的域名列表(可选)", "array<string>", false),
		WithParams("exclude", "排除特定域名的搜索结果(可选)", "array<string>", false),
		WithParams("count", "返回的最大搜索结果数量(可选)", "int", false),
	))
	tools.AddFunction(MakeFunctionCallTools(
		"webExplore",
		"根据提供的链接列表抓取网页内容或进一步探索",
		WithParams("links", "要抓取或探索的网页链接数组", "array<string>", true),
	))
	tools.AddFunction(MakeFunctionCallTools("browseHomepage", "浏览校园集市论坛主页", WithParams("fromTime", "时间戳,该时间戳前的10条帖子,输入0则表示最新的10条帖子,通过获取帖子后再输入最后帖子的时间戳来实现翻页", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("browseHot", "浏览校园集市论坛热门帖子"))
	tools.AddFunction(MakeFunctionCallTools("searchPost", "搜索校园集市论坛帖子", WithParams("keywords", "搜索关键词", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("viewComments", "浏览校园集市论坛指定帖子的评论", WithParams("postId", "帖子ID", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("viewPost", "查看校园集市论坛某个帖子详情", WithParams("postId", "帖子ID", "string", true)))
	tools.AddFunction(MakeFunctionCallTools("speak", "调用这个工具可以向用户发送一段不超过60s的语音，偶尔可以调用玩一下", WithParams("text", "要发送的文本内容", "string", true)))
	res, err := json.Marshal(tools)
	if err != nil {
		panic("error in json")
	}
	log.Println(string(res))
	log.Println("yes")
}
