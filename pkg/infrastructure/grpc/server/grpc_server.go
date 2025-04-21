package server

import (
	"fmt"
	"net"

	"github.com/dkumancev/avito-pvz/pkg/application/services/pvz"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/pb"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server     *grpc.Server
	pvzService pvz.Service
	port       int
}

func NewGRPCServer(pvzService pvz.Service, port int) *GRPCServer {
	return &GRPCServer{
		pvzService: pvzService,
		port:       port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	s.server = grpc.NewServer()

	pvzServiceServer := service.NewPVZServiceServer(s.pvzService)
	pb.RegisterPVZServiceServer(s.server, pvzServiceServer)

	reflection.Register(s.server)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
