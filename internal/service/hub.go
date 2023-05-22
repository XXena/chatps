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

	register chan Client

	unregister chan Client
	logger     *logger.Logger
}

type IHub interface {
	Run() error
	Unregister(conn Client)
	Register(conn Client)
	Broadcast(m Message)
	//GetMessage() Message
}

func NewHub(logger *logger.Logger) IHub {
	return Hub{
		broadcast:       make(chan Message),
		register:        make(chan Client),
		unregister:      make(chan Client),
		registeredChats: make(map[ChatID]Chat),
		logger:          logger,
	}
}

func (h Hub) Run() error {
	for {
		select {
		case conn := <-h.register:
			chat, ok := h.registeredChats[conn.GetChatID()]
			if ok {
				chat.Subscribe(conn.GetSendChan())
				continue
			}

			// новый чат:
			chat = NewChat(h, conn.GetChatID(), h.logger)
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

func (h Hub) Unregister(c Client) {
	h.unregister <- c
}

func (h Hub) Register(c Client) {
	h.register <- c
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
