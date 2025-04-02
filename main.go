package main

import (
	bot "rapplab.xyz/chrisbot/bot"
)

// 277025466368 perms
var TOKEN = "MTM1NzA0NzE0NTIxNTAzNzQ5MA.GogMzo.3y4tFjUaEYRt5E8MF7TArQ8OANb81WeKH300kk"

//var TOKEN string = os.Getenv("DISCORD_TOKEN") // get the token from environment variable

func main() {
	bot.BotToken = TOKEN
	bot.Run() // call the run function of bot/bot.go
}
