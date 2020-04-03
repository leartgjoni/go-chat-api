package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	app "github.com/leartgjoni/go-chat-api"
	"github.com/leartgjoni/go-chat-api/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebsocketHandler_Handle(t *testing.T) {
	t.Skip("temp")

	var ReadPumpArg, WritePumpArg *app.Client
	cs := &mock.ClientService{ReadPumpFn: func(client *app.Client) {ReadPumpArg = client}, WritePumpFn: func(client *app.Client) {WritePumpArg = client}}

	hub := &app.Hub{Register:   make(chan *app.Client)}
	var client *app.Client
	// read from the Register channel
	go func() {
		client = <-hub.Register
		return
	}()
	h := NewWebsocketHandler(cs, hub)

	room := "room1"
	id := "user-id"
	name := "user-name"

	// Create test server with the handler.
	s := httptest.NewServer(http.HandlerFunc(h.Handle))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	url := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, w, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?room=%s&name=%s&id=%s", url, room, name, id), nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// assert response
	if w.Status != "101 Switching Protocols" {
		t.Fatal("wrong response status", w.Status)
	}

	if w.Header["Connection"][0] != "Upgrade" {
		t.Fatal("expected Connection Header to be Upgrade, instead is: ", w.Header["Connection"][0])
	}

	if w.Header["Upgrade"][0] != "websocket" {
		t.Fatal("expected Upgrade protocol to be websocket, instead is: ", w.Header["Upgrade"][0])
	}

	//client := <-hub.Register

	if client.ID != id || client.Name != name || client.Room != room {
		t.Fatal("client is not as expected", client)
	}

	if cs.ReadPumpInvoked != true {
		t.Fatal("expected ReadPumpInvoked to be true", cs.ReadPumpInvoked)
	}

	if ReadPumpArg != client {
		t.Fatal("expect ReadPump arg to be the registered client", ReadPumpArg)
	}

	if cs.WritePumpInvoked != true {
		t.Fatal("expected WritePumpInvoked to be true")
	}

	if WritePumpArg != client {
		t.Fatal("expect WritePump arg to be the registered client", WritePumpArg)
	}
}

func TestWebsocketHandler_Handle_BadRequest(t *testing.T) {
	var cs mock.ClientService
	h := NewWebsocketHandler(&cs, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Content-Type", "application/json")
	httpHandler := http.HandlerFunc(h.Handle)
	httpHandler.ServeHTTP(w, r)

	bodyString := strings.TrimSpace(w.Body.String())

	if bodyString != "Bad Request" {
		t.Fatal("expected body to be bad request, instead is: ", bodyString)
	}
}