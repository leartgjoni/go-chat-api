package redis

import (
	"encoding/json"
	"fmt"
	app "github.com/leartgjoni/go-chat-api"
)

type PubSubService struct {
	db *DB
	NodeId string
}

func NewPubSubService(db *DB, NodeId string) *PubSubService {
	return &PubSubService{
		db: db,
		NodeId: NodeId,
	}
}

func (s PubSubService) Subscribe(hub *app.Hub) {
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

	// Publish a message.
	//err = rdb.Publish("mychannel1", "hello").Err()
	//if err != nil {
	//	panic(err)
	//}

	// Consume messages.
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)

		message := app.Message{}
		json.Unmarshal([]byte(msg.Payload), &message)

		fmt.Println("MESSAGE", message)
		if message.NodeId != s.NodeId {
			hub.Broadcast <- message
		}
	}
}
