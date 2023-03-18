package main

import (
	"log"
	"os"

	"time"

	"github.com/nacknime-official/gdz-ukraine/internal/controller/telegram"
	"github.com/nacknime-official/gdz-ukraine/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/vitaliy-ukiru/fsm-telebot"
	redis_storage "github.com/vitaliy-ukiru/fsm-telebot/storages/redis"
	"gopkg.in/telebot.v3"
)

func main() {
	pref := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	storage := redis_storage.NewStorage(redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0}), redis_storage.StorageSettings{})
	manager := fsm.NewManager(bot, nil, storage)

	homeworkService := service.NewMockHomeworkService()
	h := telegram.NewUserHandler(homeworkService)
	h.Register(manager)

	go bot.Start()
	log.Println("Bot started")
	for {
	}
}
