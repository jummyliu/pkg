package ldapbuilder

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

type DBConnect struct {
	*ldap.Conn
	Options *Options
}

func New(opts ...Option) (*DBConnect, error) {
	options := initOptions(opts...)
	driver := BuildDBDriver(options)
	conn, err := ldap.DialURL(driver)
	if err != nil {
		return nil, err
	}
	if options.TLS {
		if err = conn.StartTLS(&tls.Config{
			InsecureSkipVerify: true,
		}); err != nil {
			return nil, err
		}
	}
	if err = conn.Bind(options.User, options.Pass); err != nil {
		conn.Close()
		return nil, err
	}
	return &DBConnect{
		Conn:    conn,
		Options: options,
	}, nil
}

// func (db *DBConnect) DoDel() {
// 	db.Del()
// }

// func (db *DBConnect) DoModify() {
// 	ldap.NewModifyDNRequest()
// 	db.ModifyWithResult()
// }

// func (db *DBConnect) DoPasswordModify() {
// 	db.PasswordModify()
// }

// func (db *DBConnect) DoModifyDN() {
// 	db.ModifyDN()
// }
