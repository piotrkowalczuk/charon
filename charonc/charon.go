package charonc

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"golang.org/x/net/context"
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
	// RPCClient holds gRPC client from charonrpc package.
	// It's not safe to change concurrently.
	RPCClient charonrpc.AuthClient
}

// New allocates new Client instance with given options and gRPC connection.
func New(conn *grpc.ClientConn, options ...Options) *Client {
	ch := &Client{
		RPCClient: charonrpc.NewAuthClient(conn),
	}

	for _, o := range options {
		o(&ch.options)
	}

	return ch
}

// IsGranted returns true if user has granted given permission.
func (c *Client) IsGranted(ctx context.Context, userID int64, perm charon.Permission) (bool, error) {
	req := &charonrpc.IsGrantedRequest{
		UserId:     userID,
		Permission: perm.String(),
	}

	granted, err := c.RPCClient.IsGranted(ctx, req)
	if err != nil {
		return false, err
	}

	return granted.Value, nil
}

// Actor returns Actor for given token if logged in.
func (c *Client) Actor(ctx context.Context, token string) (*Actor, error) {
	resp, err := c.RPCClient.Actor(ctx, &wrappers.StringValue{Value: token})
	if err != nil {
		return nil, err
	}

	return c.mapActor(resp), nil
}

// FromContext works like Actor but retrieves access token from the context.
func (c *Client) FromContext(ctx context.Context) (*Actor, error) {
	resp, err := c.RPCClient.Actor(ctx, &wrappers.StringValue{})
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

// IsAuthenticated returns true if given access token exists.
func (c *Client) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	ok, err := c.RPCClient.IsAuthenticated(ctx, &charonrpc.IsAuthenticatedRequest{
		AccessToken: token,
	})
	if err != nil {
		return false, err
	}

	return ok.Value, nil
}

// Login is a simple wrapper around rpc Login method.
func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	token, err := c.RPCClient.Login(ctx, &charonrpc.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return token.Value, nil
}

// Logout removes given token making actor logged out.
func (c *Client) Logout(ctx context.Context, token string) error {
	_, err := c.RPCClient.Logout(ctx, &charonrpc.LogoutRequest{
		AccessToken: token,
	})
	return err
}