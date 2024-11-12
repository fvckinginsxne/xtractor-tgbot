package main

import (
	"flag"
	"log"

	"bot/internal/clients/tgclient"
	"bot/internal/consumer/eventconsumer"
	"bot/internal/listener"
	"bot/pkg/storage/sqlite"
)

const (
	batchSize       = 100
	hostname        = "api.telegram.org"
	storageBasePath = "../../data/sqlite/storage.db"
)

func main() {
	storage, err := sqlite.New(storageBasePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventprocessor := listener.New(
		tgclient.New(hostname, mustToken()),
		storage,
	)

	log.Printf("service started")

	consumer := eventconsumer.New(*eventprocessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stoped")
	}
}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
