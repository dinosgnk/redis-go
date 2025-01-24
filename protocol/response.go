package protocol

import "fmt"

func BulkStringRespone(val []byte) []byte {
	// TODO
	// account for CRLF inside value
	return []byte("+" + string(val) + "\r\n")
}

func IntResponse(num int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", num))
}

func NullBulkStringRespone() []byte {
	return []byte("$-1\r\n")
}

func SimpleStringRespone(val []byte) []byte {
	return []byte("+" + string(val) + "\r\n")
}
