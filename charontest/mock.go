package charontest

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Charon struct {
	mock.Mock
}

// IsGranted implements Charon interface.
func (c *Charon) IsGranted(ctx context.Context, perm charon.Permission, args ...interface{}) (bool, error) {
	a := c.Called(append([]interface{}{ctx, perm}, args...)...)

	return a.Bool(0), a.Error(1)
}

// IsAuthenticated implements Charon interface.
func (c *Charon) IsAuthenticated(ctx context.Context) (bool, error) {
	a := c.Called(ctx)

	return a.Bool(0), a.Error(1)
}

// Login implements Charon interface.
func (c *Charon) Login(username, password string) (*mnemosyne.Session, error) {
	a := c.Called(username, password)

	return a.Get(0).(*mnemosyne.Session), a.Error(1)
}

// Logout implements Charon interface.
func (c *Charon) Logout(token *mnemosyne.Token) error {
	a := c.Called(token)

	return a.Error(0)
}
