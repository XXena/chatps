package service

import (
	"github.com/XXena/chatps/internal/entity"
	"github.com/XXena/chatps/pkg/logger"
)

// hub maintains the set of active chats
type hub struct {
	// Registered chats grouped by id
	registeredChats map[entity.ChatID]Chat

	// Inbound messages from the connections
	broadcast chan entity.Message

	register chan ClientsService

	unregister chan ClientsService
	logger     *logger.Logger
}

func NewHub(logger *logger.Logger) Hub {
	return hub{
		broadcast:       make(chan entity.Message),
		register:        make(chan ClientsService),
		unregister:      make(chan ClientsService),
		registeredChats: make(map[entity.ChatID]Chat),
		logger:          logger,
	}
}

func (h hub) Run() error {
	for {
		select {
		case conn := <-h.register:
			chat, ok := h.registeredChats[conn.GetChatID()]
			if ok {
				chat.Subscribe(conn.GetSendChan())
				continue
			}

			// новый чат:
			chat = NewInternalChat(h, conn.GetChatID(), h.logger)
			h.registeredChats[conn.GetChatID()] = chat
			chat.Subscribe(conn.GetSendChan())

		case conn := <-h.unregister:
			if _, ok := h.registeredChats[conn.GetChatID()]; ok {
				// todo unsubscribe?
				close(conn.GetSendChan())
				delete(h.registeredChats, conn.GetChatID())
			}

		case m := <-h.broadcast:
			chat, ok := h.registeredChats[m.ChatID]
			if ok {
				chat.Publish(m.Data)
			}
		}
	}
}

func (h hub) Unregister(c ClientsService) {
	h.unregister <- c
}

func (h hub) Register(c ClientsService) {
	h.register <- c
}

// Broadcast sends new message to the channel
func (h hub) Broadcast(m entity.Message) {
	h.broadcast <- m
}
