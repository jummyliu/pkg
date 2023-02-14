package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoAdd 执行默认 ldap 查询
// 	AddOptions:
// 	WithBaseDN[*RequestOptions]   默认为创建连接时的 DN，如果创建连接时未指定，则这里必须指定
// 	WithControls[*RequestOptions]
// 	db.WithAttribute
// 	db.WithAddAttributes
func (db *DBConnect) DoAdd(ctx context.Context, opts ...AddOption) error {
	options := db.initAddOptions(opts...)
	err := db.Add(options.AddRequest)
	if err != nil {
		return err
	}
	return nil
}

type AddOption func(opts *AddOptions)
type AddOptions struct {
	*BaseOptions
	*ldap.AddRequest
}

func (db *DBConnect) WithAttribute(attrType string, attrVals []string) AddOption {
	return func(opts *AddOptions) {
		opts.Attribute(attrType, attrVals)
	}
}
func (db *DBConnect) WithAddAttributes(attrs []ldap.Attribute) AddOption {
	return func(opts *AddOptions) {
		opts.Attributes = attrs
	}
}
func (db *DBConnect) initAddOptions(opts ...AddOption) *AddOptions {
	options := &AddOptions{
		BaseOptions: &BaseOptions{
			BaseDN:   db.Options.BaseDN,
			Controls: db.Options.Controls,
		},
		AddRequest: &ldap.AddRequest{},
	}
	for _, opt := range opts {
		opt(options)
	}
	options.AddRequest.DN = options.BaseOptions.BaseDN
	options.AddRequest.Controls = options.BaseOptions.Controls
	return options
}
