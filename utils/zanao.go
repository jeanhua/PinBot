package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jeanhua/PinBot/config"
)

type Zanao struct {
	token string
}

func (z *Zanao) NewZanao(token string) {
	z.token = token
}

func getM(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = byte(r.Intn(10)) + '0'
	}
	return string(result)
}

func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getHeaders(userToken, schoolalias string) map[string]string {
	m := getM(20)
	td := time.Now().Unix()
	signString := fmt.Sprintf("%s_%s_%d_1b6d2514354bc407afdd935f45521a8c", schoolalias, m, td)
	return map[string]string{
		"X-Sc-Version":  "3.4.4",
		"X-Sc-Nwt":      "wifi",
		"X-Sc-Wf":       "",
		"X-Sc-Nd":       m,
		"X-Sc-Cloud":    "0",
		"X-Sc-Platform": "windows",
		"X-Sc-Appid":    "wx3921ddb0258ff14f",
		"X-Sc-Alias":    schoolalias,
		"X-Sc-Od":       userToken,
		"Content-Type":  "application/x-www-form-urlencoded",
		"X-Sc-Ah":       md5Hash(signString),
		"xweb_xhr":      "1",
		"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF WindowsWechat(0x63090c33)XWEB/14185",
		"X-Sc-Td":       strconv.FormatInt(td, 10),
		"Accept":        "*/*",
	}
}

func TrimSpaceAndBreakLine(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
}

const contentHeader = "获取到如下内容：\n"

func (z *Zanao) GetNewest(fromTime string) string {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.x.zanao.com/thread/v2/list?from_time=%s", fromTime), nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	headers := getHeaders(config.GetConfig().GetString("function_call_config.zanao_token"), "scu")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	var posts postsList
	err = json.Unmarshal(respBytes, &posts)
	if err != nil {
		log.Println(err)
		return ""
	}
	var back strings.Builder
	back.WriteString(contentHeader)
	for _, post := range posts.Data.List {
		back.WriteString("帖子ID：" + post.ID + "\n")
		back.WriteString("昵称: " + TrimSpaceAndBreakLine(post.Nickname) + "\n")
		back.WriteString("标题: " + TrimSpaceAndBreakLine(post.Title) + "\n")
		back.WriteString("内容: " + TrimSpaceAndBreakLine(post.Content) + "\n")
		back.WriteString("浏览量: " + strconv.Itoa(post.ViewCount) + "\n")
		back.WriteString("点赞数: " + post.LikeCount + "\n")
		back.WriteString("时间戳: " + post.PTime + "\n\n")
	}
	return back.String()
}

func (z *Zanao) GetHot() string {
	// https://api.x.zanao.com/thread/hot?count=10&type=3
	request, err := http.NewRequest(http.MethodPost, "https://api.x.zanao.com/thread/hot?count=10&type=3", nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	headers := getHeaders(config.GetConfig().GetString("function_call_config.zanao_token"), "scu")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	var posts hotList
	err = json.Unmarshal(respBytes, &posts)
	if err != nil {
		log.Println(err)
		return ""
	}
	var back strings.Builder
	back.WriteString(contentHeader)
	for _, post := range posts.Data.List {
		back.WriteString("帖子ID：" + post.ID + "\n")
		back.WriteString("昵称: " + TrimSpaceAndBreakLine(post.Nickname) + "\n")
		back.WriteString("标题: " + TrimSpaceAndBreakLine(post.Title) + "\n")
		back.WriteString("内容: " + TrimSpaceAndBreakLine(post.Content) + "\n")
		back.WriteString("浏览量: " + strconv.Itoa(post.ViewCount) + "\n\n")
	}
	return back.String()
}

func (z *Zanao) GetDetail(id string) string {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.x.zanao.com/thread/info?id=%s", id), nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	headers := getHeaders(config.GetConfig().GetString("function_call_config.zanao_token"), "scu")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	var post singlePost
	err = json.Unmarshal(respBytes, &post)
	if err != nil {
		log.Println(err)
		return ""
	}
	var back strings.Builder
	back.WriteString(contentHeader)
	back.WriteString("昵称: " + TrimSpaceAndBreakLine(post.Data.Detail.Nickname) + "\n")
	back.WriteString("标题: " + TrimSpaceAndBreakLine(post.Data.Detail.Title) + "\n")
	back.WriteString("内容: " + TrimSpaceAndBreakLine(post.Data.Detail.Content) + "\n")
	back.WriteString("浏览量: " + strconv.Itoa(post.Data.Detail.ViewCount) + "\n")
	back.WriteString("点赞数: " + post.Data.Detail.LikeCount + "\n\n")
	return back.String()
}

func (z *Zanao) Search(keyWords string) string {
	link := "https://api.x.zanao.com/thread/v2/search?wd=" + url.QueryEscape(keyWords) + "&cur_page=1&cate_id=10"
	request, err := http.NewRequest(http.MethodPost, link, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	headers := getHeaders(config.GetConfig().GetString("function_call_config.zanao_token"), "scu")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	var posts postsList
	err = json.Unmarshal(respBytes, &posts)
	if err != nil {
		log.Println(err)
		return ""
	}
	var back strings.Builder
	back.WriteString(contentHeader)
	for _, post := range posts.Data.List {
		back.WriteString("帖子ID：" + post.ID + "\n")
		back.WriteString("昵称: " + TrimSpaceAndBreakLine(post.Nickname) + "\n")
		back.WriteString("标题: " + TrimSpaceAndBreakLine(post.Title) + "\n")
		back.WriteString("内容: " + TrimSpaceAndBreakLine(post.Content) + "\n")
		back.WriteString("浏览量: " + strconv.Itoa(post.ViewCount) + "\n")
		back.WriteString("点赞数: " + post.LikeCount + "\n\n")
	}
	return back.String()
}

func (z *Zanao) GetComments(id string) string {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.x.zanao.com/comment/list?id=%s", id), nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	headers := getHeaders(config.GetConfig().GetString("function_call_config.zanao_token"), "scu")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	var comments comments
	err = json.Unmarshal(respBytes, &comments)
	if err != nil {
		log.Println(err)
		return ""
	}
	sort.Sort(&comments)
	var back strings.Builder
	back.WriteString(contentHeader)
	if len(comments.Data.List) > 10 {
		comments.Data.List = comments.Data.List[:10]
	}
	for _, c := range comments.Data.List {
		back.WriteString("评论者: " + c.Nickname + "\n")
		back.WriteString("内容: " + c.Content + "\n")
		back.WriteString("点赞数: " + c.LikeNum + "\n")
		back.WriteString("回复列表: \n")
		for _, reply := range c.ReplyList {
			back.WriteString(fmt.Sprintf("\t%s: %s\n", reply.Nickname, TrimSpaceAndBreakLine(reply.Content)))
		}
		back.WriteString("\n")
	}
	return back.String()
}

type postsList struct {
	Data struct {
		List []struct {
			ID        string `json:"thread_id"`  // 帖子ID
			Nickname  string `json:"nickname"`   // 昵称
			Title     string `json:"title"`      // 标题
			Content   string `json:"content"`    // 内容
			ViewCount int    `json:"view_count"` // 浏览量
			LikeCount string `json:"l_count"`    // 点赞数
			PTime     string `json:"p_time"`     // 时间戳
		} `json:"list"`
	} `json:"data"`
}

type singlePost struct {
	Data struct {
		Detail struct {
			Nickname  string `json:"nickname"`   // 昵称
			Title     string `json:"title"`      // 标题
			Content   string `json:"content"`    // 内容
			ViewCount int    `json:"view_count"` // 浏览量
			LikeCount string `json:"like_num"`   // 点赞数
		} `json:"detail"`
	} `json:"data"`
}

type hotList struct {
	Data struct {
		List []struct {
			ID        string `json:"thread_id"`  // 帖子ID
			Nickname  string `json:"nickname"`   // 昵称
			Title     string `json:"title"`      // 标题
			Content   string `json:"content"`    // 内容
			ViewCount int    `json:"view_count"` // 浏览量
		} `json:"list"`
	} `json:"data"`
}

type comments struct {
	Data struct {
		List []struct {
			Nickname  string `json:"nickname"`
			Content   string `json:"content"`
			LikeNum   string `json:"like_num"`
			ReplyList []struct {
				Nickname string `json:"nickname"`
				Content  string `json:"content"`
			} `json:"reply_list"`
		} `json:"list"`
	} `json:"data"`
}

func (c *comments) Len() int {
	return len(c.Data.List)
}

func (c *comments) Less(i, j int) bool {
	return c.Data.List[i].LikeNum > c.Data.List[j].LikeNum
}

func (c *comments) Swap(i, j int) {
	c.Data.List[i], c.Data.List[j] = c.Data.List[j], c.Data.List[i]
}
