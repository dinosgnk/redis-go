package main

import (
	"bufio"
	"fmt"
	"log"
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

func NewServer(cfg Config) *Server {
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
	log.Println(fmt.Sprintf("Redis-Go started, listening on %s", s.ListenAddr))
	defer s.listener.Close()
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(fmt.Sprintf("Error accepting connection: %v\n", err))
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	log.Println(fmt.Sprintf("New connection from %s", conn.RemoteAddr()))

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(fmt.Sprintf("Error writing to connection: %v\n", err))
		}

		//_, err := conn.Write([]byte("+OK\r\n"))
		//conn.Write([]byte("+OK\r\n"))
		conn.Write([]byte(msg))
	}
}

func main() {

	server := NewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
