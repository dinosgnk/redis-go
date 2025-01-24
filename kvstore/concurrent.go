package kvstore

import (
	"sync"
)

type ConcurrentMap struct {
	mutex sync.RWMutex
	data  map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[string]interface{}),
	}
}

func (cm *ConcurrentMap) Get(key []byte) ([]byte, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if val, ok := cm.data[string(key)]; ok {
		if val, ok := val.([]byte); ok {
			return val, ok
		}
	}
	return nil, false
}

func (cm *ConcurrentMap) Set(key []byte, val []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.data[string(key)] = val
}

func (cm *ConcurrentMap) Del(keys [][]byte) int {
	keysDeleted := 0
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, key := range keys {
		if _, ok := cm.data[string(key)]; ok {
			delete(cm.data, string(key))
			keysDeleted++
		}
	}
	return keysDeleted
}

func (cm *ConcurrentMap) HSet(key []byte, field []byte, val []byte) int {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if existingVal, valExists := cm.data[string(key)]; valExists {
		if hash, valueIsMap := existingVal.(map[string][]byte); valueIsMap {
			hash[string(field)] = val
			return 0
		}
	}
	cm.data[string(key)] = map[string][]byte{string(field): val}
	return 1
}

func (cm *ConcurrentMap) HGet(key []byte, field []byte) ([]byte, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	if existingVal, valExists := cm.data[string(key)]; valExists {
		if hash, valueIsMap := existingVal.(map[string][]byte); valueIsMap {
			val, ok := hash[string(field)]
			return val, ok
		}
	}
	return nil, false
}
