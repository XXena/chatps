package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/XXena/chatps/internal/entity/proto"
	"github.com/XXena/chatps/internal/service"

	"github.com/XXena/chatps/internal/handlers"

	"github.com/XXena/chatps/pkg/logger"

	"github.com/XXena/chatps/internal/config"
	grpcTransport "github.com/XXena/chatps/pkg/transport/grpc"
	httpTransport "github.com/XXena/chatps/pkg/transport/http"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	errChan := make(chan error, 1)
	hub := service.NewHub(l)
	socket := handlers.WebsocketHandler{Hub: hub, Cfg: cfg, Logger: l}
	router := socket.InitRoutes()

	go func() {
		err := hub.Run()
		l.Error(fmt.Errorf("app - Run - hub initializing error: %w", err))
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		httpTransport.RunServer(cfg, l, router, errChan)
	}()

	grpcHandler := handlers.NewGRPCService(hub, cfg, l)

	var opts []grpc.ServerOption // todo
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatPSServer(grpcServer, grpcHandler)
	go func() {
		grpcTransport.RunServer(grpcServer, cfg.GRPC.Port, l, errChan) // todo config
	}()

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	// todo завершение работы прил-я
	// todo отключение пользователей от сервера

	select {
	case s := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", s))
	case err := <-errChan:
		l.Error(fmt.Errorf("app - Run - error notify: %w", err))
		//case <-wsClient.Done:
		//	log.Println("Receiver Channel Closed! Exiting....") // todo
	}
	return

}
