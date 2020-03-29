package mock

import app "github.com/leartgjoni/go-chat-api"

type ClientService struct {
	ReadPumpFn func(client *app.Client)
	ReadPumpInvoked bool

	WritePumpFn func(hub *app.Client)
	WritePumpInvoked bool
}

func (s *ClientService) ReadPump(client *app.Client) {
	s.ReadPumpInvoked = true
	s.ReadPumpFn(client)
}

func (s *ClientService) WritePump(client *app.Client) {
	s.WritePumpInvoked = true
	s.WritePumpFn(client)
}



