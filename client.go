package app

import "github.com/gorilla/websocket"

type Client struct {
	ID   string
	Room string
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
}

type ClientService interface {
	ReadPump(client *Client)
	WritePump(client *Client)
}
