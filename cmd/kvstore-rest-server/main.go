package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "grpc-kv-store/proto" // Import the generated protobuf package

	"github.com/gorilla/mux" // Popular HTTP router
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // For insecure connection to gRPC server
)

// Response struct for JSON API responses
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

// Request struct for POST /kv/{key}
type StoreRequest struct {
	Value string `json:"value"`
}

// restAPIServer holds the gRPC client for communication with the KVStore.
type restAPIServer struct {
	grpcClient pb.KVStoreClient
}

// NewRestAPIServer creates a new instance of restAPIServer.
func NewRestAPIServer(client pb.KVStoreClient) *restAPIServer {
	return &restAPIServer{
		grpcClient: client,
	}
}

// storeHandler handles POST /kv/{key} requests.
func (s *restAPIServer) storeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var reqBody StoreRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second) // Set a timeout for gRPC call
	defer cancel()

	// Call the gRPC Store method
	resp, err := s.grpcClient.Store(ctx, &pb.StoreRequest{Key: key, Value: reqBody.Value})
	if err != nil {
		log.Printf("gRPC Store call failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, APIResponse{
		Success: resp.Success,
		Message: "Key '" + key + "' stored successfully",
	})
}

// retrieveHandler handles GET /kv/{key} requests.
func (s *restAPIServer) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	ctx, cancel := context.WithTimeout(r.Context(), time.Second) // Set a timeout for gRPC call
	defer cancel()

	// Call the gRPC Retrieve method
	resp, err := s.grpcClient.Retrieve(ctx, &pb.RetrieveRequest{Key: key})
	if err != nil {
		log.Printf("gRPC Retrieve call failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if resp.Found {
		jsonResponse(w, http.StatusOK, APIResponse{
			Success: true,
			Key:     key,
			Value:   resp.Value,
		})
	} else {
		jsonResponse(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Key '" + key + "' not found",
		})
	}
}

// deleteHandler handles DELETE /kv/{key} requests.
func (s *restAPIServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	ctx, cancel := context.WithTimeout(r.Context(), time.Second) // Set a timeout for gRPC call
	defer cancel()

	// Call the gRPC Delete method
	resp, err := s.grpcClient.Delete(ctx, &pb.DeleteRequest{Key: key})
	if err != nil {
		log.Printf("gRPC Delete call failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if resp.Success {
		jsonResponse(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Key '" + key + "' deleted successfully",
		})
	} else {
		jsonResponse(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Key '" + key + "' not found",
		})
	}
}

// jsonResponse is a helper function to send JSON responses.
func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func main() {
	// Establish a gRPC connection to the backend server.
	// The gRPC server runs on localhost:50051 (as defined in its main.go).
	// Using WithTransportCredentials(insecure.NewCredentials()) for simplicity;
	// for production, use secure credentials (e.g., TLS).
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to gRPC server: %v", err)
	}
	defer conn.Close() // Close the connection when main exits

	// Create a new gRPC client for the KVStore service.
	grpcClient := pb.NewKVStoreClient(conn)
	apiServer := NewRestAPIServer(grpcClient)

	// Initialize Gorilla Mux router.
	r := mux.NewRouter()

	// Define API routes.
	r.HandleFunc("/kv/{key}", apiServer.storeHandler).Methods("POST")
	r.HandleFunc("/kv/{key}", apiServer.retrieveHandler).Methods("GET")
	r.HandleFunc("/kv/{key}", apiServer.deleteHandler).Methods("DELETE")

	// Start the HTTP server.
	port := ":8080"
	log.Printf("REST API server listening on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
