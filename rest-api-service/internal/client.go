package client

import (
	"context"

	pb "github.com/mbasha/grpc-kv-store/kv-grpc-service/proto"
	"google.golang.org/grpc"
)

type KVClient struct {
	client pb.KVStoreClient
}

func NewKVClient(conn *grpc.ClientConn) *KVClient {
	return &KVClient{
		client: pb.NewKVStoreClient(conn),
	}
}

func (c *KVClient) Set(ctx context.Context, key string, value string) error {
	_, err := c.client.Set(ctx, &pb.SetRequest{Key: key, Value: value})
	return err
}

func (c *KVClient) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func (c *KVClient) Delete(ctx context.Context, key string) error {
	_, err := c.client.Delete(ctx, &pb.DeleteRequest{Key: key})
	return err
}

func (c *KVClient) Close() {
	// No resources to close for the client
}
