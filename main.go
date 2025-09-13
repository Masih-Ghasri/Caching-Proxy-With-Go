package main

import (
	"bufio"
	"fmt"
	"github.com/Masih-Ghasri/Caching-Proxy-With-Go.git/cache"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting TCP server:", err)
	}

	c := cache.NewCache()
	c.Set("exampleKey", []byte("exampleValue"), 5*60*1000000000) // 5 minutes

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, c)
	}
}

func handleConnection(conn net.Conn, c *cache.Cache) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected or error reading:", err)
			return
		}

		message = strings.TrimSpace(message)
		parts := strings.Fields(message)
		if len(parts) == 0 {
			continue
		}
		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("Error: SET format is 'SET key value [duration_seconds]'\n"))
				continue
			}

			duration := 24 * time.Hour // Default duration

			if len(parts) == 4 {
				seconds, err := strconv.Atoi(parts[3])
				if err != nil {
					conn.Write([]byte("Error: Invalid duration, must be a number in seconds\n"))
					continue
				}
				duration = time.Duration(seconds) * time.Second
			}

			c.Set(parts[1], []byte(parts[2]), duration)
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("Error: GET format is 'GET key'\n"))
				continue
			}
			key := parts[1]
			value, exists := c.Get(key)
			if exists {
				conn.Write(append(value, '\n'))
			} else {
				conn.Write([]byte("Key not found\n"))
			}

		case "DELETE":
			if len(parts) < 2 {
				conn.Write([]byte("Error: DELETE format is 'DELETE key'\n"))
				continue
			}
			if c.Delete(parts[1]) {
				conn.Write([]byte("1\n")) // 1 for success
			} else {
				conn.Write([]byte("0\n")) // 0 for not found
			}

		default:
			conn.Write([]byte("Error: Unknown command '" + parts[0] + "'\n"))
		}
	}
}
