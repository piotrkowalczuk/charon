package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type charonOptions struct {
	metadata metadata.MD
}

// CharonOption configures how we set up the client.
type CharonOption func(*charonOptions)

// WithMetadata sets metadata that will be attachacked to every request.
func WithMetadata(kv ...string) CharonOption {
	return func(co *charonOptions) {
		co.metadata = metadata.Pairs(kv...)
	}
}

// Charon is an interface that describes simplified client.
// It contains most commonly used methods.
// For more powerful low level API check RPCClient interface.
type Charon interface {
	IsGranted(context.Context, int64, Permission) (bool, error)
	IsAuthenticated(context.Context, string) (bool, error)
	Subject(context.Context, string) (*Subject, error)
	FromContext(context.Context) (*Subject, error)
	Login(context.Context, string, string) (string, error)
	Logout(context.Context, string) error
}

type charon struct {
	options charonOptions
	client  RPCClient
}

// New allocates new Charon instance with given options.
func New(conn *grpc.ClientConn, options ...CharonOption) Charon {
	ch := &charon{
		client: NewRPCClient(conn),
	}

	for _, o := range options {
		o(&ch.options)
	}

	return ch
}

// IsGranted implements Charon interface.
func (c *charon) IsGranted(ctx context.Context, userID int64, perm Permission) (bool, error) {
	req := &IsGrantedRequest{
		UserId:     userID,
		Permission: perm.String(),
	}

	res, err := c.client.IsGranted(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Granted, nil
}

// Subject implements Charon interface.
func (c *charon) Subject(ctx context.Context, token string) (*Subject, error) {
	resp, err := c.client.Subject(ctx, &SubjectRequest{AccessToken: token})
	if err != nil {
		return nil, err
	}

	return c.mapSubject(resp), nil
}

// FromContext implements Charon interface.
func (c *charon) FromContext(ctx context.Context) (*Subject, error) {
	resp, err := c.client.Subject(ctx, &SubjectRequest{})
	if err != nil {
		return nil, err
	}

	return c.mapSubject(resp), nil
}

func (c *charon) mapSubject(resp *SubjectResponse) *Subject {
	return &Subject{
		ID:          resp.Id,
		Username:    resp.Username,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		IsSuperuser: resp.IsSuperuser,
		IsStaff:     resp.IsStuff,
		IsConfirmed: resp.IsConfirmed,
		IsActive:    resp.IsActive,
		Permissions: NewPermissions(resp.Permissions...),
	}
}

// IsAuthenticated implements Charon interface.
func (c *charon) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	res, err := c.client.IsAuthenticated(ctx, &IsAuthenticatedRequest{
		AccessToken: token,
	})
	if err != nil {
		return false, err
	}

	return res.Authenticated, nil
}

// Login implements Charon interface.
// TODO: reimplement
func (c *charon) Login(ctx context.Context, username, password string) (string, error) {
	res, err := c.client.Login(ctx, &LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return string(res.AccessToken), nil
}

// Logout implements Charon interface.
func (c *charon) Logout(ctx context.Context, token string) error {
	_, err := c.client.Logout(ctx, &LogoutRequest{
		AccessToken: token,
	})
	return err
}

// Subject is a generic object that represent anything that can be under control of charon.
type Subject struct {
	ID          int64       `json:"id"`
	Username    string      `json:"username"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	IsSuperuser bool        `json:"isSuperuser"`
	IsActive    bool        `json:"isActive"`
	IsStaff     bool        `json:"isStaff"`
	IsConfirmed bool        `json:"isConfirmed"`
	Permissions Permissions `json:"permissions"`
}

// Name return concatenated first and last name.
func (u *User) Name() string {
	return u.FirstName + " " + u.LastName
}

// NewSubjectContext returns a new Context that carries Subject value.
func NewSubjectContext(ctx context.Context, s Subject) context.Context {
	return context.WithValue(ctx, contextKeySubject, s)
}

// SubjectFromContext returns the Subject value stored in context, if any.
func SubjectFromContext(ctx context.Context) (Subject, bool) {
	s, ok := ctx.Value(contextKeySubject).(Subject)
	return s, ok
}
