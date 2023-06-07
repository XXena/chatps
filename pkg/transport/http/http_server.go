package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/XXena/chatps/internal/config"
	"github.com/XXena/chatps/pkg/logger"
)

func RunServer(cfg *config.Config, log *logger.Logger, router *chi.Mux, listenErr chan error) {
	log.Info("starting http server on %s", cfg.WebSocket.Server.Port)

	rt, err := time.ParseDuration(cfg.WebSocket.Server.ReadTimeout)
	if err != nil {
		log.Error("failed to run http server: unable to parse read timeout", err)
		listenErr <- err
		return
	}

	wt, err := time.ParseDuration(cfg.WebSocket.Server.WriteTimeout)
	if err != nil {
		log.Error("failed to run http server: unable to parse write timeout", err)
		listenErr <- err
		return
	}

	srv := http.Server{
		Addr:         cfg.WebSocket.Server.Host + cfg.WebSocket.Server.Port,
		Handler:      router,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		listenErr <- err
	}
}
