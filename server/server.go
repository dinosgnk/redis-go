package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"redis-go/commands"
	"redis-go/protocol"
)

const defaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	listener       net.Listener
	commandHandler *commands.CommandHandler
}

func NewServer(cfg Config, cmdHandler *commands.CommandHandler) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:         cfg,
		commandHandler: cmdHandler,
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", server.ListenAddr)
	if err != nil {
		fmt.Printf("Error starting Redis-Go: %v\n", err)
	}

	server.listener = listener
	defer server.listener.Close()

	log.Printf("Redis-Go started, listening on %s", server.ListenAddr)

	server.acceptConnections()
}

func (server *Server) acceptConnections() {
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
	log.Printf("New connection from %s", conn.RemoteAddr())

	clientBuf := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		cmdArgs, err := protocol.Parse(clientBuf.Reader)
		if err != nil {
			log.Printf("Connection from %v closed by client", conn.RemoteAddr())
			break
		}
		reply := server.commandHandler.Handle(cmdArgs)
		conn.Write(reply)
	}
}
