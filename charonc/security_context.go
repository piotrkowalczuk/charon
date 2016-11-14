package charonc

import (
	"errors"

	"github.com/piotrkowalczuk/mnemosyne"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	contextKeyActor = "context_key_charon_actor"
)

// SecurityContext ....
type SecurityContext interface {
	context.Context
	oauth2.TokenSource
	// Subject ...
	Subject() (Actor, bool)
	// AccessToken ...
	AccessToken() (string, bool)
}

type securityContext struct {
	context.Context
}

// NewSecurityContext allocates new context.
func NewSecurityContext(ctx context.Context) SecurityContext {
	return &securityContext{Context: ctx}
}

// Subject implements SecurityContext interface.
func (sc *securityContext) Subject() (Actor, bool) {
	return ActorFromContext(sc)
}

// AccessToken implements SecurityContext interface.
func (sc *securityContext) AccessToken() (string, bool) {
	return mnemosyne.AccessTokenFromContext(sc.Context)
}

// Token implements oauth2.TokenSource interface.
func (sc *securityContext) Token() (*oauth2.Token, error) {
	at, ok := sc.AccessToken()
	if !ok {
		return nil, errors.New("charonc: missing access token, oauth2 token cannot be returned")
	}
	return &oauth2.Token{
		AccessToken: at,
	}, nil
}
