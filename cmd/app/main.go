package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/leartgjoni/go-chat-api/http"
	"github.com/leartgjoni/go-chat-api/redis"
	"github.com/leartgjoni/go-chat-api/websocket"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/signal"
)

func main() {
	m := NewMain()

	// Load configuration.
	if err := m.LoadConfig(); err != nil {
		_, _ = fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Execute program.
	if err := m.Run(); err != nil {
		_, _ = fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	_, _ = fmt.Fprintln(m.Stdout, "received interrupt, shutting down...")
	_ = m.Close()
}

// Main represents the main program execution.
type Main struct {
	NodeId     string // represents this process
	ConfigPath string
	Config     Config

	// Input/output streams
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	closeCh chan int
	closeFn func() error
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		NodeId: uuid.New().String(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		closeCh: make(chan int),
		closeFn: func() error { return nil },
	}
}

// Close cleans up the program.
func (m *Main) Close() error { return m.closeFn() }

// LoadConfig parses the configuration file.
func (m *Main) LoadConfig() error {

	if os.Getenv("CONFIG_PATH") != "" {
		m.ConfigPath = os.Getenv("CONFIG_PATH")
	} else {
		m.ConfigPath = ".env"
	}

	viper.SetConfigFile(m.ConfigPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("no config file present. will default to env variables")
	}

	m.Config = Config{
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		Port:          viper.GetString("PORT"),
	}

	if m.Config.RedisAddr == "" || m.Config.Port == "" {
		return errors.New("missing configs")
	}

	return nil
}

func (m *Main) Run() error {
	redisDb, err := redis.Open(m.Config.RedisAddr, m.Config.RedisPassword)
	if err != nil {
		fmt.Println(m.Stderr, err)
		os.Exit(1)
	}

	roomService := redis.NewRoomService(redisDb, m.NodeId, m.closeCh)

	clientService := websocket.NewClientService(m.NodeId)
	hubService := websocket.NewHubService(roomService)

	// Initialize Http server.
	httpServer := http.NewServer()
	httpServer.Addr = fmt.Sprintf(":%s", m.Config.Port)

	httpServer.ClientService = clientService
	httpServer.HubService = hubService

	// Start HTTP server.
	if err := httpServer.Start(); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(m.Stdout, "Listening on port: %s\n", httpServer.Addr)

	// Assign close function.
	m.closeFn = func() error {
		m.closeCh <- 0
		_ = httpServer.Close()
		_ = redisDb.Close()
		return nil
	}

	return nil
}

type Config struct {
	RedisAddr     string
	RedisPassword string
	Port          string
}
