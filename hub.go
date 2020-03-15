package app

type Message struct {
	UserID string `json:"userId"`
	Type string `json:"type"`
	Data string `json:"data"`
	Room string `json:"room"`
	NodeId string `json:"nodeId"`
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