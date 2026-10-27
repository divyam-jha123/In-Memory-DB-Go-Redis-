package main

import (
	"os"

	server "github.com/codecrafters-io/redis-starter-go/tcp"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := ":6374"

	err := server.Start(port)

	if err == nil {
		os.Exit(1)
	}
}