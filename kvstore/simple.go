package kvstore

type Map struct {
	data map[string]interface{}
}

func NewMap() *Map {
	return &Map{
		data: make(map[string]interface{}),
	}
}

func (m *Map) Get(key []byte) ([]byte, bool) {
	val, ok := m.data[string(key)]
	return val.([]byte), ok
}

func (m *Map) Set(key []byte, val []byte) {
	m.data[string(key)] = val
}

func (m *Map) Del(key []byte) int {
	delete(m.data, string(key))
	return 1
}

func (m *Map) BulkDel(keys [][]byte) int {
	keysDeleted := 0
	for _, key := range keys {
		if _, ok := m.data[string(key)]; ok {
			delete(m.data, string(key))
			keysDeleted++
		}
	}
	return keysDeleted
}
