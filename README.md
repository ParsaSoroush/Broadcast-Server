# Broadcast Server

A simple TCP broadcast server and client written in Go. Allows multiple clients to connect, send messages, and broadcast them to all connected clients in real-time.

[🔗GitHub Repository](https://github.com/ParsaSoroush/Broadcast-Server.git)  
[🔗Project Page](https://roadmap.sh/projects/broadcast-server)

---

## Features

- Start a TCP broadcast server on any port
- Connect multiple clients to the server
- Real-time message broadcasting between clients
- Simple CLI interface for server and client

---

## Installation

Make sure you have [Go](https://golang.org/dl/) installed.

```bash
git clone https://github.com/ParsaSoroush/Broadcast-Server.git
cd Broadcast-Server
go run main.go
```

## Usage

### Start the Server
```bash
# Default port 8080
go run main.go start

# Custom port
go run main.go start 9000
```

## Example
```bash
Server:
Broadcast server listening on port 8080
Client connected: 127.0.0.1:50000

Client 1:
>> Hello from Client 1

Client 2:
>> Hello from Client 1
```