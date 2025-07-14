package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WebSearch(token, query, freshness string, summary bool, include, exclude string, count int) string {
	endPoint := "https://api.bochaai.com/v1/web-search"
	bodyMap := map[string]any{}
	bodyMap["query"] = query
	bodyMap["freshness"] = freshness
	bodyMap["summary"] = summary
	bodyMap["include"] = include
	bodyMap["exclude"] = exclude
	bodyMap["count"] = count
	bodyBytes, err := json.Marshal(&bodyMap)
	if err != nil {
		fmt.Println("error in marshal json: webSearch request body")
		return "查询失败"
	}
	httpRequest, err := http.NewRequest(http.MethodPost, endPoint, bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Println("error in create request: webSearch")
		return "查询失败"
	}
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	httpRequest.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		fmt.Println("error in sendResponse: webSearch")
		return "查询失败"
	}
	if resp.StatusCode != 200 {
		fmt.Println("error in response: webSearch, status code:", resp.StatusCode)
		return "查询失败"
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in read body: webSearch")
		return "查询失败"
	}
	searchResponse := SearchResponse{}
	err = json.Unmarshal(respbytes, &searchResponse)
	if err != nil {
		fmt.Println("error in json unmarshal in response: webSearch")
		fmt.Println(string(respbytes))
		return "查询失败"
	}
	return searchResponse.ToString()
}

type SearchResponse struct {
	Code int `json:"code"`
	Data struct {
		QueryContext struct {
			OriginalQuery string `json:"originalQuery"`
		} `json:"queryContext"`
		WebPages struct {
			WebSearchUrl string `json:"webSearchUrl"`
			Value        []struct {
				Name            string `json:"name"`            // 网页的标题
				URL             string `json:"url"`             // 网页的URL
				DisplayURL      string `json:"displayUrl"`      // 网页的展示URL（url decode后的格式）
				Snippet         string `json:"snippet"`         // 网页内容的简短描述
				Summary         string `json:"summary"`         // 网页内容的文本摘要，当请求参数中 summary 为 true 时显示此属性
				SiteName        string `json:"siteName"`        // 网页的网站名称
				SiteIcon        string `json:"siteIcon"`        // 网页的网站图标
				DateLastCrawled string `json:"dateLastCrawled"` // 网页的发布时间
			} `json:"value"`
		} `json:"webPages"`
	} `json:"data"`
}

func (sr *SearchResponse) ToString() string {
	var b strings.Builder

	if sr.Data.QueryContext.OriginalQuery != "" {
		b.WriteString(fmt.Sprintf("Query Context: %s\n", sr.Data.QueryContext.OriginalQuery))
	}
	webPages := sr.Data.WebPages
	b.WriteString(fmt.Sprintf("Web Search URL: %s\n", webPages.WebSearchUrl))
	b.WriteString("Results:\n")

	for i, result := range webPages.Value {
		b.WriteString(fmt.Sprintf("  %d. Title: %s\n", i+1, result.Name))
		b.WriteString(fmt.Sprintf("     URL: %s\n", result.URL))
		b.WriteString(fmt.Sprintf("     Display URL: %s\n", result.DisplayURL))
		b.WriteString(fmt.Sprintf("     Snippet: %s\n", result.Snippet))
		if result.Summary != "" {
			b.WriteString(fmt.Sprintf("     Summary: %s\n", result.Summary))
		}
		b.WriteString(fmt.Sprintf("     Site Name: %s\n", result.SiteName))
		b.WriteString(fmt.Sprintf("     Date Last Crawled: %s\n", result.DateLastCrawled))
		b.WriteString("\n")
	}
	return b.String()
}
