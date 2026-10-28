package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Start(address string) error {
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
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {

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

		if(count == 1 && args[0] == "PING"){
			conn.Write([]byte("+PONG\r\n"))
		}

		if(count == 3 && args[0] == "SET"){
			err := writer.WriteBulkString("OK")


			if(err != nil) {
				fmt.Println("error while writing the resp encoding")
			}
		}

		fmt.Println(args)
	}
}