BINARY = redis-server-go

.PHONY: all build run clean

all: build

build: main.go resp.go
	go build -o $(BINARY) main.go resp.go store.go

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
