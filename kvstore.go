package main

import "sync"

type KVStore struct {
	mutex sync.RWMutex
	data  map[string][]byte
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: map[string][]byte{},
	}
}

func (kv *KVStore) Get(key []byte) ([]byte, bool) {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}

func (kv *KVStore) Set(key []byte, val []byte) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[string(key)] = val
	return nil
}
