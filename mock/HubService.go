package mock

import app "github.com/leartgjoni/go-chat-api"

type HubService struct {
	RunFn func(hub *app.Hub)
	RunInvoked bool

	SubscribeFn func(hub *app.Hub)
	SubscribeInvoked bool
}

func (s *HubService) Run(hub *app.Hub) {
	s.RunInvoked = true
	s.RunFn(hub)
}

func (s *HubService) Subscribe(hub *app.Hub) {
	s.SubscribeInvoked = true
	s.SubscribeFn(hub)
}


