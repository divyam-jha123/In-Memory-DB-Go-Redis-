![progress-banner](https://backend.codecrafters.io/progress/redis/95dc145e-410d-4e0c-b2df-95995c316da4)

# Redis-go

A lightweight, Redis-compatible in-memory database server written from scratch in Go.

Redis-go speaks the Redis wire protocol over TCP, so you can talk to it with standard Redis clients such as `redis-cli`. It's built to explore how Redis works under the hood: networking, the Redis Serialization Protocol (RESP), concurrent client handling and a thread-safe key-value store.

## Features

- **Redis-compatible:** works with `redis-cli` and other Redis clients
- **Concurrent:** each client connection runs in its own goroutine, sharing one mutex-protected store
- **Key expiry:** `SET key value EX seconds` or `PX milliseconds` gives a key a time to live, after which it disappears
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

redis-cli -p 6374 SET session abc EX 10
# OK   (expires after 10 seconds)

redis-cli -p 6374 SET lock 1 PX 500
# OK   (expires after 500 milliseconds)
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
| `SET key value EX seconds` | Same as `SET`, but the key expires after `seconds` seconds |
| `SET key value PX milliseconds` | Same as `SET`, but the key expires after `milliseconds` milliseconds |
| `GET key`        | Returns the value stored at `key`, or `(nil)` if it doesn't exist |
| `DEL key`        | Removes `key` and replies `OK`                               |

Calling a command with the wrong number of arguments returns `ERR wrong number of arguments`, and an unrecognised command returns `ERR unknown command`.

### Expiry

```sh
redis-cli -p 6374 SET token xyz PX 200
# OK
redis-cli -p 6374 GET token
# "xyz"
sleep 0.3
redis-cli -p 6374 GET token
# (nil)

redis-cli -p 6374 SET token xyz XX 5
# (error) ERR syntax error
redis-cli -p 6374 SET token xyz PX 0
# (error) ERR invalid expire time in 'set' command
```

- `EX` takes a whole number of seconds greater than zero. `EX 60` is one minute and `EX 1000` is about 16 minutes 40 seconds.
- `PX` takes a whole number of milliseconds greater than zero, for TTLs shorter or more precise than a second. `PX 1000` is the same as `EX 1`.
- A plain `SET` on an existing key clears any expiry it had.
- An invalid time returns `ERR invalid expire time in 'set' command`, and any option other than `EX` or `PX` returns `ERR syntax error`.
- Expiry is lazy: each key stores its deadline, and `GET` checks it. Once the deadline has passed, `GET` deletes the key and returns `(nil)`. An expired key that is never read again stays in memory until it is overwritten or deleted.

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
7. Each value is stored together with its expiry time (zero means no expiry), so a value and its TTL are always written in one locked step.

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
│       ├── store.go           # Thread-safe in-memory key-value store with per-key expiry
│       └── expiry.go          # SetExpiry: changes the TTL of an existing key (not exposed as a command yet)
├── docs/
│   └── commands/              # Command metadata in JSON (not used by the server yet)
├── go.mod
└── Makefile                   # `make dev` builds and runs the server
```

## Learn more

- [Redis protocol specification (RESP)](https://redis.io/docs/latest/develop/reference/protocol-spec/)
- [Redis commands reference](https://redis.io/docs/latest/commands/)
