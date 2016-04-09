package charon

import (
	"errors"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	contextKeySubject = "context_key_charon_subject"
)

// SecurityContext ....
type SecurityContext interface {
	context.Context
	oauth2.TokenSource
	// Subject ...
	Subject() (Subject, bool)
	// AccessToken ...
	AccessToken() (mnemosyne.AccessToken, bool)
}

type securityContext struct {
	context.Context
}

// NewSecurityContext allocates new context.
func NewSecurityContext(ctx context.Context) SecurityContext {
	return &securityContext{Context: ctx}
}

// Subject implements SecurityContext interface.
func (sc *securityContext) Subject() (Subject, bool) {
	return SubjectFromContext(sc)
}

// AccessToken implements SecurityContext interface.
func (sc *securityContext) AccessToken() (mnemosyne.AccessToken, bool) {
	return mnemosyne.AccessTokenFromContext(sc.Context)
}

// Token implements oauth2.TokenSource interface.
func (sc *securityContext) Token() (*oauth2.Token, error) {
	at, ok := sc.AccessToken()
	if !ok {
		return nil, errors.New("charon: missing access token, oauth2 token cannot be returned")
	}
	return &oauth2.Token{
		AccessToken: at.Encode(),
	}, nil
}
