package service

import "github.com/XXena/chatps/pkg/logger"

type grpcClient struct {
	// The grpc connection
	// todo
	chatID   string
	sendChan SendChan
	logger   *logger.Logger
	Hub      Hub
}
