package main

import (
	"log"

	db "github.com/fthiswrld/ToDoBot/database"
)

func main() {
	qwe := 0
	usersCollection := db.ConnectDB()
	botToken := "6328249486:AAG401acBcRinH0GzB8nRtJ98v7dZl7Tmlg"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := GetUpdates(botUrl, offset)
		if err != nil {
			log.Println(err)
		}
		for _, update := range updates {
			err := Respond(update, botUrl, usersCollection, &qwe)
			if err != nil {
				log.Println(err)
			}
			offset = update.UpdateId + 1
		}
	}
}
