package mock

import "net/http"

type WebsocketHandler struct {
	Invoked *[]string
}

func NewMockWebsocketHandler(invoked *[]string) *WebsocketHandler {
	return &WebsocketHandler{invoked}
}

func (h *WebsocketHandler) Handle(http.ResponseWriter, *http.Request) {
	*h.Invoked = append(*h.Invoked, "WebsocketHandler.Handle")
}
