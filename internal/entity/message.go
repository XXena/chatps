package entity

type ChatID string

type Message struct {
	ChatID ChatID
	Data   []byte
}
