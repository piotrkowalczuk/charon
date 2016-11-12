package charonc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func ExampleSecurityContext() {
	token := "0000000001some hash"
	subject := Actor{
		ID:       1,
		Username: "j.kowalski@gmail.com",
	}
	ctx := NewActorContext(context.Background(), subject)
	ctx = mnemosyne.NewAccessTokenContext(ctx, token)
	sctx := NewSecurityContext(ctx)

	var (
		t   *oauth2.Token
		s   Actor
		err error
		ok  bool
	)
	if t, err = sctx.Token(); err != nil {
		fmt.Printf("unexpected error: %s", err.Error())
	} else {
		fmt.Println(t.AccessToken)
	}
	if s, ok = sctx.Subject(); ok {
		fmt.Println(s.ID)
		fmt.Println(s.Username)
	}

	// Output:
	// 0000000001some hash
	// 1
	// j.kowalski@gmail.com
}

func TestNewSecurityContext(t *testing.T) {
	sctx := NewSecurityContext(context.Background())

	if _, ok := sctx.(SecurityContext); !ok {
		t.Error("result should imeplement SecurityContext interface")
	}
}

func TestSecurityContext_Subject(t *testing.T) {
	expectedSubject := Actor{ID: 1}
	ctx := NewActorContext(context.Background(), expectedSubject)
	sctx := NewSecurityContext(ctx)

	subject, ok := sctx.Subject()
	if ok {
		if !reflect.DeepEqual(subject, expectedSubject) {
			t.Error("provided and retrieved subject should be the same")
		}
	} else {
		t.Error("subject should be able retrieved")
	}
}

func TestSecurityContext_Subject_empty(t *testing.T) {
	sctx := NewSecurityContext(context.Background())

	_, ok := sctx.Subject()
	if ok {
		t.Error("subject should not be there")
	}
}

func TestSecurityContext_Token(t *testing.T) {
	expectedToken := "00000000011"
	ctx := mnemosyne.NewAccessTokenContext(context.Background(), expectedToken)
	sctx := NewSecurityContext(ctx)

	token, err := sctx.Token()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if token.AccessToken != string(expectedToken) {
		t.Error("provided and retrieved token should be the same")
	}
}

func TestSecurityContext_Token_empty(t *testing.T) {
	sctx := NewSecurityContext(context.Background())

	_, err := sctx.Token()
	if err == nil {
		t.Error("expected error, got nil")
	}
}
