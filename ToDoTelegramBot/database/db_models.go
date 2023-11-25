package database

type Users struct {
	Telegram_id int      `bson:"telegram"`
	Tasks       []string `bson:"tasks"`
}
