package mysqlbuilder

import (
	"fmt"
	"net/url"
)

type Option func(opts *Options)

type Options struct {
	User     string
	Pass     string
	Host     string
	Port     int
	PoolSize int

	DBName  string
	Charset string
	Loc     string

	DBFilePath string
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

func WithPoolSize(size int) Option {
	return func(opts *Options) {
		opts.PoolSize = size
	}
}

func WithDBName(dbName string) Option {
	return func(opts *Options) {
		opts.DBName = dbName
	}
}

func WithCharset(charset string) Option {
	return func(opts *Options) {
		opts.Charset = charset
	}
}

func WithLoc(loc string) Option {
	return func(opts *Options) {
		opts.Loc = loc
	}
}

func WithDBFilePath(DBFilePath string) Option {
	return func(opts *Options) {
		opts.DBFilePath = DBFilePath
	}
}

func initOptions(opts ...Option) *Options {
	options := &Options{
		User:     "root",
		Host:     "127.0.0.1",
		Port:     3306,
		PoolSize: 10,

		Charset: "utf8",
		Loc:     "Asia/Shanghai",

		DBFilePath: "",
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// BuildDBDriver return a driver like "user:pass@tcp(host:port)/dbName?charset=utf8&loc=Asia%%2FShanghai"
func BuildDBDriver(opts *Options) string {

	driver := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&loc=%s",
		opts.User,
		opts.Pass,
		opts.Host,
		opts.Port,
		opts.DBName,
		opts.Charset,
		url.QueryEscape(opts.Loc), // 使用 url 编码转移 /
	)
	return driver
}
