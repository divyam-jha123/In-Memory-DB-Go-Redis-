package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// read the comamnds sent by the redis client
// IF redis-cli -p 6374 SET name "divyam"
// *3\r\n
// $3\r\n
// SET\r\n
// $4\r\n
// name\r\n
// $6\r\n
// divyam\r\n

// we will implement 5 basic RESP2 data types
// + - simple string
// "-" - error
// : - integer
// $ - bulk sring
// * - Array

func ParseBulkString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	line = strings.TrimRight(line, "\r\n")
	if len(line) == 0 || line[0] != '$' {
		return "", fmt.Errorf("expected '$', got %q", line)
	}

	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", fmt.Errorf("invalid bulk string length: %q", line[1:])
	}

	// $-1\r\n -> null bulk string
	if length < 0 {
		return "", nil
	}

	// data + \r\n
	buf := make([]byte, length+2)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return "", err
	}

	if buf[length] != '\r' || buf[length+1] != '\n' {
		return "", fmt.Errorf("bulk string not terminated by CRLF")
	}

	return string(buf[:length]), nil
}
