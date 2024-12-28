package main

import (
	"bufio"
	"fmt"
	"net"
)

const defaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	listener net.Listener
}

func CreateNewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{Config: cfg}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		fmt.Printf("Error starting Redis-Go: %v\n", err)
	}
	s.listener = listener
	fmt.Printf("Redis-Go started. Listening on %s\n", s.ListenAddr)
	defer s.listener.Close()
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	message := "Welcome to Redis-Go!\n\n"
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error writing to connection: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		messageReceived := scanner.Text()
		if messageReceived == "exit" {
			fmt.Printf("Client %s disconnected.\n", conn.RemoteAddr())
			return
		}
		fmt.Printf("Received: %s\n", messageReceived)
	}
}

func main() {

	server := CreateNewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
