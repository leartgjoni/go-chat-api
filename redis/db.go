package redis

import (
	"github.com/go-redis/redis/v7"
)

type DB struct {
	*redis.Client
}

func Open(addr, password string) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,  // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &DB{client}, nil
}
