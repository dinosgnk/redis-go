package tcp

import (
	"fmt"
	"log"
	"net"
	"redis-go/commands"
	"redis-go/message"
	"redis-go/protocol"
)

const defaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	listener  net.Listener
	cmdRouter *commands.Router
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		cmdRouter: commands.NewRouter(),
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
	log.Printf("New connection from %s", conn.RemoteAddr())

	clientMsgCh := make(chan *message.Message)

	go server.handleClientCommandsLoop(clientMsgCh)

	parser := protocol.NewParser(conn)
	parser.Parse(conn, clientMsgCh)
}

func (server *Server) handleClientCommandsLoop(clientMsgCh <-chan *message.Message) {
	var msg *message.Message
	for {
		msg = <-clientMsgCh
		go server.cmdRouter.Route(msg)
	}
}
