package main

import (
	"log"
)

func main() {
	store := NewStore()
	server := NewServer(":6379", store)

	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
