package esbuilder

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
)

type DBConnect struct {
	*elasticsearch.Client
	Options *elasticsearch.Config
}

// New return a new elasticsearch client, and try ping.
func New(opts *elasticsearch.Config) (*DBConnect, error) {
	client, err := elasticsearch.NewClient(*opts)
	if err != nil {
		return nil, err
	}

	resp, err := client.Ping(
		client.Ping.WithContext(context.Background()),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, err
	}

	return &DBConnect{
		Client:  client,
		Options: opts,
	}, nil
}
