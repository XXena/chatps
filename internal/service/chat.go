package service

import (
	"sync"

	"github.com/XXena/chatps/pkg/logger"
)

type Chat struct {
	hub    IHub
	ID     ChatID
	mu     sync.RWMutex
	subs   map[ChatID][]chan []byte
	logger *logger.Logger
}

func (c *Chat) Publish(msg []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, ch := range c.subs[c.ID] {
		ch <- msg
	}
}

func (c *Chat) Subscribe(ch chan []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subs[c.ID] = append(c.subs[c.ID], ch)
}

func NewChat(hub IHub, ID ChatID, logger *logger.Logger) Chat {
	return Chat{
		hub:    hub,
		ID:     ID,
		subs:   make(map[ChatID][]chan []byte),
		logger: logger,
	}
}
