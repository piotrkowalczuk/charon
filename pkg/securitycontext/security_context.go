package securitycontext

import (
	"context"
	"errors"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/oauth2"
)

// Context ....
type Context interface {
	context.Context
	oauth2.TokenSource
	// Actor ...
	Actor() (Actor, bool)
	// AccessToken ...
	AccessToken() (string, bool)
}

type securityContext struct {
	context.Context
}

// NewSecurityContext allocates new context.
func NewSecurityContext(ctx context.Context) Context {
	return &securityContext{Context: ctx}
}

// Actor implements Context interface.
func (sc *securityContext) Actor() (Actor, bool) {
	return ActorFromContext(sc)
}

// AccessToken implements Context interface.
func (sc *securityContext) AccessToken() (string, bool) {
	return mnemosyne.AccessTokenFromContext(sc.Context)
}

// Token implements oauth2.TokenSource interface.
func (sc *securityContext) Token() (*oauth2.Token, error) {
	at, ok := sc.AccessToken()
	if !ok {
		return nil, errors.New("securitycontext: missing access token, oauth2 token cannot be returned")
	}
	return &oauth2.Token{
		AccessToken: at,
	}, nil
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

type key struct{}

var contextKeyActor = key{}

// NewActorContext returns a new Context that carries Actor value.
func NewActorContext(ctx context.Context, act Actor) context.Context {
	return context.WithValue(ctx, contextKeyActor, act)
}

// ActorFromContext returns the Actor value stored in context, if any.
func ActorFromContext(ctx context.Context) (Actor, bool) {
	act, ok := ctx.Value(contextKeyActor).(Actor)
	return act, ok
}
