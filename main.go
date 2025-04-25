package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type HandlerFunc func([]Value) Value

func main() {
	store := NewStore()
	var handlers = map[string]HandlerFunc{
		"PING": store.ping,
		"SET":  store.set,
		"GET":  store.get,
		"HSET": store.hset,
		"HGET": store.hget,
	}

	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Println("Error starting redis server: ", err.Error())
		os.Exit(1)
	}

	log.Println("Redis Server running on PORT: 6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleClient(conn, handlers)
	}

}

func handleClient(conn net.Conn, handlers map[string]HandlerFunc) {
	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error reading from client: ", err.Error())
			return
		}

		if value.typ != "array" {
			log.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			log.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := handlers[command]
		if !ok {
			log.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
