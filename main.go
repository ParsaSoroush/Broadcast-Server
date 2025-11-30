package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]bool)
	clientsMu sync.Mutex
)

func broadcast(sender net.Conn, msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	
	for client := range clients {
		if client != sender {
			_, err := client.Write([]byte(msg))
			if err != nil {
				delete(clients, client)
				client.Close()
			}
		}
	}
}

func handleClient(conn net.Conn) {
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
	}()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Printf("Message from %s: %s", conn.RemoteAddr(), msg)
		broadcast(conn, msg)
	}
}

func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()
	
	fmt.Printf("Broadcast server listening on port %s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
		go handleClient(conn)
	}
}

// Client functions
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nDisconnected from server")
			os.Exit(0)
		}
		fmt.Printf(">> %s", msg)
	}
}

func connectClient(host string, port string) {
	addr := host + ":" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to broadcast server at %s\n", addr)
	fmt.Println("Type your messages and press Enter to send. Press Ctrl+C to exit.")

	go receiveMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			continue
		}
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return
		}
	}
}

func printUsage() {
	fmt.Println("Broadcast Server CLI Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  broadcast-server start [port]           Start the broadcast server (default port: 8080)")
	fmt.Println("  broadcast-server connect [host] [port]  Connect to a broadcast server (default: localhost:8080)")
	fmt.Println("\nExamples:")
	fmt.Println("  broadcast-server start")
	fmt.Println("  broadcast-server start 9000")
	fmt.Println("  broadcast-server connect")
	fmt.Println("  broadcast-server connect localhost 9000")
	fmt.Println("  broadcast-server connect 192.168.1.100 8080")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "start":
		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		startServer(port)

	case "connect":
		host := "localhost"
		port := "8080"
		
		if len(os.Args) > 2 {
			host = os.Args[2]
		}
		if len(os.Args) > 3 {
			port = os.Args[3]
		}
		
		connectClient(host, port)

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}