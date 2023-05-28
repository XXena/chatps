package service

type Hub interface {
	Run() error
	Unregister(conn Client)
	Register(conn Client)
	Broadcast(m Message)
}
