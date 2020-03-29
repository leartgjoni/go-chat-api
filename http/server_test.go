package http

import (
	"fmt"
	app "github.com/leartgjoni/go-chat-api"
	"github.com/leartgjoni/go-chat-api/mock"
	httpMock "github.com/leartgjoni/go-chat-api/mock/http"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerListeningIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	fmt.Println("SERVER IS:")

	server := NewServer()
	server.Addr = ":1234"
	//server.ClientService = &mock.ClientService{}
	var hs mock.HubService
	hs.RunFn = func(hub *app.Hub) {}
	hs.SubscribeFn = func(hub *app.Hub) {}
	server.HubService = &hs


	if err := server.Start(); err != nil {
		t.Fatal("Error on server.Open()", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			t.Errorf("error closing server: %s", err)
		}
	}()

	if hs.RunInvoked != true {
		t.Fatalf("expected RunInvoked to be true")
	}

	if hs.SubscribeInvoked != true {
		t.Fatalf("expected SubscribeInvoked to be true")
	}

	resp, err := http.Get("http://localhost:1234/health")
	if err != nil {
		t.Fatal("http get failed", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("http get failed", err)
	}

	if string(body) != "healthy" {
		t.Fatalf("Expected 'healthy' but got %s", string(body))
	}
}

func TestServerRoutes(t *testing.T) {
	var tests = []struct {
		method          string
		route           string
		expectedInvoked []string
	}{
		{
			"POST",
			"/ws",
			[]string{"WebsocketHandler.Handle"},
		},
		{
			"POST",
			"/ws",
			[]string{"WebsocketHandler.Handle"},
		},
	}

	for _, test := range tests {
		server := NewServer()
		invoked := &[]string{}
		// mock handlers
		server.websocketHandler = httpMock.NewMockWebsocketHandler(invoked)

		router := server.router()

		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.route, nil)
		if err != nil {
			t.Fatal("creating request failed", err)
		}

		router.ServeHTTP(w, r)

		if !reflect.DeepEqual(*invoked, test.expectedInvoked) {
			t.Errorf("Expect %s but got %v", test.expectedInvoked, *invoked)
		}
	}
}
