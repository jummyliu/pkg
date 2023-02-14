package ldapbuilder

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

// DoPasswordModify 执行默认 ldap 查询
// 	PasswordModifyOptions:
// 	db.WithUserIdentity 必须
// 	db.WithOldPassword  必须
// 	db.WithNewPassword  必须
func (db *DBConnect) DoPasswordModify(ctx context.Context, opts ...PasswordModifyOption) (*ldap.PasswordModifyResult, error) {
	options := db.initPasswordModifyOptions(opts...)
	result, err := db.PasswordModify(options.PasswordModifyRequest)
	if err != nil {
		return nil, err
	}
	return result, err
}

type PasswordModifyOption func(opts *PasswordModifyOptions)
type PasswordModifyOptions struct {
	*ldap.PasswordModifyRequest
}

func (db *DBConnect) WithUserIdentity(userIdentity string) PasswordModifyOption {
	return func(opts *PasswordModifyOptions) {
		opts.UserIdentity = userIdentity
	}
}
func (db *DBConnect) WithOldPassword(oldPassword string) PasswordModifyOption {
	return func(opts *PasswordModifyOptions) {
		opts.OldPassword = oldPassword
	}
}
func (db *DBConnect) WithNewPassword(newPassword string) PasswordModifyOption {
	return func(opts *PasswordModifyOptions) {
		opts.NewPassword = newPassword
	}
}
func (db *DBConnect) initPasswordModifyOptions(opts ...PasswordModifyOption) *PasswordModifyOptions {
	options := &PasswordModifyOptions{
		PasswordModifyRequest: &ldap.PasswordModifyRequest{},
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
