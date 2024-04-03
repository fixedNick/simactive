package grpc

import (
	"fmt"
	"log"
	"net"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/config"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	port    int
	timeout time.Duration
}

func NewGRPCServer(cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		port:    cfg.GRPC.Port,
		timeout: cfg.GRPC.Timeout,
	}
}

func (gs *GRPCServer) MustRun(simService SimService) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", gs.port))
	if err != nil {
		log.Fatalf("error starting gRPC server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSimServer(s, NewGRPCSimService(simService, gs.timeout))

	log.Print("Starting gRPC server...")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("error serving gRPC requests: %v", err)
	}
}
