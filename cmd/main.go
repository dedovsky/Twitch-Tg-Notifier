package main

import (
	"awesomeProject/services/telegram"
	"awesomeProject/services/twitch"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
)

func main() {
	pref := tele.Settings{
		Token: os.Getenv("TELEGRAM_TOKEN"),
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	telegram.PushRoutes(b)
	go twitch.ListenAndServe(b)

	log.Println("Bot started")

	b.Start()
}
