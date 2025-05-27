package main

import (
	"context"
	"log"
	"net"
	"sync" // For thread-safe map access

	pb "grpc-kv-store/proto" // Import the generated protobuf package

	"google.golang.org/grpc"
)

// kvStoreServer implements the KVStore service.
type kvStoreServer struct {
	pb.UnimplementedKVStoreServer                   // Embed for forward compatibility
	mu                            sync.RWMutex      // Mutex for concurrent map access
	data                          map[string]string // In-memory key-value store
}

// NewKVStoreServer creates a new instance of kvStoreServer.
func NewKVStoreServer() *kvStoreServer {
	return &kvStoreServer{
		data: make(map[string]string),
	}
}

// Store implements the Store RPC method.
func (s *kvStoreServer) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	s.mu.Lock()         // Acquire write lock
	defer s.mu.Unlock() // Release write lock when function exits

	log.Printf("Storing key: %s, value: %s", req.Key, req.Value)
	s.data[req.Key] = req.Value // Store the key-value pair
	return &pb.StoreResponse{Success: true}, nil
}

// Retrieve implements the Retrieve RPC method.
func (s *kvStoreServer) Retrieve(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	s.mu.RLock()         // Acquire read lock
	defer s.mu.RUnlock() // Release read lock when function exits

	log.Printf("Retrieving key: %s", req.Key)
	value, found := s.data[req.Key] // Retrieve the value
	if !found {
		log.Printf("Key '%s' not found", req.Key)
		return &pb.RetrieveResponse{Value: "", Found: false}, nil
	}
	log.Printf("Key '%s' found with value: %s", req.Key, value)
	return &pb.RetrieveResponse{Value: value, Found: true}, nil
}

// Delete implements the Delete RPC method.
func (s *kvStoreServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()         // Acquire write lock
	defer s.mu.Unlock() // Release write lock when function exits

	log.Printf("Deleting key: %s", req.Key)
	_, found := s.data[req.Key] // Check if key exists
	if !found {
		log.Printf("Key '%s' not found for deletion", req.Key)
		return &pb.DeleteResponse{Success: false}, nil
	}
	delete(s.data, req.Key) // Delete the key
	log.Printf("Key '%s' deleted successfully", req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func main() {
	// Listen on TCP port 50051 for gRPC connections.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("gRPC server listening on %v", lis.Addr())

	// Create a new gRPC server instance.
	s := grpc.NewServer()

	// Register our KVStore service implementation with the gRPC server.
	pb.RegisterKVStoreServer(s, NewKVStoreServer())

	// Start serving gRPC requests. This is a blocking call.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
