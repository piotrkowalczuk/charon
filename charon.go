package charon

import (
	"errors"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	contextKeyUser  = "context_key_charon_user"
	DecisionUnknown = 0
	DecisionGranted = 1
)

var (
	// ErrMissingTokenInContext can be returned by functions
	// that are using arbitrary token taken from a context if it missing.
	ErrMissingTokenInContext = errors.New("charon: missing token in context")
)

// NewContext returns a new Context that carries User value.
func NewContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// FromContext returns the User value stored in context, if any.
func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(contextKeyUser).(User)
	return u, ok
}

// Error ...
type Error struct {
	StatusCode   int
	InternalCode int
	Message      string
	Validation   map[string][]string
}

// Error ...
func (e *Error) Error() string {
	return e.Message
}

// AddValidation ...
func (e *Error) AddValidation(key, value string) {
	if e.Validation[key] == nil {
		e.Validation[key] = make([]string, 0, 1)
	}

	e.Validation[key] = append(e.Validation[key], value)
}

// AuthorizationChecker ..
// TODO: unstable
type AuthorizationChecker func(context.Context, Permission, ...interface{}) (bool, error)

// Charon ...
type Charon interface {
	IsGranted(context.Context, mnemosyne.Token, Permission) (bool, error)
	IsAuthenticated(context.Context, mnemosyne.Token) (bool, error)
	Subject(context.Context, mnemosyne.Token) (*Subject, error)
	Login(context.Context, string, string) (*mnemosyne.Token, error)
	Logout(context.Context, mnemosyne.Token) error
}

type charon struct {
	meta   metadata.MD
	client RPCClient
}

// CharonOpts ...
type CharonOpts struct {
	Metadata []string
}

// New allocates new Charon instance.
func New(conn *grpc.ClientConn, options CharonOpts) Charon {
	return &charon{
		meta:   metadata.Pairs(options.Metadata...),
		client: NewRPCClient(conn),
	}
}

// IsGranted implements Charon interface.
func (c *charon) IsGranted(ctx context.Context, token mnemosyne.Token, perm Permission) (bool, error) {
	req := &IsGrantedRequest{
		Token:      &token,
		Permission: perm.String(),
	}

	res, err := c.client.IsGranted(ctx, req)
	if err != nil {
		return false, err
	}

	return res.IsGranted, nil
}

// Subject implements Charon interface.
func (c *charon) Subject(ctx context.Context, token mnemosyne.Token) (*Subject, error) {
	resp1, err := c.client.GetUser(ctx, &GetUserRequest{Token: &token})
	if err != nil {
		return nil, err
	}
	resp2, err := c.client.GetUserPermissions(ctx, &GetUserPermissionsRequest{Token: &token})
	if err != nil {
		return nil, err
	}

	return &Subject{
		ID:          resp1.User.Id,
		Name:        resp1.User.Name(),
		Email:       resp1.User.Username,
		Permissions: NewPermissions(resp2.Permissions),
	}, nil
}

// IsAuthenticated implements Charon interface.
func (c *charon) IsAuthenticated(ctx context.Context, token mnemosyne.Token) (bool, error) {
	res, err := c.client.IsAuthenticated(ctx, &IsAuthenticatedRequest{
		Token: &token,
	})
	if err != nil {
		return false, err
	}

	return res.IsAuthenticated, nil
}

// Login implements Charon interface.
func (c *charon) Login(ctx context.Context, username, password string) (*mnemosyne.Token, error) {
	res, err := c.client.Login(ctx, &LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return res.Token, nil
}

// Logout implements Charon interface.
func (c *charon) Logout(ctx context.Context, token mnemosyne.Token) error {
	_, err := c.client.Logout(ctx, &LogoutRequest{
		Token: &token,
	})
	return err
}

// Subject is a generic object that represent anything that can be under control of charon.
type Subject struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Permissions Permissions `json:"permissions"`
}

// Name return concatenated first and last name.
func (u *User) Name() string {
	return u.FirstName + " " + u.LastName
}
