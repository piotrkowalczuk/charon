package grpcerr_test

import (
	"errors"
	"io"
	"testing"

	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

func TestE_noArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = grpcerr.E()
}

func TestE(t *testing.T) {
	cases := map[string]struct {
		args []interface{}
		exp  string
	}{
		"message-only": {
			args: []interface{}{"a message"},
			exp:  "a message",
		},
		"error-only": {
			args: []interface{}{errors.New("an error")},
			exp:  "an error",
		},
		"error-and-code": {
			args: []interface{}{errors.New("an error"), codes.InvalidArgument},
			exp:  "InvalidArgument: an error",
		},
		"error-and-code-and-description": {
			args: []interface{}{errors.New("an error"), codes.InvalidArgument, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "example",
						Description: "something went wrong",
					},
				},
			}},
			exp: "InvalidArgument: an error",
		},
		"error-and-message": {
			args: []interface{}{"a message", errors.New("an error")},
			exp:  "a message: an error",
		},
		"error-and-message-and-kind": {
			args: []interface{}{"a message", errors.New("an error"), grpcerr.Kind("a kind")},
			exp:  "a kind: a message: an error",
		},
		"error-and-message-and-kind-and-operator": {
			args: []interface{}{
				"a message",
				errors.New("an error"),
				grpcerr.Kind("a kind"),
				grpcerr.Op("creativeservrpc.AccountManager/Get")},
			exp: "creativeservrpc.AccountManager/Get: a kind: a message: an error",
		},
		"error-and-message-and-fields": {
			args: []interface{}{
				"a message",
				errors.New("an error"),
				grpcerr.Kind("a kind"),
				zap.String("zap.string", "a string"),
				[]zapcore.Field{zap.String("zap.int64", "in integer")},
			},
			exp: "a kind: a message: an error",
		},
		"nested-error": {
			args: []interface{}{grpcerr.E("a message", errors.New("an error"), grpcerr.Kind("a kind"))},
			exp:  "a kind: a message: an error",
		},
		"nested-error-with-same-properties": {
			args: []interface{}{grpcerr.E("a message", errors.New("an error"), grpcerr.Kind("a kind")), grpcerr.Kind("a kind")},
			exp:  "a kind:\n\ta message: an error",
		},
		"nested-error-unknown-kind": {
			args: []interface{}{grpcerr.E("a message", errors.New("an error"))},
			exp:  "a message: an error",
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			got := grpcerr.E(c.args...)
			if got.Error() != c.exp {
				t.Errorf("wrong output, expected '%s' but got '%s'", c.exp, got.Error())
			}
		})
	}
}

const (
	op = grpcerr.Op("Op")
)

func TestMatch(t *testing.T) {
	cases := []struct {
		err1, err2 error
		matched    bool
	}{
		// Errors not of type *Error fail outright.
		{nil, nil, false},
		{io.EOF, io.EOF, false},
		{grpcerr.E(io.EOF), io.EOF, false},
		{io.EOF, grpcerr.E(io.EOF), false},
		// Success. We can drop fields from the first argument and still match.
		{grpcerr.E(io.EOF), grpcerr.E(io.EOF), true},
		{grpcerr.E(op, grpcerr.Kind("invalid"), io.EOF), grpcerr.E(op, grpcerr.Kind("invalid"), io.EOF), true},
		{grpcerr.E(op), grpcerr.E(op, grpcerr.Kind("invalid"), io.EOF), true},
		{grpcerr.E(grpcerr.Kind("invalid")), grpcerr.E(op, grpcerr.Kind("invalid"), io.EOF), true},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			matched := grpcerr.Match(c.err1, c.err2)
			if matched != c.matched {
				t.Errorf("Match(%q, %q)=%t; want %t", c.err1, c.err2, matched, c.matched)
			}
		})
	}
}

func TestIs(t *testing.T) {
	cases := []struct {
		err     error
		kind    grpcerr.Kind
		matched bool
	}{
		{io.EOF, grpcerr.Other, false},
		{grpcerr.E("a message", grpcerr.Kind("a kind")), grpcerr.Other, false},
		{grpcerr.E("a message", grpcerr.Kind("a kind")), grpcerr.Kind("a kind"), true},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			matched := grpcerr.Is(c.kind, c.err)
			if matched != c.matched {
				t.Errorf("Is(%q, %q)=%t; want %t", c.kind, c.err, matched, c.matched)
			}
		})
	}
}
