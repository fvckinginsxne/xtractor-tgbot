package eventconsumer

import (
	"errors"
	"log"
	"sync"
	"time"

	"bot/internal/core"
	"bot/internal/listener"
	"bot/internal/tech/e"
)

type Consumer struct {
	listener     listener.Listener
	batchSize    int
	msgQueueSize int
	mu           sync.Mutex
	queues       map[int]chan core.Event
}

func New(listener listener.Listener, batchSize, msgQueueSize int) *Consumer {
	return &Consumer{
		listener:     listener,
		batchSize:    batchSize,
		msgQueueSize: msgQueueSize,
		queues:       make(map[int]chan core.Event),
	}
}

func (c *Consumer) Start() {
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

		c.handleEvents(gotEvents)
	}
}

func (c *Consumer) handleEvents(events []core.Event) {
	for _, event := range events {
		c.mu.Lock()
		userQueue, exists := c.queues[event.ChatID]
		if !exists {
			userQueue = make(chan core.Event, c.msgQueueSize)

			c.queues[event.ChatID] = userQueue
			go c.processUserEvents(userQueue)
		}
		c.mu.Unlock()

		userQueue <- event
	}
}

func (c *Consumer) processUserEvents(queue chan core.Event) {
	for event := range queue {
		if err := c.listener.Process(event); err != nil {
			log.Print("can't handle event for user", err.Error())
		}
	}
}
