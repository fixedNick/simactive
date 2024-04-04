package grpc

import (
	"fmt"
	"log"
	"net"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/config"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errors
var (
	ErrInternal = status.Error(codes.Internal, "Internal Server Error")
)

type GRPCServer struct {
	port    int
	timeout time.Duration

	// gRPC services
	server *grpc.Server
}

func NewGRPCServer(cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		port:    cfg.GRPC.Port,
		timeout: cfg.GRPC.Timeout,
	}
}

func (gs *GRPCServer) MustRun(simService SimService, serviceService ServiceService) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", gs.port))
	if err != nil {
		log.Fatalf("error starting gRPC server: %v", err)
	}

	s := grpc.NewServer()
	gs.server = s
	pb.RegisterSimServer(s, NewGRPCSimService(simService, gs.timeout))
	pb.RegisterServiceServer(s, NewGRPCServiceService(serviceService, gs.timeout))

	log.Print("Starting gRPC server...")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("error serving gRPC requests: %v", err)
	}
}

func (gs *GRPCServer) Stop() {
	gs.server.GracefulStop()
}
