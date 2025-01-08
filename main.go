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
	cmdRouter *CommandRouter
	cmdCh     chan *Command
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		cmdRouter: NewCommandRouter(),
		cmdCh:     make(chan *Command),
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
	client := NewClient(conn)

	resp := NewParser(client.conn)
	resp.Parse(client, server.cmdCh)
}

func (server *Server) handleCommandLoop() {
	var cmd *Command
	for {
		cmd = <-server.cmdCh
		go server.cmdRouter.Dispatch(cmd)
	}
}

func main() {

	server := NewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
