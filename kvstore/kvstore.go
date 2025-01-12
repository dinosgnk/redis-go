package kvstore

type KVStore interface {
	Get(key []byte) ([]byte, bool)
	Set(key []byte, val []byte)
	Del(key []byte) int
	BulkDel(keys [][]byte) int
}
