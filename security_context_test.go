package charon_test

import (
	"reflect"
	"testing"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
)

func TestNewSecurityContext(t *testing.T) {
	sctx := charon.NewSecurityContext(context.Background())

	if _, ok := sctx.(charon.SecurityContext); !ok {
		t.Errorf("result should imeplement SecurityContext interface")
	}
}

func TestSecurityContext_Subject(t *testing.T) {
	expectedSubject := charon.Subject{ID: 1}
	ctx := charon.NewSubjectContext(context.Background(), expectedSubject)
	sctx := charon.NewSecurityContext(ctx)

	subject, ok := sctx.Subject()
	if ok {
		if !reflect.DeepEqual(subject, expectedSubject) {
			t.Error("provided and retrieved subject should be the same")
		}
	} else {
		t.Errorf("subject should be able retrieved")
	}
}

func TestSecurityContext_Subject_empty(t *testing.T) {
	sctx := charon.NewSecurityContext(context.Background())

	_, ok := sctx.Subject()
	if ok {
		t.Errorf("subject should not be there")
	}
}

func TestSecurityContext_Token(t *testing.T) {
	expectedToken := mnemosyne.NewToken([]byte("0000000001"), []byte("1"))
	ctx := mnemosyne.NewTokenContext(context.Background(), expectedToken)
	sctx := charon.NewSecurityContext(ctx)

	token, ok := sctx.Token()
	if ok {
		if !reflect.DeepEqual(token, expectedToken) {
			t.Error("provided and retrieved token should be the same")
		}
	} else {
		t.Errorf("token should be able retrieved")
	}
}


func TestSecurityContext_Token_empty(t *testing.T) {
	sctx := charon.NewSecurityContext(context.Background())

	_, ok := sctx.Token()
	if ok {
		t.Errorf("token should not be there")
	}
}