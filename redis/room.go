package redis

import (
	"encoding/json"
	"fmt"
	app "github.com/leartgjoni/go-chat-api"
)

// Ensure service implements interface.
var _ app.RoomService = &RoomService{}

type RoomService struct {
	db     *DB
	NodeId string
}

func NewRoomService(db *DB, NodeId string) *RoomService {
	return &RoomService{
		db:     db,
		NodeId: NodeId,
	}
}

// Update updates the list of users in a room stored in redis and broadcasts the updated list
func (s RoomService) Update(c *app.Client, action app.Action) {
	roomKey := fmt.Sprintf("room_%s", c.Room)

	if action == app.ActionUnregister {
		s.db.HDel(roomKey, c.ID)
	} else if action == app.ActionRegister {
		clientValue, err := json.Marshal(struct {
			Name string `json:"name"`
		}{c.Name})
		if err != nil {
			fmt.Println("err parsing json", err)
			return
		}

		s.db.HSet(roomKey, c.ID, clientValue)
	}

	clientsList, _ := s.db.HGetAll(roomKey).Result()

	jsonClientsList, err := json.Marshal(clientsList)
	if err != nil {
		fmt.Println("err parsing json", err)
		return
	}

	// sending to channel is blocking until the channel reads, but the channel can read only if the loop continues.
	// so we send to the channel from another goroutine, allowing the loop to continue
	go func() {
		// no client because we want to send to everybody in the room, include the client
		c.Hub.Broadcast <- app.Message{
			Type:   "user:list",
			Data:   string(jsonClientsList),
			Room:   c.Room,
			NodeId: s.NodeId,
		}
	}()
}

// Publish publishes a message to the redis `room-messages` channel
func (s RoomService) Publish(message app.Message) {
	if message.NodeId != s.NodeId {
		return
	}

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		fmt.Println("err parsing json", err)
		return
	}

	s.db.Publish("room-messages", string(jsonMessage))
}

// Subscribe subscribes to the `room-messages` channel
func (s RoomService) Subscribe(hub *app.Hub) {
	pubsub := s.db.Subscribe("room-messages")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = pubsub.Close()
	}()

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Consume messages.
	for msg := range ch {
		message := app.Message{}
		json.Unmarshal([]byte(msg.Payload), &message)

		if message.NodeId != s.NodeId {
			hub.Broadcast <- message
		}
	}
}
