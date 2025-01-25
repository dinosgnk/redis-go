package protocol

import "fmt"

func BulkStringRespone(val []byte) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), string(val)))
}

func IntResponse(num int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", num))
}

func NumOfArgumentsErrorResponse(cmd []byte) []byte {
	return []byte(fmt.Sprintf("-ERR wrong number of arguments for '%v' command", cmd))
}

func NullBulkStringRespone() []byte {
	return []byte("$-1\r\n")
}

func SimpleStringRespone(val []byte) []byte {
	return []byte("+" + string(val) + "\r\n")

}
