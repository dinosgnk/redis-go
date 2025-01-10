package main

import (
	"sync"
)

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

func (kv *KVStore) Set(key []byte, val []byte) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[string(key)] = val
}

func (kv *KVStore) Del(key []byte) int {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	delete(kv.data, string(key))
	return 1
}

func (kv *KVStore) BulkDel(keys [][]byte) int {
	keysDeleted := 0
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	for _, key := range keys {
		if _, ok := kv.data[string(key)]; ok {
			delete(kv.data, string(key))
			keysDeleted++
		}
	}
	return keysDeleted
}
