package service

type ChatID string

type Chat interface {
	Publish(msg []byte)
	Subscribe(ch chan []byte)
}
