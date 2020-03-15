package redis

import (
	"encoding/json"
	"fmt"
	app "github.com/leartgjoni/go-chat-api"
)

// Ensure service implements interface.
var _ app.RoomService = &RoomService{}

type RoomService struct {
	db *DB
	NodeId string
}

func NewRoomService(db *DB, NodeId string) *RoomService {
	return &RoomService{
		db: db,
		NodeId: NodeId,
	}
}

func (s RoomService) PersistAndBroadcast(c *app.Client, action app.Action) {
	roomKey := fmt.Sprintf("room_%s", c.Room)

	if action == app.ActionUnregister {
		s.db.HDel(roomKey, c.ID)
	} else if action == app.ActionRegister {
		clientValue, err := json.Marshal(struct{Name string `json:"name"`}{c.Name})
		if err != nil {
			fmt.Println("err parsing json", err)
			return
		}

		s.db.HSet(roomKey, c.ID, clientValue)
	}

	clientsList, _ :=  s.db.HGetAll(roomKey).Result()

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
			Type: "user:list",
			Data: string(jsonClientsList),
			Room: c.Room,
			NodeId: s.NodeId,
		}
	}()
}

func (s RoomService) BroadcastMessage(message app.Message) {
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