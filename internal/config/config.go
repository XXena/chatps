package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		App       `env:",prefix=APP_"`
		HTTP      `env:",prefix=HTTP_"`
		WebSocket `env:",prefix=WS_"`
		GRPC      `env:",prefix=GRPC_"`
		Log
		Chat
	}

	App struct {
		Name    string `env-required:"true"  env:"NAME"`
		Version string `env-required:"true" env:"VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" env:"PORT"`
	}

	WebSocket struct {
		Server WebSocketServer
	}

	WebSocketServer struct {
		Host         string `env:"HOST,default=localhost"`
		Port         string `env-required:"true" env:"PORT"`
		ReadTimeout  string `env-required:"true" env:"READ_TIMEOUT"`
		WriteTimeout string `env-required:"true" env:"WRITE_TIMEOUT"`
	}

	GRPC struct {
		Port string `env-required:"true" env:"PORT"`
	}

	Log struct {
		Level string `env-required:"true"  env:"LOG_LEVEL"`
	}

	Chat struct {
		SendBufferSize int `env:"SEND_BUFFER_SIZE,default=256"` // todo
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	var envFiles []string
	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")
		envFiles = append(envFiles, ".env")
	}
	if os.Getenv("APP_ENV") != "" {
		appEnvName := fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))
		if _, err := os.Stat(appEnvName); err == nil {
			log.Println("found", appEnvName, "file, adding it to env config files list")
			envFiles = append(envFiles, appEnvName)
		}
	}
	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			return nil, errors.New("error while opening env config: %s")
		}
	}

	ctx := context.Background()
	err := envconfig.Process(ctx, cfg)

	if err != nil {
		return nil, errors.New("unable to load config from file")
	}

	return cfg, nil
}
