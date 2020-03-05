package websocket

import (
	app "github.com/leartgjoni/go-chat-api"
)

// Ensure service implements interface.
var _ app.HubService = &HubService{}

type HubService struct {
}

func NewHubService() *HubService {
	return &HubService{}
}

func (s *HubService) Run(h *app.Hub) {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.Room] == nil {
				h.Rooms[client.Room] = make(map[*app.Client]bool)
			}
			h.Rooms[client.Room][client] = true
		case client := <-h.Unregister:
			clients := h.Rooms[client.Room]
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.Rooms, client.Room)
				}
			}
		case message := <-h.Broadcast:
			clients := h.Rooms[message.Room]
			for client := range clients {
				select {
				case client.Send <- message.Data:
				default:
					// Send chan close
					close(client.Send)
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.Rooms, message.Room)
					}
				}
			}
		}
	}
}
