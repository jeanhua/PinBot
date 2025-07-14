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

// AI
var aiModelMap map[uint]aicommunicate.AiModel

// 大模型速率限制
var llmLock sync.Mutex

func Register() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Println("error: ", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config.ConfigInstance)
	if err != nil {
		log.Println("error config: ", err)
	}
	aiModelMap = make(map[uint]aicommunicate.AiModel, 0)
	http.HandleFunc("/Pinbot", Handler)
	log.Println("Server starting on http://localhost:7823...")
	log.Fatal(http.ListenAndServe(":7823", nil))
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
		for _, uin := range config.ConfigInstance.Friend.Exclude {
			if uin == strconv.Itoa(friendmsg.UserId) {
				return
			}
		}
		for _, uin := range config.ConfigInstance.Friend.Include {
			if uin == "all" || uin == strconv.Itoa(friendmsg.UserId) {
				onPrivateMessage(friendmsg)
				return
			}
		}

	} else if friendmsg.MessageType == "group" {
		groupmsg := model.GroupMessage{}
		err := json.Unmarshal(message, &groupmsg)
		if err != nil {
			log.Println("error:", err)
			return
		}
		for _, uin := range config.ConfigInstance.Group.Exclude {
			if uin == strconv.Itoa(groupmsg.GroupId) {
				return
			}
		}
		for _, uin := range config.ConfigInstance.Group.Include {
			if uin == "all" || uin == strconv.Itoa(groupmsg.GroupId) {
				onGroupMessage(groupmsg)
				return
			}
		}
	}
}
