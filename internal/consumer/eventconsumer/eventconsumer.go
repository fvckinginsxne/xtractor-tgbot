package eventconsumer

import (
	"errors"
	"log"
	"time"

	"bot/internal/core"
	"bot/internal/listener"
)

type Consumer struct {
	eventprocessor listener.Listener
	batchSize      int
}

func New(eventprocessor listener.Listener, batchSize int) *Consumer {
	return &Consumer{
		eventprocessor: eventprocessor,
		batchSize:      batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.eventprocessor.Fetch(c.batchSize)
		if err != nil && !errors.Is(err, listener.ErrNoUpdates) {
			log.Printf("consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err.Error())
		}
	}
}

func (c Consumer) handleEvents(events []core.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.eventprocessor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
		}
	}

	return nil
}
