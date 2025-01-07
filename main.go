package main

import (
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
	listener  net.Listener
	kv        *KVStore
	cmdArgsCh chan [][]byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		kv:        NewKVStore(),
		cmdArgsCh: make(chan [][]byte),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		fmt.Printf("Error starting Redis-Go: %v\n", err)
	}

	s.listener = listener
	defer s.listener.Close()

	go s.handleCommandLoop()

	log.Println(fmt.Sprintf("Redis-Go started, listening on %s", s.ListenAddr))

	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println(fmt.Sprintf("New connection from %s", conn.RemoteAddr()))

	resp := NewParser(conn)
	resp.Parse(s.cmdArgsCh)
}

func (s *Server) handleCommandLoop() {
	var cmdArgs [][]byte
	for {
		cmdArgs = <-s.cmdArgsCh
		log.Println("Received command: ", string(cmdArgs[0]))
	}
}

func main() {

	server := NewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
