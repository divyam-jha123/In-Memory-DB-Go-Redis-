package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := 6370
	flag.Parse()

	fmt.Println("connecting...")

	l, err := net.Listen("tcp", "0.0.0:6370")
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
		os.Exit(1)
	} else {
		fmt.Printf("TCP connection binded to port %d\n", port)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)

	for {
		_, err := conn.Read(buff)
		if err != nil {
			// io.EOF means the client disconnected
			fmt.Println("client disconnected:", err)
			return
		}

		fmt.Println("msg got from the client")

		conn.Write([]byte("+PONG\r\n"))
	}
}
