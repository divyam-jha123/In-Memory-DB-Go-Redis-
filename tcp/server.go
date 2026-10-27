package server

import (
	"bytes"
	"fmt"
	"net"
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

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

  // Process and use the data (here, we'll just print it)
	data := buffer[:n]

	if  bytes.Contains(bytes.ToUpper(data), []byte("PING")) {
		conn.Write([]byte("+PONG\r\n"))
	}

    fmt.Println("New client connected:", conn.RemoteAddr())
}