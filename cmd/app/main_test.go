package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	_ = os.Setenv("CONFIG_PATH", "../../test.env")

	m := NewMain()
	if err := m.LoadConfig(); err != nil {
		t.Fatal("cannot load config", err)
	}
	if err := m.Run(); err != nil {
		t.Fatal("cannot run main", err)
	}
	defer func() {
		if err := m.Close(); err != nil {
			t.Fatal("cannot close main", err)
		}
	}()

	resp, err := http.Get("http://localhost:9090/health")
	if err != nil {
		t.Fatal("http get failed", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal("cannot close resp body", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("http get failed", err)
	}

	if string(body) != "healthy" {
		t.Fatalf("Expected 'healthy' but got %s", string(body))
	}
}
