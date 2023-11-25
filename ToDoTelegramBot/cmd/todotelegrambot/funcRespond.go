package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	db "github.com/fthiswrld/ToDoBot/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUpdates(boturl string, offset int) ([]Update, error) {
	resp, err := http.Get(boturl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func Respond(update Update, botUrl string, usersCollection *mongo.Collection, qwe *int) error {
	if update.Message.Text == "/start" {
		var sendMessage SendMessage
		var button1 Keyboardbutton
		var button2 Keyboardbutton
		var button3 Keyboardbutton
		var buttons [][]Keyboardbutton
		var b []Keyboardbutton
		button1.Text = "Список задач"
		button2.Text = "Добавить задачу"
		button3.Text = "Удалить задачу"
		b = append(b, button1)
		b = append(b, button2)
		b = append(b, button3)
		buttons = append(buttons, b)
		sendMessage.ChatId = update.Message.Chat.Id
		sendMessage.Text = "Привет, я помощник по твоим задачам, можешь выбрать что мне делать:)"
		sendMessage.ReplyMarkup.Keyboard = buttons
		sendMessage.ReplyMarkup.OneTimeKeyboard = false
		sendMessage.ReplyMarkup.ResizeKeyboard = true
		err := sMessage(sendMessage, botUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else if update.Message.Text == "Список задач" {
		var sm SendMessage2
		result := listToDo(usersCollection, update)
		if len(result.Tasks) == 0 {
			sm.ChatId = update.Message.Chat.Id
			sm.Text = "Пока что у вас нету задач:("
			err := sMessage2(sm, botUrl)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			text := ""
			for i := 0; i < len(result.Tasks); i++ {
				example := fmt.Sprintf("%d) %s.\n", i+1, result.Tasks[i])
				text += example
				fmt.Println(text)
			}
			sm.ChatId = update.Message.Chat.Id
			sm.Text = text
			err := sMessage2(sm, botUrl)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		}
	} else if update.Message.Text == "Добавить задачу" || *qwe == 1 {
		if update.Message.Text == "Добавить задачу" {
			var sm SendMessage2
			sm.ChatId = update.Message.Chat.Id
			sm.Text = "Напишите вашу задачу, она будет добавленно в список!"
			err := sMessage2(sm, botUrl)
			if err != nil {
				panic(err)
			}
			*qwe = 1
		} else if *qwe == 1 {
			addTask(usersCollection, update)
			var sm SendMessage2
			sm.ChatId = update.Message.Chat.Id
			sm.Text = "Задание успешно добавленно!"
			err := sMessage2(sm, botUrl)
			if err != nil {
				panic(err)
			}
			*qwe = 0
		}
	} else if update.Message.Text == "Удалить задачу" || *qwe == 2 {
		if update.Message.Text == "Удалить задачу" {
			var sm SendMessage2
			sm.ChatId = update.Message.Chat.Id
			sm.Text = "Введите номер задачи"
			err := sMessage2(sm, botUrl)
			if err != nil {
				panic(err)
			}
			*qwe = 2
		} else if *qwe == 2 {
			user := listToDo(usersCollection, update)
			key, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				var sm SendMessage2
				sm.ChatId = update.Message.Chat.Id
				sm.Text = "Введите корректный номер задачи!"
				err := sMessage2(sm, botUrl)
				if err != nil {
					panic(err)
				}
			}
			if len(user.Tasks) < key {
				var sm SendMessage2
				sm.ChatId = update.Message.Chat.Id
				sm.Text = "Введите корректный номер задачи!"
				err := sMessage2(sm, botUrl)
				if err != nil {
					panic(err)
				}
			} else {
				user.Tasks = append(user.Tasks[:key-1], user.Tasks[key:]...)
				fmt.Println(user.Tasks)
				filter := bson.D{{Key: "telegram", Value: update.Message.From.Id}}
				replace := bson.D{{Key: "telegram", Value: update.Message.From.Id}, {Key: "tasks", Value: user.Tasks}}
				usersCollection.FindOneAndReplace(context.TODO(), filter, replace)
				var sm SendMessage2
				sm.ChatId = update.Message.Chat.Id
				sm.Text = "Задача успешно удаленна!"
				err = sMessage2(sm, botUrl)
				if err != nil {
					panic(err)
				}
				*qwe = 0
			}

		}
	} else {
		var sm SendMessage2
		sm.ChatId = update.Message.Chat.Id
		sm.Text = "Пожалуйста, выберите действие!"
		err := sMessage2(sm, botUrl)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func listToDo(collection *mongo.Collection, update Update) db.Users {
	var result db.Users
	filter := bson.D{bson.E{Key: "telegram", Value: update.Message.From.Id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println(result)
	if err != nil {
		_, err := collection.InsertOne(context.TODO(), db.Users{Telegram_id: update.Message.From.Id, Tasks: []string{}})
		if err != nil {
			panic(err)
		}
	}
	return result
}

func sMessage(sendMessage SendMessage, botUrl string) error {
	buf, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
func sMessage2(sendMessage SendMessage2, botUrl string) error {
	buf, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
func addTask(collection *mongo.Collection, update Update) {
	var result db.Users
	filter := bson.D{{Key: "telegram", Value: update.Message.From.Id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		_, err := collection.InsertOne(context.TODO(), db.Users{Telegram_id: update.Message.From.Id, Tasks: []string{}})
		if err != nil {
			panic(err)
		}
	}
	adding_task := append(result.Tasks, update.Message.Text)
	replace := bson.D{{Key: "telegram", Value: update.Message.From.Id}, {Key: "tasks", Value: adding_task}}
	err = collection.FindOneAndReplace(context.TODO(), filter, replace).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := collection.InsertOne(context.TODO(), db.Users{Telegram_id: update.Message.From.Id, Tasks: adding_task})
			if err != nil {
				fmt.Println(184)
				panic(err)
			}
		}
	}
	collection.UpdateOne(context.TODO(), filter, replace)
}
