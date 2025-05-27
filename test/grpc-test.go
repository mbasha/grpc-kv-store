package test

import (
	"context"
	"log"
	"testing"

	pb "github.com/mbasha/grpc-kv-store/kv-grpc-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:50051" // gRPC server address
)

var client pb.KVStoreClient

func setup() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client = pb.NewKVStoreClient(conn)
}

func TestStoreAndRetrieve(t *testing.T) {
	setup()
	defer client.Close()

	key := "testKey"
	value := "testValue"

	// Store the value
	_, err := client.Store(context.Background(), &pb.StoreRequest{Key: key, Value: value})
	if err != nil {
		t.Fatalf("could not store value: %v", err)
	}

	// Retrieve the value
	resp, err := client.Retrieve(context.Background(), &pb.RetrieveRequest{Key: key})
	if err != nil {
		t.Fatalf("could not retrieve value: %v", err)
	}

	if resp.Value != value {
		t.Errorf("expected value %s, got %s", value, resp.Value)
	}
}

func TestDelete(t *testing.T) {
	setup()
	defer client.Close()

	key := "testKey"

	// Delete the key
	_, err := client.Delete(context.Background(), &pb.DeleteRequest{Key: key})
	if err != nil {
		t.Fatalf("could not delete key: %v", err)
	}

	// Try to retrieve the deleted key
	_, err = client.Retrieve(context.Background(), &pb.RetrieveRequest{Key: key})
	if err == nil {
		t.Errorf("expected error when retrieving deleted key, got none")
	}

	if status.Code(err) != status.NotFound {
		t.Errorf("expected NotFound error, got %v", status.Code(err))
	}
}
