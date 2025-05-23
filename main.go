package main

import (
	"fmt"
	"os"
)

func main() {
	serverIDs := []int{0, 1, 2, 3, 4, 5, 6}
	filename := "res.txt"
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	for _, id := range serverIDs {
		server := NewServer(id, filename)
		go server.Start() // Start each server in its own goroutine.
	}
	client := NewClient(filename)
	go client.Start()
	select {}
}
