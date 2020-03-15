package websocket

import (
	"encoding/json"
	"fmt"
	app "github.com/leartgjoni/go-chat-api"
)

// Ensure service implements interface.
var _ app.HubService = &HubService{}

type HubService struct {
	RoomService app.RoomService
}

func NewHubService(rs app.RoomService) *HubService {
	return &HubService{
		RoomService: rs,
	}
}

func (s *HubService) Run(h *app.Hub) {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.Room] == nil {
				h.Rooms[client.Room] = make(map[*app.Client]bool)
			}
			h.Rooms[client.Room][client] = true

			s.RoomService.PersistAndBroadcast(client, app.ActionRegister)
		case client := <-h.Unregister:
			clients := h.Rooms[client.Room]
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.Rooms, client.Room)
				}

				s.RoomService.PersistAndBroadcast(client, app.ActionUnregister)
			}
		case message := <-h.Broadcast:
			fmt.Println("broadcast")
			s.RoomService.BroadcastMessage(message)

			clients := h.Rooms[message.Room]
			for client := range clients {
				if message.UserID == client.ID {
					continue
				}

				jsonMsg, err := json.Marshal(message)
				if err != nil {
					fmt.Println("err parsing json", err)
					return
				}

				select {
				case client.Send <- jsonMsg:
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
