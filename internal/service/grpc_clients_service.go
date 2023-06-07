package service

import (
	"errors"
	"fmt"
	"io"

	"github.com/XXena/chatps/internal/entity"

	pb "github.com/XXena/chatps/internal/entity/proto"
	"github.com/XXena/chatps/pkg/logger"
)

type grpcClientsService struct {
	grpc     pb.ChatPS_ExchangeServer
	chatID   string
	sendChan SendChan
	logger   *logger.Logger
	Hub      Hub
}

func NewGRPCClientService(grpcStream pb.ChatPS_ExchangeServer, hub Hub, logger *logger.Logger) ClientsService {
	return &grpcClientsService{
		Hub:      hub,
		grpc:     grpcStream,
		sendChan: make(SendChan),
		logger:   logger,
	}
}

func (g grpcClientsService) GetChatID() entity.ChatID {
	return entity.ChatID(g.chatID)

}

func (g grpcClientsService) SetChatID(id string) {
	g.chatID = id
}

func (g grpcClientsService) GetSendChan() chan []byte {
	return g.sendChan

}

func (g grpcClientsService) GetHub() Hub {
	return g.Hub
}

func (g grpcClientsService) SendMessage() error {
	defer func() {
		g.Hub.Unregister(g)
	}()
	for {
		msg, err := g.grpc.Recv()
		if err != nil {
			if err == io.EOF {
				break // todo
			}
			g.logger.Error(fmt.Errorf("error recieving message from grpc stream: %w", err))
			return err
		}

		if msg.ChatID == "" {
			g.logger.Error(fmt.Errorf("chat ID not specified, http response write error: %w", err))
			return errors.New("you must specify chat ID")
		} // todo валидировать proto?

		g.SetChatID(msg.ChatID)

		m := entity.Message{
			ChatID: entity.ChatID(msg.ChatID),
			Data:   []byte(msg.Data), // todo
		}

		g.Hub.Broadcast(m)
	}
	return nil
}

// PullMessage promotes message from the send channel to the grpc connection
func (g grpcClientsService) PullMessage() error {
	for {
		select {
		case message, ok := <-g.sendChan:
			if !ok {
				return errors.New("grpc message not received")
			}
			if err := g.grpc.Send(&pb.Message{
				ChatID: g.chatID,
				Data:   string(message),
			}); err != nil {
				g.logger.Error(fmt.Errorf("grpc message recieved but send error: %w", err))
			}
		}
	}

}
