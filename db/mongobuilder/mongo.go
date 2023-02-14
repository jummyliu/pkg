package mongobuilder

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnect struct {
	*mongo.Database
	Options *Options
}

// New return a new mongo client.
func New(opts ...Option) (*DBConnect, error) {
	initOpts := initOptions(opts...)
	driver := BuildDBDriver(initOpts)

	clientOpts := options.Client().ApplyURI(driver).SetMaxPoolSize(initOpts.PoolSize)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)

	if err != nil {
		return nil, err
	}
	db := client.Database(initOpts.DBName)
	return &DBConnect{
		Database: db,
		Options:  initOpts,
	}, nil
}
