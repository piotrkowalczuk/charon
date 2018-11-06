package grpcerr

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"context"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

type Error struct {
	// The underlying error that triggered this one, if any.
	Err    error
	Msg    string
	Fields []zapcore.Field
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	Op   Op
	Code codes.Code
	// Details can be mapped into grpc status details.
	Details []proto.Message
}

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == "" && e.Err == nil
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	errorek.Op
//		The operation being performed, usually the method
//		being invoked (Get, Put, etc.).
//	codes.Code
//		GRPC response status code.
//	proto.Message
//		GRPC response status details.
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.Str. To avoid a common
//		class of misuse, if the string contains an @, it will be
//		treated as a PathName or UserName, as appropriate. Use
//		errors.Str explicitly to avoid this special-casing.
//	errorek.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func E(args ...interface{}) *Error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case Op:
			e.Op = arg
		case codes.Code:
			e.Code = arg
		case proto.Message:
			e.Details = append(e.Details, arg)
		case zapcore.Field:
			e.Fields = append(e.Fields, arg)
		case []zapcore.Field:
			e.Fields = append(e.Fields, arg...)
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		case string:
			e.Msg = arg
		default:
			// TODO: implement
			//_, file, line, _ := runtime.Caller(1)
			//log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return E(fmt.Sprintf("unknown type %T, value %v in error call", arg, arg))
		}
	}

	if e.Err == nil && e.Msg != "" {
		e.Err = errors.New(e.Msg)
		e.Msg = ""
	}

	switch e.Err {
	case context.DeadlineExceeded:
		e.Code = codes.DeadlineExceeded
	case context.Canceled:
		e.Code = codes.Canceled
	}

	if sts, ok := status.FromError(e.Err); ok {
		if sts.Code() != codes.OK && sts.Code() != codes.Unknown {
			e.Code = sts.Code()
		}
	}

	e.Fields = append(e.Fields, zap.Error(e.Err))

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, file name or user name
	// twice.
	if prev.Kind == e.Kind {
		prev.Kind = Other
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	return e
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}
	if e.Code != 0 {
		pad(b, ": ")
		b.WriteString(e.Code.String())
	}
	if e.Kind != Other && e.Kind != "" {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Msg != "" {
		pad(b, ": ")
		b.WriteString(e.Msg)
	}
	if e.Err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// Match compares its two error arguments. It can be used to check
// for expected errors in tests. Both arguments must have underlying
// type *Error or Match will return false. Otherwise it returns true
// iff every non-zero element of the first error is equal to the
// corresponding element of the second.
// If the Err field is a *Error, Match recurs on that field;
// otherwise it compares the strings returned by the Error methods.
// Elements that are in the second argument but not present in
// the first are ignored.
//
// For example,
//	Match(errors.E(errors.Permission), err)
// tests whether err is an Error with Kind=Permission.
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != "" && e2.Kind != e1.Kind {
		return false
	}
	if e2.Code != e1.Code {
		return false
	}
	if !reflect.DeepEqual(e1.Details, e2.Details) {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// String implements fmt Stringer interface.
func (o Op) String() string {
	return string(o)
}

var Other Kind = "other"

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind string

// String implements fmt Stringer interface.
func (k Kind) String() string {
	return string(k)
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
