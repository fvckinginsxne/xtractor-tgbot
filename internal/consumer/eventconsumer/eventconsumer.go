package eventconsumer

import (
	"errors"
	"log"
	"time"

	"bot/internal/core"
	"bot/internal/listener"
	"bot/pkg/tech/e"
)

type Consumer struct {
	listener  listener.Listener
	batchSize int
}

func New(listener listener.Listener, batchSize int) *Consumer {
	return &Consumer{
		listener:  listener,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.listener.Fetch(c.batchSize)
		if err != nil && !errors.Is(err, e.ErrNoUpdates) {
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
		go func() {
			if err := c.listener.Process(event); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}()
	}

	return nil
}
