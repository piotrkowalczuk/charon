package ldap

import "github.com/go-ldap/ldap"

type Client struct {
	SearchDN string
	conn     *ldap.Conn
}

//func (c *Client) Search(username string) {
//	var filter string
//	if strings.Contains(username, "@") {
//		filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(mail=%s))", ldap.EscapeFilter(username))
//	} else {
//		parts := strings.Split(username, "@")
//		if len(parts) != 2 {
//			return nil, grpc.Errorf(codes.InvalidArgument, "invalid email address")
//		}
//		filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(parts[0]))
//	}
//	res, err := c.conn.Search(ldap.NewSearchRequest(
//		c.SearchDN,
//		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//		filter,
//		[]string{"dn", "givenName", "sn", "mail", "cn", "ou", "dc"},
//		nil,
//	))
//	if err != nil {
//		return nil, fmt.Errorf("ldap search failure: %s", err.Error())
//	}
//}
