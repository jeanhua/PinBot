package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const failText = "页面打开失败"

func WebExplore(links []string, token string) string {
	const requestUrl = "https://api.tavily.com/extract"
	postBody := &WebExploreRequestBody{
		Urls:         links,
		ExtractDepth: "advanced",
	}
	postBytes, err := json.Marshal(postBody)
	if err != nil {
		log.Println("error in json marshal: webExplore", err)
		return failText
	}
	httpRequest, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(postBytes))
	if err != nil {
		log.Println("error in create request: webExplore")
		return failText
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Println("error in send request: webExplore", err)
		return failText
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("error in response: webSearch, status code:", resp.StatusCode)
		return "查询失败"
	}
	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error in read response body: webExplore", err)
		return failText
	}
	var response WebExploreResponse
	err = json.Unmarshal(respbytes, &response)
	if err != nil {
		log.Println("error in jsonUnMarshal: webExplore", err)
		return failText
	}
	return response.ToReadableString()
}

type WebExploreRequestBody struct {
	Urls         []string `json:"urls"`
	ExtractDepth string   `json:"extract_depth"`
}

type WebExploreResponse struct {
	Requests []struct {
		Url        string `json:"url"`
		RawContent string `json:"raw_content"`
	} `json:"results"`
	FailedResults []struct {
		Url   string `json:"url"`
		Error string `json:"error"`
	} `json:"failed_results"`
	ResponseTime float32 `json:"response_time"`
}

func (r *WebExploreResponse) ToReadableString() string {
	var output string

	output += "=== Successful Requests ===\n"
	for i, req := range r.Requests {
		output += fmt.Sprintf("Request %d:\n", i+1)
		output += fmt.Sprintf("  URL: %s\n", req.Url)
		contentSummary := req.RawContent
		if len(contentSummary) > 20000 {
			contentSummary = contentSummary[:20000] + "..."
		}
		output += fmt.Sprintf("  Raw Content (first 20000 chars): %s\n", contentSummary)
		output += "\n"
	}

	output += "=== Failed Requests ===\n"
	for i, failed := range r.FailedResults {
		output += fmt.Sprintf("Failed Request %d:\n", i+1)
		output += fmt.Sprintf("  URL: %s\n", failed.Url)
		output += fmt.Sprintf("  Error: %s\n", failed.Error)
		output += "\n"
	}

	output += "=== Response Time ===\n"
	output += fmt.Sprintf("Total response time: %.2f seconds\n", r.ResponseTime)

	return output
}
