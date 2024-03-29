package clickhousebuilder

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type DBConnect struct {
	clickhouse.Conn
	Options *clickhouse.Options
}

// New return a new clickhouse client, and try ping.
//
//		conn, err := clickhousebuilder.New(&clickhouse.Options{
//		Addr: []string{fmt.Sprintf("%s:%d", "localhost", 9000)},
//		Auth: clickhouse.Auth{
//			Database: "database",
//			Username: "default",
//			Password: "default",
//		},
//		DialTimeout:     time.Second,
//		MaxOpenConns:    10,
//		MaxIdleConns:    10 / 2,
//		ConnMaxLifetime: time.Hour,
//	})
func New(opts *clickhouse.Options) (*DBConnect, error) {
	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &DBConnect{
		Conn:    conn,
		Options: opts,
	}, nil
}

// AllowExperimentalObjectType enable or disable object(json) type
func (db *DBConnect) AllowExperimentalObjectType(enable bool) error {
	return db.Conn.Exec(context.Background(), "SET allow_experimental_object_type = ?", enable)
}
