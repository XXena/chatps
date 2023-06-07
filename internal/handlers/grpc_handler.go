package handlers

import (
	"fmt"

	"github.com/XXena/chatps/internal/config"
	pb "github.com/XXena/chatps/internal/entity/proto"
	"github.com/XXena/chatps/internal/service"
	"github.com/XXena/chatps/pkg/logger"
)

type Service struct {
	pb.UnimplementedChatPSServer
	Hub    service.Hub
	Cfg    *config.Config
	Logger *logger.Logger
}

func NewGRPCService(Hub service.Hub, cfg *config.Config, l *logger.Logger) *Service {
	return &Service{
		Hub:    Hub,
		Cfg:    cfg,
		Logger: l,
	}
}

func (s Service) Exchange(stream pb.ChatPS_ExchangeServer) error {
	grpcClientService := service.NewGRPCClientService(stream, s.Hub, s.Logger)
	grpcClientService.GetHub().Register(grpcClientService)
	errChan := make(chan error, 1)
	go func() {
		err := grpcClientService.PullMessage()
		if err != nil {
			s.Logger.Error(fmt.Errorf("grpc handler - PullMessage - error: %w", err))
		}
	}()

	s.Logger.Debug("goroutine PullMessage passed")
	go func() {
		err := grpcClientService.SendMessage()
		if err != nil {
			s.Logger.Error(fmt.Errorf("grpc handler - SendMessage - error: %w", err))
		}
	}()
	s.Logger.Debug("goroutine SendMessage passed")

	select {
	case err := <-errChan:
		s.Logger.Error(fmt.Errorf("app - Run - error notify: %w", err))
	}

	s.Logger.Debug("grpc handler finished")
	return nil
}
