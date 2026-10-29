package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// every command function has this shape, so the server can call any of them the same way
type Handler func(s *store.Store, w *resp.Writer, args []string) error

// PING -> +PONG
func Ping(s *store.Store, w *resp.Writer, args []string) error {
	return w.EncodeSimpleString("PONG")
}

// SET key value -> +OK
func Set(s *store.Store, w *resp.Writer, args []string) error {
	if len(args) != 3 {
		return w.EncodeError("ERR wrong number of arguments for 'set' command")
	}

	s.Set(args[1], args[2])

	return w.EncodeSimpleString("OK")
}

// GET key -> bulk string, or null if the key does not exist
func Get(s *store.Store, w *resp.Writer, args []string) error {
	if len(args) != 2 {
		return w.EncodeError("ERR wrong number of arguments for 'get' command")
	}

	value, exists := s.Get(args[1])
	if !exists {
		return w.EncodeNull()
	}

	return w.EncodeBulkString(value)
}

func Del(s *store.Store, w *resp.Writer, args []string) error {
	if len(args) != 2 {
		return w.EncodeError("ERR wrong number of arguments for 'del' command")
	}

	s.Delete(args[1])

	return w.EncodeSimpleString("OK")
}
