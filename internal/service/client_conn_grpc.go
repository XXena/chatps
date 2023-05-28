package service

type GrpcConnection struct {
	hub Hub

	// The grpc connection
	// todo

	// Buffered channel for outbound messages
	send chan []byte
}
