package app

type Action int

const (
	ActionRegister   Action = 0
	ActionUnregister Action = 2
)

type RoomService interface {
	Update(c *Client, a Action)
	Publish(m Message)
	Subscribe(hub *Hub)
}
