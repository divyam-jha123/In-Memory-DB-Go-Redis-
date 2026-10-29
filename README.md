![progress-banner](https://backend.codecrafters.io/progress/redis/95dc145e-410d-4e0c-b2df-95995c316da4)

# Redis-go

A lightweight, Redis-compatible in-memory database server written from scratch in Go.

Redis-go speaks the Redis wire protocol over TCP, so you can talk to it with standard Redis clients such as `redis-cli`. It's built to explore how Redis works under the hood: networking, the Redis Serialization Protocol (RESP), concurrent client handling and a thread-safe key-value store.

## Features

- **Redis-compatible:** works with `redis-cli` and other Redis clients
- **Concurrent:** each client connection runs in its own goroutine, sharing one mutex-protected store
- **Easy to extend:** adding a command is one function and one line in `main.go`
- **Dependency-free:** built only on the Go standard library
- **Small:** easy to read and learn from

## Quick start

### Prerequisites

- [Go](https://go.dev/dl/) 1.27 or later

### Run the server

```sh
git clone <your-repo-url>
cd Redis-go
make dev
```

or:

```sh
go run ./cmd/redis-server
```

### Connect

```sh
redis-cli -p 6374 PING
# PONG

redis-cli -p 6374 SET name redis
# OK

redis-cli -p 6374 GET name
# "redis"
```

Or without `redis-cli`:

```sh
printf '*1\r\n$4\r\nPING\r\n' | nc localhost 6374
# +PONG
```

## Supported commands

| Command          | Description                                                  |
| ---------------- | ------------------------------------------------------------ |
| `PING`           | Checks the server is alive and replies `PONG`                |
| `SET key value`  | Stores `value` under `key` and replies `OK`                  |
| `GET key`        | Returns the value stored at `key`, or `(nil)` if it doesn't exist |
| `DEL key`        | Removes `key` and replies `OK`                               |

Calling a command with the wrong number of arguments returns `ERR wrong number of arguments`, and an unrecognised command returns `ERR unknown command`.

## How it works

```
client ──TCP──▶ net.Listener ──Accept()──▶ goroutine per connection
                                               │
                                               ├─ parse RESP array  →  ["SET", "name", "redis"]
                                               ├─ look up handlers["SET"]
                                               ├─ handler reads/writes the shared store
                                               └─ write RESP reply (+OK\r\n)
```

1. `main.go` creates a single `store.Store` and a map from command names to handler functions, then passes both to `server.Start`.
2. The server opens a TCP listener with `net.Listen` and accepts connections in a loop.
3. Each connection is handled in its own goroutine, so one slow client never blocks the others.
4. Requests are parsed from RESP arrays of bulk strings, the format `redis-cli` sends.
5. The command name is upper-cased and looked up in the handler map. The matching handler checks its arguments, reads or writes the store, and encodes a RESP reply.
6. The store guards its map with a `sync.Mutex`, so concurrent clients can safely read and write the same keys.

## Adding a command

Every command handler has the same signature:

```go
type Handler func(s *store.Store, w *resp.Writer, args []string) error
```

To add one:

1. Write the function in `internal/commands/implementation.go`. `args[0]` is the command name, and the rest are its arguments.
2. Register it in the `handlers` map in `cmd/redis-server/main.go`:

   ```go
   handlers := map[string]commands.Handler{
       "PING": commands.Ping,
       "SET":  commands.Set,
       "GET":  commands.Get,
       "DEL":  commands.Del,
       // "NEW": commands.New,
   }
   ```

## Project structure

```
.
├── cmd/
│   └── redis-server/
│       └── main.go            # Entry point: creates the store, registers commands, starts the server
├── internal/
│   ├── commands/
│   │   └── implementation.go  # Handler type and command implementations (PING, SET, GET, DEL)
│   ├── resp/
│   │   ├── parser.go          # Reads RESP bulk strings from the client
│   │   └── writer.go          # Encodes RESP replies (simple strings, errors, bulk strings, null, arrays)
│   ├── server/
│   │   └── server.go          # TCP listener, per-connection goroutines, command dispatch
│   └── store/
│       └── store.go           # Thread-safe in-memory key-value store
├── docs/
│   └── commands/              # Command metadata in JSON (not used by the server yet)
├── go.mod
└── Makefile                   # `make dev` builds and runs the server
```

## Learn more

- [Redis protocol specification (RESP)](https://redis.io/docs/latest/develop/reference/protocol-spec/)
- [Redis commands reference](https://redis.io/docs/latest/commands/)
