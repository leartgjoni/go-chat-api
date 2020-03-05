package http

import (
	app "github.com/leartgjoni/go-chat-api"
	"net"
	"net/http"
)

type Server struct {
	ln net.Listener

	// Services
	ClientService app.ClientService
	HubService    app.HubService

	// Handlers
	websocketHandler WebsocketHandler

	// Server options.
	Addr string // bind address
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	go http.Serve(s.ln, s.router())

	return nil
}

func (s *Server) Start() error {
	s.initializeHandlers()

	return s.Open()
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// initialize handlers server needs
func (s *Server) initializeHandlers() {
	hub := &app.Hub{
		Register:   make(chan *app.Client),
		Unregister: make(chan *app.Client),
		Rooms:      make(map[string]map[*app.Client]bool),
		Broadcast:  make(chan app.Message),
	}

	go s.HubService.Run(hub)

	s.websocketHandler = NewWebsocketHandler(s.ClientService, hub)
}

// handlePing handles health check from kubernetes.
func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("healthy"))
}
