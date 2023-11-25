package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	From      User   `json:"from"`
}

type Chat struct {
	Id int `json:"id"`
}

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type SendMessage struct {
	ChatId      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}
type SendMessage2 struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]Keyboardbutton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}
type Keyboardbutton struct {
	Text string `json:"text"`
}
