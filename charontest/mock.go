package charontest

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Charon struct {
	mock.Mock
}

func (c *Charon) IsGranted(ctx context.Context, userID int64, perm charon.Permission) (bool, error) {
	a := c.Called(ctx, userID, perm)

	return a.Bool(0), a.Error(1)
}

// IsAuthenticated implements Charon interface.
func (c *Charon) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	a := c.Called(ctx, token)

	return a.Bool(0), a.Error(1)
}

// Subject implements Charon interface.
func (c *Charon) Subject(ctx context.Context, token string) (*charon.Subject, error) {
	a := c.Called(ctx, token)

	subj, err := a.Get(0), a.Error(1)
	if err != nil {
		return nil, err
	}

	return subj.(*charon.Subject), nil
}

// Context implements Charon interface.
func (c *Charon) FromContext(ctx context.Context) (*charon.Subject, error) {
	a := c.Called(ctx)

	subj, err := a.Get(0), a.Error(1)
	if err != nil {
		return nil, err
	}

	return subj.(*charon.Subject), nil
}

// Login implements Charon interface.
func (c *Charon) Login(ctx context.Context, username, password string) (*mnemosynerpc.AccessToken, error) {
	a := c.Called(ctx, username, password)

	ses, err := a.Get(0), a.Error(1)
	if err != nil {
		return nil, err
	}

	return ses.(*mnemosynerpc.AccessToken), nil
}

// Logout implements Charon interface.
func (c *Charon) Logout(ctx context.Context, token string) error {
	a := c.Called(ctx, token)

	return a.Error(0)
}
