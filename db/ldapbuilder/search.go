package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoSearch 执行默认 ldap 查询
//
//	ReuqestOptions:
//	WithBaseDN[*RequestOptions]   默认为创建连接时的 DN，如果创建连接时未指定，则这里必须指定
//	WithControls[*RequestOptions]
//	db.WithSearchAttributes       默认查询全部字段
//	db.WithScope
//	db.WithDerefAliases
//	db.WithSizeLimit
//	db.WithTimeLimit
//	db.WithTypesOnly
//	db.WithPageSize
func (db *DBConnect) DoSearch(ctx context.Context, filter string, opts ...RequestOption) (results []*ldap.Entry, err error) {
	options := db.initRequestOptions(opts...)
	options.Filter = filter
	result, err := db.SearchWithPaging(options.SearchRequest, options.PageSize)
	if err != nil {
		return nil, err
	}
	return result.Entries, nil
}

type RequestOption func(opts *RequestOptions)
type RequestOptions struct {
	PageSize uint32
	*BaseOptions
	*ldap.SearchRequest
}

func (db *DBConnect) WithScope(scope int) RequestOption {
	return func(opts *RequestOptions) {
		opts.Scope = scope
	}
}
func (db *DBConnect) WithDerefAliases(derefAliases int) RequestOption {
	return func(opts *RequestOptions) {
		opts.DerefAliases = derefAliases
	}
}
func (db *DBConnect) WithSizeLimit(sizeLimit int) RequestOption {
	return func(opts *RequestOptions) {
		opts.SizeLimit = sizeLimit
	}
}
func (db *DBConnect) WithTimeLimit(timeLimit int) RequestOption {
	return func(opts *RequestOptions) {
		opts.TimeLimit = timeLimit
	}
}
func (db *DBConnect) WithTypesOnly(typesOnly bool) RequestOption {
	return func(opts *RequestOptions) {
		opts.TypesOnly = typesOnly
	}
}
func (db *DBConnect) WithFilter(filter string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Filter = filter
	}
}
func (db *DBConnect) WithPageSize(pageSize uint32) RequestOption {
	return func(opts *RequestOptions) {
		opts.PageSize = pageSize
	}
}
func (db *DBConnect) WithSearchAttributes(attributes []string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Attributes = attributes
	}
}
func (db *DBConnect) initRequestOptions(opts ...RequestOption) *RequestOptions {
	options := &RequestOptions{
		PageSize: 200,
		BaseOptions: &BaseOptions{
			BaseDN:   db.Options.BaseDN,
			Controls: db.Options.Controls,
		},
		SearchRequest: &ldap.SearchRequest{
			Scope:        ldap.ScopeWholeSubtree,
			DerefAliases: ldap.NeverDerefAliases,
			SizeLimit:    0,
			TimeLimit:    0,
			TypesOnly:    false,
			Filter:       "",
			Attributes:   nil,
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	options.SearchRequest.BaseDN = options.BaseOptions.BaseDN
	options.SearchRequest.Controls = options.BaseOptions.Controls
	return options
}
