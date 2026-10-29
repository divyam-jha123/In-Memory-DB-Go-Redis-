package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func Start(address string, store *store.Store, handlers map[string]commands.Handler) error {
	fmt.Println("connecting...")

    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }

    defer listener.Close()

    fmt.Println("Redis server listening on", address)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }

		// goroutine -> handle multiple clients
        go handleConnection(conn, store, handlers)
    }
}

func handleConnection(conn net.Conn, store *store.Store, handlers map[string]commands.Handler) {

	fmt.Println("New client connected:", conn.RemoteAddr())
	defer conn.Close()

	// one reader per connection, so buffered bytes are not lost between commands
	reader := bufio.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		// redis-cli sends commands as arrays: *N\r\n followed by N bulk strings
		header, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		header = strings.TrimRight(header, "\r\n")
		if len(header) == 0 || header[0] != '*' {
			conn.Write([]byte("-ERR expected array\r\n"))
			continue
		}

		count, err := strconv.Atoi(header[1:])
		if err != nil || count < 1 {
			conn.Write([]byte("-ERR invalid array length\r\n"))
			continue
		}

		fmt.Println("count:" , count)

		// []string
		args := make([]string, 0, count)
		for range count {
			arg, err := resp.ParseBulkString(reader)
			if err != nil {
				fmt.Println("Parse error:", err)
				return
			}

			args = append(args, arg)
		}

		command := strings.ToUpper(args[0])

		handler, ok := handlers[command]
		if !ok {
			writer.EncodeError("ERR unknown command '" + args[0] + "'")
			continue
		}

		if err := handler(store, writer, args); err != nil {
			fmt.Println("Write error:", err)
			return
		}

		fmt.Println(args)
	}
}