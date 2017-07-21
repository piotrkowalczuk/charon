package ldap

import "github.com/go-ldap/ldap"

type Client struct {
	SearchDN string
	conn     *ldap.Conn
}
