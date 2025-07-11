package functioncall

func CallFunction(functionCall *FunctionCall) string {
	switch functionCall.Action {
	case "browse_homepage":
		{
			return browseHomepage()
		}
	case "browse_hot":
		{
			return browseHot()
		}
	case "search":
		{
			return search(functionCall.Param["keyword"].(string))
		}
	case "view_post":
		{
			return viewPost(functionCall.Param["post_id"].(string))
		}
	case "view_comments":
		{
			return viewComments(functionCall.Param["post_id"].(string))
		}
	}
	return ""
}
