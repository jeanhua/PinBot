package utils

import (
	"log"

	"github.com/jeanhua/PinBot/config"
)

func LogErr(msg string) {
	dbg := config.ConfigInstance.Debug
	if dbg {
		log.Println(msg)
	}
}
