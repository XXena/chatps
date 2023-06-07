package grpc

import (
	"fmt"
	"net"

	"github.com/XXena/chatps/pkg/logger"

	"google.golang.org/grpc"
)

func RunServer(s *grpc.Server, port string, log *logger.Logger, listenErr chan error) {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	log.Info("starting grpc server on %s", l.Addr())
	if err := s.Serve(l); err != nil {
		listenErr <- err
	}
}
