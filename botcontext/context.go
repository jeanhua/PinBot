package botcontext

import (
	"io"
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/botplugin"
	"github.com/jeanhua/PinBot/config"
)

/**------------------------------**/
/**
* 插件注册
**/
func RegisterPlugin(instance *BotContext) {
	instance.Plugins.AddFriendPlugin(botplugin.ExampleFriendPlugin)
	instance.Plugins.AddGroupPlugin(botplugin.ExampleGroupPlugin)
}

/**------------------------------**/

type BotContext struct {
	Plugins *botplugin.BotPlugin
}

func (bot *BotContext) NewBot() *BotContext {
	instance := &BotContext{
		Plugins: &botplugin.BotPlugin{},
	}
	RegisterPlugin(instance)
	return instance
}

func (bot *BotContext) Run() {
	config.LoadConfig()
	InitAIModelMap()
	startHTTPServer(bot)
}

func startHTTPServer(bot *BotContext) {
	http.HandleFunc("/Pinbot", bot.handler)
	log.Println("Server starting on http://localhost:7823...")
	log.Fatal(http.ListenAndServe(":7823", nil))
}

func (bot *BotContext) handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	HandleMessage(body, bot)
}
