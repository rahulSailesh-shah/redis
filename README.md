# Simple Redis in Go

A minimalist, in-memory Redis-compatible server written in Go. It implements the RESP protocol and supports basic commands like `PING`, `SET`/`GET`, and hash operations (`HSET`, `HGET`, `HGETALL`).

## Features

- RESP protocol parsing and serialization
- Thread-safe key-value store (`SET` / `GET`)
- Thread-safe hash map store (`HSET` / `HGET` / `HGETALL`)
- Easy to build and run with a `Makefile`

## Project Structure

- **main.go**
  Entry point: initializes the store and server, then starts listening on port 6379.
- **resp.go**
  Implements the Redis Serialization Protocol (RESP) reader and writer (`NewResp`, `Read`, `Value.Marshall()`, etc.).
- **server.go**
  TCP server that accepts connections, parses RESP arrays into commands, dispatches to handlers, and writes RESP responses.
- **store.go**
  In-memory data store with `sync.RWMutex` for concurrent access; supports string keys and hash maps.
- **Makefile**
  Helper targets for building, running, and cleaning the project.

## Prerequisites

- Go 1.18 or later installed
- `make` (optional, for convenience)

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/simple-redis-go.git
   cd simple-redis-go
   ```

2. **Build**

   ```bash
   make build
   # or:
   go build -o simple-redis main.go server.go resp.go store.go
   ```

3. **Run**
   ```bash
   make run
   # or:
   ./simple-redis
   ```
   The server will listen on port `6379`.

## Usage

You can interact using `redis-cli`, `telnet`, or any RESP-compatible client.

```bash
# Using redis-cli
redis-cli -p 6379 PING
# → PONG

redis-cli -p 6379 SET foo bar
# → OK

redis-cli -p 6379 GET foo
# → "bar"

redis-cli -p 6379 HSET myhash field1 value1
# → OK

redis-cli -p 6379 HGET myhash field1
# → "value1"

redis-cli -p 6379 HGETALL myhash
# → 1) "field1" 2) "value1"
```

Or with `telnet`:

```bash
printf "*1\r\n$4\r\nPING\r\n" | telnet localhost 6379
```

## Makefile Targets

- `make build` – compile the binary
- `make run` – build (if needed) and run the server
- `make clean` – remove built artifacts

(_Adjust the Makefile as needed for your environment._)
