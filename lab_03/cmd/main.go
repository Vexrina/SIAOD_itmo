package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"vexrina/siaod_itmo/lab_03/internal"
	"vexrina/siaod_itmo/lab_03/internal/api"
	pb "vexrina/siaod_itmo/lab_03/lab_03/pkg/lab_03"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	uc := api.NewLabUseCase()
	server := internal.NewService(uc)
	pb.RegisterLab03Server(s, server)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
