package charonc

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type charonOptions struct {
	metadata metadata.MD
}

// Options configures how we set up the Client.
type Options func(*charonOptions)

// WithMetadata sets metadata that will be attachable to every request.
func WithMetadata(kv ...string) Options {
	return func(co *charonOptions) {
		co.metadata = metadata.Pairs(kv...)
	}
}

// Client is simplified version of rpc AuthClient.
// It contains most commonly used methods.
// For more powerful low level API check RPCClient interface.
type Client struct {
	options charonOptions
	auth    charonrpc.AuthClient
}

// New allocates new Charon instance with given options.
func New(conn *grpc.ClientConn, options ...Options) *Client {
	ch := &Client{
		auth: charonrpc.NewAuthClient(conn),
	}

	for _, o := range options {
		o(&ch.options)
	}

	return ch
}

// IsGranted implements Charon interface.
func (c *Client) IsGranted(ctx context.Context, userID int64, perm charon.Permission) (bool, error) {
	req := &charonrpc.IsGrantedRequest{
		UserId:     userID,
		Permission: perm.String(),
	}

	granted, err := c.auth.IsGranted(ctx, req)
	if err != nil {
		return false, err
	}

	return granted.Value, nil
}

// Actor implements Charon interface.
func (c *Client) Actor(ctx context.Context, token string) (*Actor, error) {
	resp, err := c.auth.Actor(ctx, &wrappers.StringValue{Value: token})
	if err != nil {
		return nil, err
	}

	return c.mapActor(resp), nil
}

// FromContext implements Charon interface.
func (c *Client) FromContext(ctx context.Context) (*Actor, error) {
	resp, err := c.auth.Actor(ctx, &wrappers.StringValue{})
	if err != nil {
		return nil, err
	}

	return c.mapActor(resp), nil
}

func (c *Client) mapActor(resp *charonrpc.ActorResponse) *Actor {
	return &Actor{
		ID:          resp.Id,
		Username:    resp.Username,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		IsSuperuser: resp.IsSuperuser,
		IsStaff:     resp.IsStuff,
		IsConfirmed: resp.IsConfirmed,
		IsActive:    resp.IsActive,
		Permissions: charon.NewPermissions(resp.Permissions...),
	}
}

// IsAuthenticated implements Charon interface.
func (c *Client) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	ok, err := c.auth.IsAuthenticated(ctx, &charonrpc.IsAuthenticatedRequest{
		AccessToken: token,
	})
	if err != nil {
		return false, err
	}

	return ok.Value, nil
}

// Login is a simple wrapper around rpc Login method.
func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	token, err := c.auth.Login(ctx, &charonrpc.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return token.Value, nil
}

// Logout implements Charon interface.
func (c *Client) Logout(ctx context.Context, token string) error {
	_, err := c.auth.Logout(ctx, &charonrpc.LogoutRequest{
		AccessToken: token,
	})
	return err
}

// Actor is a generic object that represent anything that can be under control of charon.
type Actor struct {
	ID          int64              `json:"id"`
	Username    string             `json:"username"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	IsSuperuser bool               `json:"isSuperuser"`
	IsActive    bool               `json:"isActive"`
	IsStaff     bool               `json:"isStaff"`
	IsConfirmed bool               `json:"isConfirmed"`
	Permissions charon.Permissions `json:"permissions"`
}

// NewActorContext returns a new Context that carries Actor value.
func NewActorContext(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, contextKeyActor, a)
}

// ActorFromContext returns the Actor value stored in context, if any.
func ActorFromContext(ctx context.Context) (Actor, bool) {
	s, ok := ctx.Value(contextKeyActor).(Actor)
	return s, ok
}
