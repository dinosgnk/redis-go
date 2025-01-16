package message

import "net"

type Message struct {
	Header  []byte
	CmdArgs [][]byte
	Conn    net.Conn
}

func NewMessage(cmdArgs [][]byte, conn net.Conn) *Message {
	return &Message{
		Header:  cmdArgs[0],
		CmdArgs: cmdArgs,
		Conn:    conn,
	}
}
