package main

import (
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	server "github.com/codecrafters-io/redis-starter-go/internal/server"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := ":6374"

	// one store shared by all clients
	db := store.NewStore()

	// command name -> function that handles it
	handlers := map[string]commands.Handler{
		"PING": commands.Ping,
		"SET":  commands.Set,
		"GET":  commands.Get,
		"DEL": commands.Del,
	}

	err := server.Start(port, db, handlers)

	if err != nil {
		os.Exit(1)
	}
}
