package main

import (
	"log"
	"net"

	"/kv-grpc-service/internal"

	pb "github.com/mbasha/grpc-kv-store/kv-grpc-service/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterKVStoreServer(s, internal.NewKVStoreServer())

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
