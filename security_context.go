package charon

import (
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
)

const (
	contextKeySubject = "context_key_charon_subject"
)

// SecurityContext ....
type SecurityContext interface {
	context.Context
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
	return mnemosyne.AccessTokenFromContext(sc)
}
