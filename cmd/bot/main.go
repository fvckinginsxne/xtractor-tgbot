package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"bot/internal/clients/tgclient"
	"bot/internal/consumer/eventconsumer"
	"bot/internal/listener"
	"bot/pkg/postgres/audiostorage"
	"bot/pkg/postgres/urlstorage"
	"bot/pkg/postgres/userstorage"
)

const (
	BATCH_SIZE     = 100
	MSG_QUEUE_SIZE = 100

	ENV_PATH = "/Users/madw3y/petprojects/xtractor-tgbot/.env"
)

func main() {
	if err := godotenv.Load(ENV_PATH); err != nil {
		log.Fatal("can't loading .env file: ", err)
	}

	dbConn := os.Getenv("DB_URL")

	userStorage, err := userstorage.New(dbConn)
	if err != nil {
		log.Fatal("can't connect to user storage: ", err)
	}

	urlStorage, err := urlstorage.New(dbConn)
	if err != nil {
		log.Fatal("can't connect to user storage: ", err)
	}

	audioStorage, err := audiostorage.New(dbConn, urlStorage)
	if err != nil {
		log.Fatal("can't connect to audio storage: ", err)
	}

	if err := userStorage.Init(); err != nil {
		log.Fatal("can't init user storage: ", err)
	}

	if err := urlStorage.Init(); err != nil {
		log.Fatal("can't init url storage: ", err)
	}

	if err := audioStorage.Init(); err != nil {
		log.Fatal("can't init audio storage: ", err)
	}

	tgclient := tgclient.New(os.Getenv("HOSTNAME"), os.Getenv("TOKEN"))

	if err := tgclient.SetCommandsList(); err != nil {
		log.Fatal("can't set commands list", err)
	}

	listener := listener.New(tgclient, audioStorage, userStorage, urlStorage)

	consumer := eventconsumer.New(*listener, BATCH_SIZE, MSG_QUEUE_SIZE)

	consumer.Start()
}
