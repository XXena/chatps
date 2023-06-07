package service

type Chat interface {
	Publish(msg []byte)
	Subscribe(ch chan []byte)
}
