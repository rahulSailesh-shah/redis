package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type HandlerFunc func([]Value) Value

type Server struct {
	addr     string
	handlers map[string]HandlerFunc
}

func NewServer(addr string, store *Store) *Server {
	return &Server{
		addr: addr,
		handlers: map[string]HandlerFunc{
			"PING":    store.ping,
			"SET":     store.set,
			"GET":     store.get,
			"HSET":    store.hset,
			"HGET":    store.hget,
			"HGETALL": store.hgetall,
		},
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer ln.Close()
	log.Println("Redis Server running on PORT: 6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	writer := NewWriter(conn)

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("read error: %v", err)
			writer.Write(Value{typ: "error", str: "ERR internal error"})
			return
		}

		cmd, args, err := parseRequest(value)
		if err != nil {
			log.Printf("bad request: %v", err)
			writer.Write(Value{typ: "error", str: fmt.Sprintf("ERR %v", err)})
			return
		}

		handler, ok := s.handlers[cmd]
		if !ok {
			writer.Write(Value{typ: "error", str: fmt.Sprintf("ERR unknown command '%s'", cmd)})
			continue
		}
		result := handler(args)
		writer.Write(result)
	}
}

func parseRequest(v Value) (cmd string, args []Value, err error) {
	if v.typ != "array" {
		return "", nil, fmt.Errorf("expected array, got %s", v.typ)
	}

	if len(v.array) == 0 {
		return "", nil, fmt.Errorf("empty array")
	}

	cmd = strings.ToUpper(v.array[0].bulk)
	return cmd, v.array[1:], nil
}
