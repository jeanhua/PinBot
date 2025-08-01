package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func WebSearch(token, query string, timeRange *string, include, exclude []string, count int) string {
	endPoint := "https://api.tavily.com/search"
	bodyMap := map[string]any{}
	bodyMap["query"] = query
	if timeRange != nil {
		bodyMap["time_range"] = *timeRange
	}
	if len(include) != 0 {
		bodyMap["include_domains"] = include
	}
	if len(exclude) != 0 {
		bodyMap["exclude_domains"] = exclude
	}
	bodyMap["max_results"] = count
	bodyBytes, err := json.Marshal(&bodyMap)
	if err != nil {
		log.Println("error in marshal json: webSearch request body")
		return "查询失败"
	}
	httpRequest, err := http.NewRequest(http.MethodPost, endPoint, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("error in create request: webSearch")
		return "查询失败"
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Println("error in sendResponse: webSearch")
		return "查询失败"
	}
	if resp.StatusCode != 200 {
		log.Println("error in response: webSearch, status code:", resp.StatusCode)
		return "查询失败"
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error in read body: webSearch")
		return "查询失败"
	}
	searchResponse := searchResponse{}
	err = json.Unmarshal(respbytes, &searchResponse)
	if err != nil {
		log.Println("error in json unmarshal in response: webSearch")
		log.Println(string(respbytes))
		return "查询失败"
	}
	return searchResponse.toString()
}

type searchResponse struct {
	Query   string `json:"query"`
	Results []struct {
		URL     string  `json:"url"`
		Title   string  `json:"title"`
		Content string  `json:"content"`
		Score   float64 `json:"score"`
	} `json:"results"`
	ResponseTime float64 `json:"response_time"` // 响应时间
}

func (sr *searchResponse) toString() string {
	var output string

	output += "Search Query: " + sr.Query + "\n"
	output += fmt.Sprintf("Total Results: %d\n", len(sr.Results))
	output += fmt.Sprintf("Response Time: %.3f seconds\n", sr.ResponseTime)
	output += "\nResults:\n"

	for i, result := range sr.Results {
		output += fmt.Sprintf("Result #%d:\n", i+1)
		output += "  标题:   " + result.Title + "\n"
		output += "  URL:     " + result.URL + "\n"
		output += "  内容: " + result.Content + "\n"
		output += "  置信度:   " + strconv.FormatFloat(result.Score, 'f', -1, 64) + "\n"
		output += "\n"
	}
	return output
}
