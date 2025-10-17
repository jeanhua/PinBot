package botcontext

import (
	"fmt"
	"github.com/jeanhua/PinBot/botcontext/plugin"
	"io"
	"log"
	"net/http"

	"github.com/jeanhua/PinBot/config"
)

type BotContext struct {
	Plugins *plugin.BotPlugin
}

func NewBot() *BotContext {
	instance := &BotContext{
		Plugins: &plugin.BotPlugin{},
	}
	return instance
}

func (bot *BotContext) Run() {
	config.LoadConfig()
	startHTTPServer(bot)
}

func startHTTPServer(bot *BotContext) {
	http.HandleFunc("/Pinbot", bot.handler)
	log.Printf("Server starting on http://localhost:%d...\n", config.GetConfig().LocalListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().LocalListenPort), nil))
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
