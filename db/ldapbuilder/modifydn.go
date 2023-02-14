package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoModifyDN 执行默认 ldap 查询
// 	ModifyDNOptions:
// 	WithBaseDN[*RequestOptions]   默认为创建连接时的 DN，如果创建连接时未指定，则这里必须指定
// 	WithControls[*RequestOptions]
// 	db.WithNewRDN
// 	db.WithDeleteOldRDN
// 	db.WithNewSuperior
func (db *DBConnect) DoModifyDN(ctx context.Context, opts ...ModifyDNOption) error {
	options := db.initModifyDNOptions(opts...)
	err := db.ModifyDN(options.ModifyDNRequest)
	if err != nil {
		return err
	}
	return nil
}

type ModifyDNOption func(opts *ModifyDNOptions)
type ModifyDNOptions struct {
	*BaseOptions
	*ldap.ModifyDNRequest
}

func (db *DBConnect) WithNewRDN(rdn string) ModifyDNOption {
	return func(opts *ModifyDNOptions) {
		opts.NewRDN = rdn
	}
}
func (db *DBConnect) WithDeleteOldRDN(delOld bool) ModifyDNOption {
	return func(opts *ModifyDNOptions) {
		opts.DeleteOldRDN = delOld
	}
}
func (db *DBConnect) WithNewSuperior(newSup string) ModifyDNOption {
	return func(opts *ModifyDNOptions) {
		opts.NewSuperior = newSup
	}
}
func (db *DBConnect) initModifyDNOptions(opts ...ModifyDNOption) *ModifyDNOptions {
	options := &ModifyDNOptions{
		BaseOptions: &BaseOptions{
			BaseDN:   db.Options.BaseDN,
			Controls: db.Options.Controls,
		},
		ModifyDNRequest: &ldap.ModifyDNRequest{},
	}
	for _, opt := range opts {
		opt(options)
	}
	options.ModifyDNRequest.DN = options.BaseOptions.BaseDN
	options.ModifyDNRequest.Controls = options.BaseOptions.Controls
	return options
}
