package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPVZServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPVZList(ctx, &pb.GetPVZListRequest{})
	if err != nil {
		log.Fatalf("Failed to get PVZ list: %v", err)
	}

	fmt.Printf("Received %d PVZs\n", len(resp.Pvzs))
	for i, pvz := range resp.Pvzs {
		fmt.Printf("PVZ %d: ID=%s, City=%s, Registration Date=%s\n",
			i+1, pvz.Id, pvz.City, pvz.RegistrationDate.AsTime().Format("2006-01-02"))
	}
}
