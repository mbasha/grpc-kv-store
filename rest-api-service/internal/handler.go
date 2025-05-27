package handler

import (
    "encoding/json"
    "net/http"

    "google.golang.org/grpc"
    "go-kv-store/rest-api-service/internal/client"
)

type KeyValueHandler struct {
    kvClient client.KVClient
}

func NewKeyValueHandler(kvClient client.KVClient) *KeyValueHandler {
    return &KeyValueHandler{kvClient: kvClient}
}

func (h *KeyValueHandler) StoreValue(w http.ResponseWriter, r *http.Request) {
    var req client.StoreRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.kvClient.Store(req.Key, req.Value); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *KeyValueHandler) RetrieveValue(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value, err := h.kvClient.Retrieve(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    response := map[string]string{"key": key, "value": value}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *KeyValueHandler) DeleteValue(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if err := h.kvClient.Delete(key); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}