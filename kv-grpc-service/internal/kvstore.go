package kvstore

import (
    "sync"
)

type KVStore struct {
    sync.RWMutex
    store map[string]string
}

func NewKVStore() *KVStore {
    return &KVStore{
        store: make(map[string]string),
    }
}

func (kvs *KVStore) Set(key string, value string) {
    kvs.Lock()
    defer kvs.Unlock()
    kvs.store[key] = value
}

func (kvs *KVStore) Get(key string) (string, bool) {
    kvs.RLock()
    defer kvs.RUnlock()
    value, exists := kvs.store[key]
    return value, exists
}

func (kvs *KVStore) Delete(key string) {
    kvs.Lock()
    defer kvs.Unlock()
    delete(kvs.store, key)
}