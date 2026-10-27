package server

import (
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
    defer conn.Close()

    fmt.Println("New client connected:", conn.RemoteAddr())
}