package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	app "github.com/leartgjoni/go-chat-api"
	"net/http"
)

type WebsocketHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type websocketHandler struct {
	ClientService app.ClientService
	Hub           *app.Hub
}

func NewWebsocketHandler(cs app.ClientService, hub *app.Hub) *websocketHandler {
	return &websocketHandler{ClientService: cs, Hub: hub}
}

func (h *websocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	userName := r.URL.Query().Get("name")
	userId := r.URL.Query().Get("id")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("UPGRADE ERR: %+v\n", err)
		return
	}

	client := &app.Client{
		ID:   userId,
		Name: userName,
		Room: room,
		Conn: conn,
		Hub:  h.Hub,
		Send: make(chan []byte, 256),
	}

	h.Hub.Register <- client

	go h.ClientService.ReadPump(client)
	go h.ClientService.WritePump(client)
}
