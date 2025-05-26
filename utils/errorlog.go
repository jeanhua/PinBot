package utils

import (
	"log"

	"github.com/jeanhua/PinBot/config"
)

func LogErr(msg string) {
	config.ConfigInstance_mu.RLock()
	dbg := config.ConfigInstance.Debug
	config.ConfigInstance_mu.RUnlock()
	if dbg {
		log.Println(msg)
	}
}
