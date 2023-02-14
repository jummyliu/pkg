package redisbuilder

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type DBConnect struct {
	*redis.Client
	Options *redis.Options
}

// New return a new redis client, and try ping.
func New(opts *redis.Options) (*DBConnect, error) {
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &DBConnect{
		Client: client,
		Options: opts,
	}, nil
}