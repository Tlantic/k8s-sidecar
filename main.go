package main

import (
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	"github.com/Tlantic/k8s-sidecar/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterK8SServiceServer(s, server.NewK8sService())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("listening on port %s\n", port)
}
