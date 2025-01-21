package resp

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
)

func ParseArray(arrayHeader []byte, reader *bufio.Reader) ([][]byte, error) {
	// Remove CRLF of header
	arrayHeader = bytes.TrimSuffix(arrayHeader, []byte{'\r', '\n'})
	numOfExpectedElements, err := strconv.ParseInt(string(arrayHeader[1:]), 10, 64)
	if err != nil {
		return nil, nil
	}

	elements := make([][]byte, 0, numOfExpectedElements)
	for i := int64(0); i < numOfExpectedElements; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, nil
		}

		// get line length to get all digits of $ header
		lineLength := len(line)

		// read after $ until before CRLF
		elementLength, err := strconv.ParseInt(string(line[1:lineLength-2]), 10, 64)
		if err != nil {
			return nil, nil
		}

		// element length + CRLF
		body := make([]byte, elementLength+2)

		// read next element + CRLF bytes
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, nil
		}

		// append the element to the list, without the trailing CRLF
		elements = append(elements, body[:len(body)-2])
	}

	return elements, nil
}

func Parse(reader *bufio.Reader) ([][]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, nil
	}

	switch line[0] {
	case '*':
		return ParseArray(line, reader)
	default:
		log.Println("TODO")
		return nil, nil
	}
}
