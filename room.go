package app

type Action int

const (
	ActionRegister Action = 0
	ActionUnregister Action = 1
)

type RoomService interface {
	PersistAndBroadcast(c *Client, a Action)
	BroadcastMessage(m Message)
}

type PubSubService interface {
	Subscribe(hub *Hub)
}