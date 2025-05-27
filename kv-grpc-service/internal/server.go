package kvstore

import (
	"context"
	"log"
	"net"

	pb "github.com/mbasha/grpc-kv-store/kv-grpc-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type server struct {
	pb.UnimplementedKVStoreServer
	store map[string]string
}

func NewServer() *server {
	return &server{
		store: make(map[string]string),
	}
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.store[req.Key] = req.Value
	return &pb.SetResponse{Success: true}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value, exists := s.store[req.Key]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "key not found")
	}
	return &pb.GetResponse{Value: value}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	delete(s.store, req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterKVStoreServer(grpcServer, NewServer())
	log.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
