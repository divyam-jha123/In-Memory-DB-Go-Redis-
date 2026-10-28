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

		// goroutine -> handle multiple clients
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {

	fmt.Println("New client connected:", conn.RemoteAddr())
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
	
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		data := buffer[:n]

		if  bytes.Contains(bytes.ToUpper(data), []byte("PING")) {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}