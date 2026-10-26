![progress-banner](https://backend.codecrafters.io/progress/redis/95dc145e-410d-4e0c-b2df-95995c316da4)

# Redis-go

A lightweight, Redis-compatible in-memory database server written from scratch in Go.

Redis-go speaks the Redis wire protocol over TCP, so you can talk to it with standard Redis clients such as `redis-cli`. It's built to explore how Redis works under the hood: networking, the Redis Serialization Protocol (RESP) and concurrent client handling.

## Features

- **Redis-compatible:** works with `redis-cli` and other Redis clients
- **Concurrent:** each client connection runs in its own goroutine
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
go run ./app
```

### Connect

```sh
redis-cli -p 6370 PING
# PONG
```

Or without `redis-cli`:

```sh
printf '*1\r\n$4\r\nPING\r\n' | nc localhost 6370
# +PONG
```

## Supported commands

| Command | Description                     |
| ------- | ------------------------------- |
| `PING`  | Checks the server is alive and replies `PONG` |

## How it works

```
client ──TCP──▶ net.Listener ──Accept()──▶ goroutine per connection
                                               │
                                               ├─ read request from socket
                                               └─ write RESP reply (+PONG\r\n)
```

1. The server opens a TCP listener with `net.Listen`.
2. It accepts incoming connections in a loop.
3. Each connection is handled in its own goroutine, so one slow client never blocks the others.
4. Replies are encoded in RESP, the same protocol real Redis uses.

## Project structure

```
.
├── app/
│   └── main.go     # Server entry point and connection handling
├── go.mod
└── Makefile        # `make dev` builds and runs the server
```

## Learn more

- [Redis protocol specification (RESP)](https://redis.io/docs/latest/develop/reference/protocol-spec/)
- [Redis commands reference](https://redis.io/docs/latest/commands/)
