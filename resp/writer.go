package resp

import (
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}	

func (w *Writer) WriteBulkString(value string) error {	

	message := "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"

	_, err := w.writer.Write([]byte(message))

    return err	
}
