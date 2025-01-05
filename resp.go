package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type RESPParser struct {
	reader *bufio.Reader
}

func NewRESPParser(rd io.Reader) *RESPParser {
	return &RESPParser{
		reader: bufio.NewReader(rd),
	}
}

func (r *RESPParser) Parse() {
	var position int

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			break
		}
		fmt.Println(b, int(b))
		switch b {
		case '*':
			position = 0
			fmt.Println("position set to 0")
		case '\r':
			position++
		case '\n':
			position++
		default:
			if position == 1 {
				// expectedNumOfArguments = int(b)
				expectedNumOfArguments, _ := strconv.ParseInt(string(b), 10, 64)
				fmt.Println(int(expectedNumOfArguments))
			}

			fmt.Println("dLALALALALALALALALA")
			position++
		}
	}
}

// func (r *RESPParser) readLine() {
// 	for {
// 		b, err
// 	}
// }

// func (r *RESPParser) readNumberOfElements() {
// 	line
// }

// func (r *RESPParser) readArray() {
// 	for {
// 		b, err := r.reader.ReadByte()
// 		if err != nil {
// 			break
// 		}

// 		fmt.Print(string(b))

// 		if b == '\n' {
// 			break
// 		}
// 	}
// }
