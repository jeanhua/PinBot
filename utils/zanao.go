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
	"strconv"
	"time"

	"github.com/jeanhua/PinBot/config"
)

type Zanao struct{}

func getM(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = byte(r.Intn(10)) + '0'
	}
	return string(result)
}

func getH() int64 {
	return time.Now().Unix()
}

func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getB(input string) string {
	return md5Hash(input)
}

func getResult(userToken, schoolalias string) map[string]string {
	m := getM(20)
	td := getH()
	signString := fmt.Sprintf("%s_%s_%d_1b6d2514354bc407afdd935f45521a8c", schoolalias, m, td)
	b := getB(signString)
	return map[string]string{
		"X-Sc-Ah":    b,
		"X-Sc-Alias": schoolalias,
		"X-Sc-Nd":    m,
		"X-Sc-Od":    userToken,
		"X-Sc-Td":    strconv.FormatInt(td, 10),
	}
}

type SimpleResponse struct {
	Data struct {
		List []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"list"`
	} `json:"data"`
}

func getResponse(token, url string) (*SimpleResponse, error) {
	header := getResult(token, "scu")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req, _ := http.NewRequest(http.MethodGet, url+"&from_time="+timestamp, nil)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var jsonResp SimpleResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}
	return &jsonResp, nil
}

func (Zanao) GetNewest() (*SimpleResponse, error) {
	config.Config_mu.RLock()
	token := config.Config.ZanaoToken
	config.Config_mu.RUnlock()
	return getResponse(token, "https://api.x.zanao.com/thread/v2/list?from_time=0&hot=1&with_comment=true&with_reply=true")
}

func (Zanao) GetHot() (*SimpleResponse, error) {
	config.Config_mu.RLock()
	token := config.Config.ZanaoToken
	config.Config_mu.RUnlock()
	return getResponse(token, "https://api.x.zanao.com/thread/hot?count=10&type=3&with_comment=true&with_reply=true")
}
