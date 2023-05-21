package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/XXena/chatps/internal/service"

	ws "github.com/XXena/chatps/internal/websocket"

	"github.com/XXena/chatps/pkg/logger"

	"github.com/XXena/chatps/internal/config"
	httpTransport "github.com/XXena/chatps/internal/transport/http"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	errChan := make(chan error, 1)
	hub := service.NewHub(l)
	socket := ws.Handler{Hub: hub, Cfg: cfg, Logger: l}
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

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
