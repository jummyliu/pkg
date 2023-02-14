package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoModify 执行默认 ldap 查询
// 	ModifyOptions:
// 	WithBaseDN[*RequestOptions]   默认为创建连接时的 DN，如果创建连接时未指定，则这里必须指定
// 	WithControls[*RequestOptions]
// 	db.WithAdd
// 	db.WithDelete
// 	db.WithReplace
// 	db.WithIncrement
// 	db.WithChanges
func (db *DBConnect) DoModify(ctx context.Context, opts ...ModifyOption) (*ldap.ModifyResult, error) {
	options := db.initModifyOptions(opts...)
	result, err := db.ModifyWithResult(options.ModifyRequest)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type ModifyOption func(opts *ModifyOptions)
type ModifyOptions struct {
	*BaseOptions
	*ldap.ModifyRequest
}

func (db *DBConnect) WithAdd(attrType string, attrVals []string) ModifyOption {
	return func(opts *ModifyOptions) {
		opts.Add(attrType, attrVals)
	}
}
func (db *DBConnect) WithDelete(attrType string, attrVals []string) ModifyOption {
	return func(opts *ModifyOptions) {
		opts.Delete(attrType, attrVals)
	}
}
func (db *DBConnect) WithReplace(attrType string, attrVals []string) ModifyOption {
	return func(opts *ModifyOptions) {
		opts.Replace(attrType, attrVals)
	}
}
func (db *DBConnect) WithIncrement(attrType string, attrVals string) ModifyOption {
	return func(opts *ModifyOptions) {
		opts.Increment(attrType, attrVals)
	}
}
func (db *DBConnect) WithChanges(changes []ldap.Change) ModifyOption {
	return func(opts *ModifyOptions) {
		opts.Changes = changes
	}
}
func (db *DBConnect) initModifyOptions(opts ...ModifyOption) *ModifyOptions {
	options := &ModifyOptions{
		BaseOptions: &BaseOptions{
			BaseDN:   db.Options.BaseDN,
			Controls: db.Options.Controls,
		},
		ModifyRequest: &ldap.ModifyRequest{},
	}
	for _, opt := range opts {
		opt(options)
	}
	options.ModifyRequest.DN = options.BaseOptions.BaseDN
	options.ModifyRequest.Controls = options.BaseOptions.Controls
	return options
}
