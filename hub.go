package app

type Message struct {
	Data []byte
	Room string
}

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
}

type HubService interface {
	Run(hub *Hub)
}
