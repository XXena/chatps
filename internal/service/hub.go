package service

import "github.com/XXena/chatps/internal/entity"

type Hub interface {
	Run() error
	Unregister(conn ClientsService)
	Register(conn ClientsService)
	Broadcast(m entity.Message)
}
