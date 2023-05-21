package main

import (
	"log"

	"github.com/XXena/chatps/internal/app"

	"github.com/XXena/chatps/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
