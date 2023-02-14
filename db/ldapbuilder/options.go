package ldapbuilder

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type Option func(opts *Options)

type Options struct {
	Protocol string // ldap
	User     string
	Pass     string
	Host     string
	Port     int  // 389
	TLS      bool // false
	*BaseOptions
}

// WithProtocol eg: ldap
func WithProtocol(protocol string) Option {
	return func(opts *Options) {
		opts.Protocol = protocol
	}
}

// WithUser eg: administrator@test.com
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

// WithHost eg: 127.0.0.1
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

func WithTLS(tls bool) Option {
	return func(opts *Options) {
		opts.TLS = tls
	}
}

func initOptions(opts ...Option) *Options {
	options := &Options{
		Protocol:    "ldap",
		Port:        389,
		TLS:         false,
		BaseOptions: &BaseOptions{},
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// BuildDBDriver return a driver like "protocol://host:port"
func BuildDBDriver(opts *Options) string {
	driver := fmt.Sprintf(
		"%s://%s:%d",
		opts.Protocol,
		opts.Host,
		opts.Port,
	)
	return driver
}

type IBaseDN interface {
	setBaseDN(baseDN string)
}
type IControls interface {
	setControls(controls []ldap.Control)
}

type BaseOption func(opt *BaseOptions)
type BaseOptions struct {
	BaseDN   string
	Controls []ldap.Control
}

func (opt *BaseOptions) setBaseDN(baseDN string) {
	opt.BaseDN = baseDN
}
func (opt *BaseOptions) setControls(controls []ldap.Control) {
	opt.Controls = controls
}

func WithBaseDN[T IBaseDN](baseDN string) func(T) {
	return func(t T) {
		t.setBaseDN(baseDN)
	}
}
func WithControls[T IControls](controls []ldap.Control) func(T) {
	return func(t T) {
		t.setControls(controls)
	}
}
