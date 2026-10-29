package resp

import (
	"io"
	"strconv"
	"strings"
)

// encodes replies sent back to the redis client
// + - simple string -> +OK\r\n
// "-" - error       -> -ERR message\r\n
// : - integer       -> :1\r\n
// $ - bulk string   -> $5\r\nhello\r\n  (null -> $-1\r\n)
// * - array         -> *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n  (null -> *-1\r\n)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) EncodeSimpleString(value string) error {
	_, err := w.writer.Write([]byte("+" + value + "\r\n"))

	return err
}

// message should start with an error prefix, e.g. "ERR unknown command"
func (w *Writer) EncodeError(message string) error {
	_, err := w.writer.Write([]byte("-" + message + "\r\n"))

	return err
}

func (w *Writer) EncodeInteger(value int) error {
	_, err := w.writer.Write([]byte(":" + strconv.Itoa(value) + "\r\n"))

	return err
}

func (w *Writer) EncodeBulkString(value string) error {
	_, err := w.writer.Write([]byte(bulkString(value)))

	return err
}

// null bulk string, e.g. GET on a missing key
func (w *Writer) EncodeNull() error {
	_, err := w.writer.Write([]byte("$-1\r\n"))

	return err
}

// array of bulk strings, e.g. KEYS or MGET replies
func (w *Writer) EncodeArray(values []string) error {
	var message strings.Builder

	message.WriteString("*")
	message.WriteString(strconv.Itoa(len(values)))
	message.WriteString("\r\n")

	for _, value := range values {
		message.WriteString(bulkString(value))
	}

	_, err := w.writer.Write([]byte(message.String()))

	return err
}

// null array
func (w *Writer) EncodeNullArray() error {
	_, err := w.writer.Write([]byte("*-1\r\n"))

	return err
}

func bulkString(value string) string {
	return "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
}
