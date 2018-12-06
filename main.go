package main

import (
	"fmt"
	"github.com/Tlantic/k8s-sidecar/internal/manager"
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	"github.com/Tlantic/k8s-sidecar/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	port := ":50051"
	if os.Getenv("SIDECAR_PORT") != "" {
		port = fmt.Sprintf(":%s", os.Getenv("SIDECAR_PORT"))
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	kubeManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    os.Getenv("KUBECONFIG"),
		Namespace: os.Getenv("K8S_NAMESPACE"),
		Timeout:   10,
	})
	if err != nil {
		panic(err)
	}

	pb.RegisterK8SServiceServer(s, server.NewK8sService(kubeManager))
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("listening on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
