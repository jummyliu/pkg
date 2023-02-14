package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoAdd 执行默认 ldap 查询
// 	AddOptions:
// 	WithBaseDN[*RequestOptions]   默认为创建连接时的 DN，如果创建连接时未指定，则这里必须指定
// 	WithControls[*RequestOptions]
func (db *DBConnect) DoDel(ctx context.Context, opts ...DelOption) error {
	options := db.initDelOptions(opts...)
	err := db.Del(options.DelRequest)
	if err != nil {
		return err
	}
	return nil
}

type DelOption func(*DelOptions)
type DelOptions struct {
	*BaseOptions
	*ldap.DelRequest
}

func (db *DBConnect) initDelOptions(opts ...DelOption) *DelOptions {
	options := &DelOptions{
		BaseOptions: &BaseOptions{
			BaseDN:   db.Options.BaseDN,
			Controls: db.Options.Controls,
		},
		DelRequest: &ldap.DelRequest{},
	}
	for _, opt := range opts {
		opt(options)
	}
	options.DelRequest.DN = options.BaseOptions.BaseDN
	options.DelRequest.Controls = options.BaseOptions.Controls
	return options
}
