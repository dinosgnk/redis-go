package kvstore

import "sync"

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
	val, ok := cm.data[string(key)]
	return val.([]byte), ok
}

func (cm *ConcurrentMap) Set(key []byte, val []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.data[string(key)] = val
}

func (cm *ConcurrentMap) Del(key []byte) int {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.data, string(key))
	return 1
}

func (cm *ConcurrentMap) BulkDel(keys [][]byte) int {
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
