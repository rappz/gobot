package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	bot "rapplab.xyz/chrisbot/bot"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Note: .env file not loaded (run from project root or set env vars): ", err)
	}
	fmt.Println("Starting bot...")
	fmt.Println("Token: ", os.Getenv("TOKEN"))
	bot.BotToken = os.Getenv("TOKEN")
	bot.Run() // call the run function of bot/bot.go
	fmt.Println("Bot started")
}
