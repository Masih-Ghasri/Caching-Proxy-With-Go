package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type SafeCache struct {
	mu    sync.Mutex
	cache map[string]string
}

func handleConnection(conn net.Conn, cache *SafeCache) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(conn)

	fmt.Printf("Handling new connection from %v\n", conn.RemoteAddr())

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Printf("Message received: %s", message)

	message = strings.Trim(message, "\r\n")
	parts := strings.Fields(message)
	if len(parts) == 0 {
		conn.Write([]byte("Error: Empty command\n"))
		return
	}
	command := strings.ToUpper(parts[0])

	switch command {
	case "SET":
		if len(parts) < 3 {
			conn.Write([]byte("Error: Key or value is missing\n"))
			return
		}
		cache.mu.Lock()
		key := parts[1]
		value := parts[2]
		cache.cache[key] = value
		cache.mu.Unlock()
		write, err := conn.Write([]byte("OK\n"))
		if err != nil {
			log.Println("Error writing to connection:", err)
		}
		fmt.Printf("Wrote %d bytes to connection\n", write)
	case "GET":
		if len(parts) < 2 {
			conn.Write([]byte("Error: Key is missing\n"))
			return
		}
		cache.mu.Lock()
		key := parts[1]
		value, exists := cache.cache[key]
		cache.mu.Unlock()
		if exists {
			write, err := conn.Write([]byte(value + "\n"))
			if err != nil {
				log.Println("Error writing to connection:", err)
			}
			fmt.Printf("Wrote %d bytes to connection\n", write)
		} else {
			write, err := conn.Write([]byte("Key not found\n"))
			if err != nil {
				log.Println("Error writing to connection:", err)
			}
			fmt.Printf("Wrote %d bytes to connection\n", write)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting TCP server:", err)
	}

	safeCache := SafeCache{
		cache: make(map[string]string),
	}
	safeCache.cache["example"] = "This is a cached value"
	fmt.Println("Cache initialized with example key.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, &safeCache)
	}
}
