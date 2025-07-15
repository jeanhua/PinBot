package logic

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/jeanhua/PinBot/ai/aicommunicate"
	"github.com/jeanhua/PinBot/config"
	"github.com/jeanhua/PinBot/model"
	"gopkg.in/yaml.v3"
)

// 全局变量
var (
	repeatLock sync.Mutex
	repeat     = &repeatTuple{}
	llmLock    sync.Mutex
	aiModelMap = make(map[uint]aicommunicate.AiModel)
)

// Register 初始化并启动HTTP服务
func Register() {
	loadConfig()
	initAIModelMap()
	startHTTPServer()
}

// loadConfig 加载配置文件
func loadConfig() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config.ConfigInstance); err != nil {
		log.Fatalf("Failed to decode config: %v", err)
	}
}

// initAIModelMap 初始化AI模型映射
func initAIModelMap() {
	aiModelMap = make(map[uint]aicommunicate.AiModel)
}

// startHTTPServer 启动HTTP服务器
func startHTTPServer() {
	http.HandleFunc("/Pinbot", Handler)
	log.Println("Server starting on http://localhost:7823...")
	log.Fatal(http.ListenAndServe(":7823", nil))
}

// Handler 处理HTTP请求
func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	handleMessage(body)
}

// 处理接收到的消息
func handleMessage(message []byte) {
	basicMsg, err := parseBasicMessage(message)
	if err != nil || basicMsg.PostType != "message" {
		return
	}

	if isFriendMessage(message) {
		handleFriendMessage(message)
	} else if isGroupMessage(message) {
		handleGroupMessage(message)
	}
}

// 解析基本消息结构
func parseBasicMessage(message []byte) (*model.Message, error) {
	basic := &model.Message{}
	if err := json.Unmarshal(message, basic); err != nil {
		log.Printf("Error parsing basic message: %v", err)
		return nil, err
	}
	return basic, nil
}

// 检查是否为好友消息
func isFriendMessage(message []byte) bool {
	friendMsg := model.FriendMessage{}
	if err := json.Unmarshal(message, &friendMsg); err != nil {
		log.Printf("Error parsing friend message: %v", err)
		return false
	}
	return friendMsg.MessageType == "private"
}

// 检查是否为群组消息
func isGroupMessage(message []byte) bool {
	groupMsg := model.GroupMessage{}
	if err := json.Unmarshal(message, &groupMsg); err != nil {
		log.Printf("Error parsing group message: %v", err)
		return false
	}
	return groupMsg.MessageType == "group"
}

// 处理好友消息
func handleFriendMessage(message []byte) {
	friendMsg := model.FriendMessage{}
	if err := json.Unmarshal(message, &friendMsg); err != nil {
		log.Printf("Error parsing friend message: %v", err)
		return
	}

	if isExcludedFriend(friendMsg.UserId) {
		return
	}

	if isIncludedFriend(friendMsg.UserId) {
		onPrivateMessage(friendMsg)
	}
}

// 处理群组消息
func handleGroupMessage(message []byte) {
	groupMsg := model.GroupMessage{}
	if err := json.Unmarshal(message, &groupMsg); err != nil {
		log.Printf("Error parsing group message: %v", err)
		return
	}

	if isExcludedGroup(groupMsg.GroupId) {
		return
	}

	if isIncludedGroup(groupMsg.GroupId) {
		onGroupMessage(groupMsg)
	}
}

// 检查好友是否在排除列表中
func isExcludedFriend(userId int) bool {
	for _, uin := range config.ConfigInstance.Friend.Exclude {
		if uin == strconv.Itoa(userId) {
			return true
		}
	}
	return false
}

// 检查好友是否在包含列表中
func isIncludedFriend(userId int) bool {
	for _, uin := range config.ConfigInstance.Friend.Include {
		if uin == "all" || uin == strconv.Itoa(userId) {
			return true
		}
	}
	return false
}

// 检查群组是否在排除列表中
func isExcludedGroup(groupId int) bool {
	for _, uin := range config.ConfigInstance.Group.Exclude {
		if uin == strconv.Itoa(groupId) {
			return true
		}
	}
	return false
}

// 检查群组是否在包含列表中
func isIncludedGroup(groupId int) bool {
	for _, uin := range config.ConfigInstance.Group.Include {
		if uin == "all" || uin == strconv.Itoa(groupId) {
			return true
		}
	}
	return false
}
