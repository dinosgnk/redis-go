package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const defaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	listener   net.Listener
	cmdHandler *CommandHandler
	cmdArgsCh  chan [][]byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:     cfg,
		cmdHandler: NewCommandHandler(),
		cmdArgsCh:  make(chan [][]byte),
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", server.ListenAddr)
	if err != nil {
		fmt.Printf("Error starting Redis-Go: %v\n", err)
	}

	server.listener = listener
	defer server.listener.Close()

	go server.handleCommandLoop()

	log.Println(fmt.Sprintf("Redis-Go started, listening on %s", server.ListenAddr))

	server.acceptLoop()
}

func (server *Server) acceptLoop() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println(fmt.Sprintf("New connection from %s", conn.RemoteAddr()))

	resp := NewParser(conn)
	resp.Parse(server.cmdArgsCh)
}

func (server *Server) handleCommandLoop() {
	var cmdArgs [][]byte
	for {
		cmdArgs = <-server.cmdArgsCh
		cmd := strings.ToUpper(string(cmdArgs[0]))
		switch cmd {
		case "SET":
			server.cmdHandler.Set(cmdArgs)
		case "GET":
			server.cmdHandler.Get(cmdArgs)
		default:
			log.Println("Unknown command")
		}
	}
}

func main() {

	server := NewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
