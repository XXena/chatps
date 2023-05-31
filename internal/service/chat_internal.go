package service

import (
	"fmt"
	"sync"

	"github.com/XXena/chatps/pkg/logger"
)

type chat struct {
	hub    Hub
	ID     ChatID
	mu     sync.RWMutex
	subs   map[ChatID][]SendChan
	logger *logger.Logger
}

// Publish чат рассылает сообщение клиента всем остальным подписанным на чат клиентам
func (c *chat) Publish(msg []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, ch := range c.subs[c.ID] {
		ch <- msg
		fmt.Println("ch", ch, "msg", msg)
	}
}

// Subscribe клиент подписывается на чат
func (c *chat) Subscribe(ch chan []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subs[c.ID] = append(c.subs[c.ID], ch)
}

func NewInternalChat(hub Hub, ID ChatID, logger *logger.Logger) Chat {

	return &chat{
		hub:    hub,
		ID:     ID,
		subs:   make(map[ChatID][]SendChan),
		logger: logger,
	}
}
