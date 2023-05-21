package service

import (
	"github.com/XXena/chatps/pkg/logger"
)

type Chat struct {
	hub    IHub
	ID     ChatID
	queue  chan Message
	logger *logger.Logger
}

func (c *Chat) Enqueue(msg Message) {
	c.queue <- msg
}

func (c *Chat) Dequeue() Message {
	return <-c.queue
}

func NewChat(hub IHub, ID ChatID, logger *logger.Logger) Chat {
	return Chat{
		hub:    hub,
		ID:     ID,
		queue:  make(chan Message),
		logger: logger,
	}
}
