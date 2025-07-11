package functioncall

import "github.com/jeanhua/PinBot/utils"

type FunctionCall struct {
	Action string         `json:"action"`
	Param  map[string]any `json:"parameters"`
}

func browseHomepage() string {
	zanao := &utils.Zanao{}
	return zanao.GetNewest()
}

func search(keywords string) string {
	zanao := &utils.Zanao{}
	return zanao.Search(keywords)
}

func viewPost(postId string) string {
	zanao := &utils.Zanao{}
	return zanao.GetDetail(postId)
}

func browseHot() string {
	zanao := &utils.Zanao{}
	return zanao.GetHot()
}

func viewComments(postId string) string {
	zanao := &utils.Zanao{}
	return zanao.GetComments(postId)
}
