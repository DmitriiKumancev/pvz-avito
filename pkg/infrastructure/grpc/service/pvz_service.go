package service

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/application/services/pvz"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZServiceServer struct {
	pb.UnimplementedPVZServiceServer
	pvzService pvz.Service
}

func NewPVZServiceServer(pvzService pvz.Service) *PVZServiceServer {
	return &PVZServiceServer{
		pvzService: pvzService,
	}
}

func (s *PVZServiceServer) GetPVZList(ctx context.Context, req *pb.GetPVZListRequest) (*pb.GetPVZListResponse, error) {
	filter := repositories.PVZFilter{
		Page:  1,
		Limit: 100, // Set a reasonable limit
	}

	pvzs, err := s.pvzService.ListPVZs(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoPVZs := make([]*pb.PVZ, 0, len(pvzs))
	for _, p := range pvzs {
		protoPVZs = append(protoPVZs, &pb.PVZ{
			Id:               p.ID,
			RegistrationDate: timestamppb.New(p.RegistrationDate),
			City:             p.City,
		})
	}

	return &pb.GetPVZListResponse{
		Pvzs: protoPVZs,
	}, nil
}
