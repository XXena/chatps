package service

import (
	"github.com/XXena/chatps/pkg/logger"
)

type ChatID string

// Hub maintains the set of active chats
type Hub struct {
	// Registered chats grouped by id
	registeredChats map[ChatID]Chat

	// Inbound messages from the connections
	broadcast chan Message

	register chan ChatID

	unregister chan Connection
	logger     *logger.Logger
}

type IHub interface {
	Run() error
	Unregister(connection Connection)
	Register(ChatID)
	Broadcast(m Message)
	//GetMessage() Message
}

func NewHub(logger *logger.Logger) IHub {
	return Hub{
		broadcast:       make(chan Message),
		register:        make(chan ChatID),
		unregister:      make(chan Connection),
		registeredChats: make(map[ChatID]Chat),
		logger:          logger,
	}
}

func (h Hub) Run() error {
	for {
		select {
		case chatID := <-h.register:
			chat, ok := h.registeredChats[chatID]
			if ok {
				continue
			}

			// новый чат:
			chat = NewChat(h, chatID, h.logger)
			h.registeredChats[chatID] = chat
			// todo subscribe - подписаться на chat

		case conn := <-h.unregister:
			if _, ok := h.registeredChats[conn.GetChatID()]; ok {
				close(conn.GetSendChan())
				delete(h.registeredChats, conn.GetChatID())
			}

		case m := <-h.broadcast:
			chat, ok := h.registeredChats[m.ChatID]
			if ok {
				chat.Enqueue(m)
				// todo пришло новое сообщ-е, это событие для других подписчиков - для них dequeue
			}
		}
	}
}

func (h Hub) Unregister(c Connection) {
	h.unregister <- c
}

func (h Hub) Register(id ChatID) {
	h.register <- id
}

// Broadcast sends new message to the channel
func (h Hub) Broadcast(m Message) {
	h.broadcast <- m
}

//// GetMessage reads message from the channel
//func (h Hub) GetMessage() (m Message) {
//	m = <-h.broadcast
//
//	return m
//}
