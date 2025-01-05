package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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

func (p *Parser) ParseArray(arrayHeader []byte) error {
	// CRLF of 1st line has been removed before calling this function
	numOfExpectedElements, err := strconv.ParseInt(string(arrayHeader[1:]), 10, 64)
	if err != nil {
		return err
	}

	elements := make([][]byte, 0, numOfExpectedElements)
	for i := int64(0); i < numOfExpectedElements; i++ {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		// get line length to get all digits of $ header
		lineLength := len(line)

		// read after $ until before CRLF
		elementLength, err := strconv.ParseInt(string(line[1:lineLength-2]), 10, 64)
		fmt.Println("Element length", elementLength)
		if err != nil {
			return err
		}

		// element length + CRLF
		body := make([]byte, elementLength+2)

		// read next element + CRLF bytes
		_, err = io.ReadFull(p.reader, body)
		if err != nil {
			return err
		}

		// append the element to the list, without the trailing CRLF
		elements = append(elements, body[:len(body)-2])
	}

	for i, line := range elements {
		fmt.Printf("Row %d: %s\n", i, string(line))
	}

	return nil
}

func (p *Parser) Parse() error {

	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return nil
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		_type := line[0]

		switch _type {
		case '*':
			p.ParseArray(line)
		default:
			log.Println("TODO")
			continue
		}
	}
}
