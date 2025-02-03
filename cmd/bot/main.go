package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"bot/internal/clients/tgclient"
	"bot/internal/consumer/eventconsumer"
	"bot/internal/listener"
	"bot/pkg/postgres/audiostorage"
	"bot/pkg/postgres/userstorage"
)

const (
	batchSize = 100
)

func main() {
	log.Printf("service started")

	if err := godotenv.Load("/Users/madw3y/petprojects/extracter-tgbot/.env"); err != nil {
		log.Fatal("can't loading .env file: ", err)
	}

	userStorage, err := userstorage.New(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("can't connect to user storage: ", err)
	}

	audioStorage, err := audiostorage.New(os.Getenv("DB_URL"), userStorage)
	if err != nil {
		log.Fatal("can't connect to audio storage: ", err)
	}

	if err := userStorage.Init(); err != nil {
		log.Fatal("can't init user storage: ", err)
	}

	if err := audioStorage.Init(); err != nil {
		log.Fatal("can't init audio storage: ", err)
	}

	tgclient := tgclient.New(os.Getenv("HOSTNAME"), os.Getenv("TOKEN"))

	if err := tgclient.SetCommandsList(); err != nil {
		log.Fatal("can't set commands list", err)
	}

	listener := listener.New(tgclient, audioStorage, userStorage)

	consumer := eventconsumer.New(*listener, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped")
	}
}
