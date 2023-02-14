package mongobuilder

import "fmt"

type Option func(opts *Options)

type Options struct {
	User     string
	Pass     string
	Host     string
	Port     int
	PoolSize uint64

	DBName string
}

func WithUser(user string) Option {
	return func(opts *Options) {
		opts.User = user
	}
}

func WithPass(pass string) Option {
	return func(opts *Options) {
		opts.Pass = pass
	}
}

func WithHost(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

func WithPort(port int) Option {
	return func(opts *Options) {
		opts.Port = port
	}
}

func WithPoolSize(size uint64) Option {
	return func(opts *Options) {
		opts.PoolSize = size
	}
}

func WithDBName(dbName string) Option {
	return func(opts *Options) {
		opts.DBName = dbName
	}
}

func initOptions(opts ...Option) *Options {
	options := &Options{
		User:     "mongo",
		Host:     "127.0.0.1",
		Port:     27017,
		PoolSize: 10,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// BuildDBDriver return a driver like "mongodb://user:pass@host:port/dbName"
func BuildDBDriver(opts *Options) string {
	driver := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		opts.User,
		opts.Pass,
		opts.Host,
		opts.Port,
		opts.DBName,
	)
	return driver
}
