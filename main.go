package main

import (
	"os"

	"github.com/subosito/gotenv"
	bot "rapplab.xyz/chrisbot/bot"
)

func init() {
	gotenv.Load()
}

func main() {
	bot.BotToken = os.Getenv("TOKEN")
	bot.Run() // call the run function of bot/bot.go
}
