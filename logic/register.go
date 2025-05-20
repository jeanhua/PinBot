package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	llm "github.com/jeanhua/PinBot/LLM"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
	"gopkg.in/yaml.v3"
)

// 智谱AI
var zhipu *llm.ZhiPu

// 大模型速率限制
var llmLock sync.RWMutex
var ready bool

// 配置
var config_mu sync.RWMutex
var config model.Config

func Register() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	config_mu.Lock()
	err = decoder.Decode(&config)
	config_mu.Unlock()
	if err != nil {
		fmt.Println("error config: ", err)
	}
	go watchConfig()
	zhipu = llm.NewZhiPu()
	http.HandleFunc("/Pinbot", Handler)
	log.Println("Server starting on http://localhost:7823...")
	ready = true
	log.Fatal(http.ListenAndServe(":7823", nil))
}

func watchConfig() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	filemd5 := ""
	for range ticker.C {
		hash, err := utils.FileMD5("./config.yaml")
		if err != nil {
			log.Println("read config error: ", err)
		}
		if hash == filemd5 {
			continue
		}
		filemd5 = hash
		file, err := os.Open("./config.yaml")
		if err != nil {
			log.Println("error: ", err)
		}
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		config_mu.Lock()
		err = decoder.Decode(&config)
		config_mu.Unlock()
		if err != nil {
			log.Println("error config: ", err)
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error: ", err)
	}
	handleMessage(body)
	r.Body.Close()
}

func handleMessage(message []byte) {
	basic := &model.Message{}
	err := json.Unmarshal(message, &basic)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if basic.PostType != "message" {
		return
	}
	friendmsg := model.FriendMessage{}
	err = json.Unmarshal(message, &friendmsg)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if friendmsg.MessageType == "private" {
		config_mu.RLock()
		defer config_mu.RUnlock()
		for _, uin := range config.Group.Exclude {
			if uin == strconv.Itoa(friendmsg.UserId) {
				return
			}
		}
		for _, uin := range config.Friend.Include {
			if uin == "all" || uin == strconv.Itoa(friendmsg.UserId) {
				onPrivateMessage(friendmsg)
				return
			}
		}

	} else if friendmsg.MessageType == "group" {
		config_mu.RLock()
		defer config_mu.RUnlock()
		groupmsg := model.GroupMessage{}
		err := json.Unmarshal(message, &groupmsg)
		if err != nil {
			log.Println("error:", err)
			return
		}
		for _, uin := range config.Group.Exclude {
			if uin == strconv.Itoa(groupmsg.GroupId) {
				return
			}
		}
		for _, uin := range config.Group.Include {
			if uin == "all" || uin == strconv.Itoa(groupmsg.GroupId) {
				onGroupMessage(groupmsg)
				return
			}
		}
	}
}
