package charon

import (
	"github.com/piotrkowalczuk/mnemosyne"
	"code.google.com/p/go.net/context"
)

const (
	contextKeySubject = "context_key_charon_subject";
)

// SecurityContext ....
type SecurityContext interface{
	context.Context
	// Subject ...
	Subject() (Subject, bool)
	// Token ...
	Token() (mnemosyne.Token, bool)
}

type securityContext struct {
	context.Context
}

// NewSecurityContext allocates new context.
func NewSecurityContext(ctx context.Context) SecurityContext {
	return &securityContext{Context: ctx}
}

// Subject implements SecurityContext interface.
func (sc *securityContext) Subject() (Subject, bool){
	return SubjectFromContext(sc)
}

// Token implements SecurityContext interface.
func (sc *securityContext) Token() (mnemosyne.Token, bool) {
	return mnemosyne.TokenFromContext(sc)
}
