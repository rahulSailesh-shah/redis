package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting redis server: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Redis Server running on PORT: 6379")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	handleClient(conn)
}

func handleClient(conn net.Conn) {
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(value)

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
