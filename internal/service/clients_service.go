package service

import (
	"time"

	"github.com/XXena/chatps/internal/entity"
)

const (
	// Time allowed to Write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// sendChan pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type ClientsService interface {
	SendMessage() error
	PullMessage() error
	GetChatID() entity.ChatID
	GetSendChan() chan []byte
	GetHub() Hub
}

type SendChan chan []byte
