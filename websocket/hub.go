package websocket

import (
	"encoding/json"
	"fmt"
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

			// send everybody the list of all users
			type StaticUser struct {
				Name string `json:"name"`
				Id string `json:"id"`
			}
			var clients []StaticUser
			for client := range h.Rooms[client.Room] {
				clients = append(clients, StaticUser{Name: client.Name, Id: client.ID})
			}
			fmt.Println("all users", clients)

			jsonClients, err := json.Marshal(clients)
			if err != nil {
				fmt.Println("err parsing json", err)
				return
			}

			// sending to channel is blocking until the channel reads, but the channel can read only if the loop continues.
			// so we send to the channel from another goroutine, allowing the loop to continue
			go func() {
				client.Hub.Broadcast <- app.Message{
					Type: "user-update",
					Data: string(jsonClients),
					Room: client.Room,
				}
			}()
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
			fmt.Print("here", message)
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
