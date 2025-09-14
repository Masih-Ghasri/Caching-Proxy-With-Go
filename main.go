package main

import (
	"bufio"
	"fmt"
	"github.com/Masih-Ghasri/Caching-Proxy-With-Go.git/cache"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting TCP server:", err)
	}

	c := cache.NewCache(100)
	c.Set("exampleKey", []byte("exampleValue"))

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
			if len(parts) != 3 {
				conn.Write([]byte("Error: SET format is 'SET key value'\n"))
				continue
			}
			c.Set(parts[1], []byte(parts[2]))
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) != 2 {
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
			if len(parts) != 2 {
				conn.Write([]byte("Error: DELETE format is 'DELETE key'\n"))
				continue
			}
			if c.Delete(parts[1]) {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}

		default:
			conn.Write([]byte("Error: Unknown command '" + parts[0] + "'\n"))
		}
	}
}
