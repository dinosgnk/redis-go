package protocol

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"redis-go/message"
	"strconv"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(rd),
	}
}

func (p *Parser) ParseArray(arrayHeader []byte) [][]byte {
	// CRLF of 1st line has been removed before calling this function
	numOfExpectedElements, err := strconv.ParseInt(string(arrayHeader[1:]), 10, 64)
	if err != nil {
		return nil
	}

	elements := make([][]byte, 0, numOfExpectedElements)
	for i := int64(0); i < numOfExpectedElements; i++ {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return nil
		}

		// get line length to get all digits of $ header
		lineLength := len(line)

		// read after $ until before CRLF
		elementLength, err := strconv.ParseInt(string(line[1:lineLength-2]), 10, 64)
		if err != nil {
			return nil
		}

		// element length + CRLF
		body := make([]byte, elementLength+2)

		// read next element + CRLF bytes
		_, err = io.ReadFull(p.reader, body)
		if err != nil {
			return nil
		}

		// append the element to the list, without the trailing CRLF
		elements = append(elements, body[:len(body)-2])
	}

	return elements
}

func (p *Parser) Parse(conn net.Conn, clientMsgCh chan<- *message.Message) error {

	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return nil
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		_type := line[0]

		switch _type {
		case '*':
			cmdArgs := p.ParseArray(line)
			msg := message.NewMessage(cmdArgs, conn)
			clientMsgCh <- msg
		default:
			log.Println("TODO")
			continue
		}
	}
}
