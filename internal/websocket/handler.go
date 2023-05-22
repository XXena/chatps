package websocket

import (
	"fmt"
	"net/http"

	"github.com/XXena/chatps/internal/config"
	"github.com/XXena/chatps/internal/service"

	"github.com/XXena/chatps/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{} // use default options
	Done     chan interface{}       // todo
)

type Handler struct {
	Hub    service.IHub
	Cfg    *config.Config
	Logger *logger.Logger
}

func (h Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})
	r.Get("/ws/{chat}", h.socketHandler)
	return r
}

func (h Handler) socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error(fmt.Errorf("error during connection upgradation: %w", err))
		return
	}

	chatID := chi.URLParam(r, "chat")

	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("you must specify chat ID"))
		if err != nil {
			h.Logger.Error(fmt.Errorf("chat ID not specified, http response write error: %w", err))
			return
		}
		return
	}

	client := service.NewWebsocketClient(
		h.Hub,
		wsConn,
		chatID,
		h.Logger)

	client.GetHub().Register(client)

	errChan := make(chan error, 1)
	go func() {
		err := client.PullMessage()
		h.Logger.Error(fmt.Errorf("socket handler - PullMessage - error: %w", err))
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		err := client.SendMessage()
		h.Logger.Error(fmt.Errorf("socket handler - SendMessage - error: %w", err))
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		h.Logger.Error(fmt.Errorf("app - Run - error notify: %w", err))
	}

}
