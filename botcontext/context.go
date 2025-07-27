package botcontext

import "github.com/jeanhua/PinBot/logic"

type BotContext struct{}

func (bot *BotContext) Run() {
	logic.Register()
}
