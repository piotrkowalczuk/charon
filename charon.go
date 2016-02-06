package charon

import (
	"errors"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	DecisionUnknown = 0
	DecisionGranted = 1
)

var (
	// ErrMissingTokenInContext can be returned by functions
	// that are using arbitrary token taken from a context if it missing.
	ErrMissingTokenInContext = errors.New("charon: missing token in context")
)

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

// Charon ...
type Charon interface {
	IsGranted(context.Context, mnemosyne.Token, Permission) (bool, error)
	IsAuthenticated(context.Context, mnemosyne.Token) (bool, error)
	Subject(context.Context, mnemosyne.Token) (*Subject, error)
	FromContext(context.Context) (*Subject, error)
	Login(context.Context, string, string) (*mnemosyne.Token, error)
	Logout(context.Context, mnemosyne.Token) error
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
func (c *charon) IsGranted(ctx context.Context, token mnemosyne.Token, perm Permission) (bool, error) {
	req := &IsGrantedRequest{
		Token:      &token,
		Permission: perm.String(),
	}

	res, err := c.client.IsGranted(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Granted, nil
}

// Subject implements Charon interface.
func (c *charon) Subject(ctx context.Context, token mnemosyne.Token) (*Subject, error) {
	resp, err := c.client.Subject(ctx, &SubjectRequest{Token: &token})
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
func (c *charon) IsAuthenticated(ctx context.Context, token mnemosyne.Token) (bool, error) {
	res, err := c.client.IsAuthenticated(ctx, &IsAuthenticatedRequest{
		Token: &token,
	})
	if err != nil {
		return false, err
	}

	return res.Authenticated, nil
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
