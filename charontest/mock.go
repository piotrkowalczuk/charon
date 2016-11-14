package charontest

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonc"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Client struct {
	mock.Mock
}

func (c *Client) IsGranted(ctx context.Context, userID int64, perm charon.Permission) (bool, error) {
	a := c.Called(ctx, userID, perm)

	return a.Bool(0), a.Error(1)
}

// IsAuthenticated implements Client interface.
func (c *Client) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	a := c.Called(ctx, token)

	return a.Bool(0), a.Error(1)
}

// Actor implements Client interface.
func (c *Client) Actor(ctx context.Context, token string) (*charonc.Actor, error) {
	a := c.Called(ctx, token)

	act, err := a.Get(0), a.Error(1)
	if err != nil {
		return nil, err
	}

	return act.(*charonc.Actor), nil
}

// Context implements Client interface.
func (c *Client) FromContext(ctx context.Context) (*charonc.Actor, error) {
	a := c.Called(ctx)

	subj, err := a.Get(0), a.Error(1)
	if err != nil {
		return nil, err
	}

	return subj.(*charonc.Actor), nil
}

// Login implements Client interface.
func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	a := c.Called(ctx, username, password)

	ses, err := a.Get(0), a.Error(1)
	if err != nil {
		return "", err
	}

	return ses.(string), nil
}

// Logout implements Client interface.
func (c *Client) Logout(ctx context.Context, token string) error {
	a := c.Called(ctx, token)

	return a.Error(0)
}
