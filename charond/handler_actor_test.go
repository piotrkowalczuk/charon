package charond

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/mnemosyne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestActorHandler_Actor(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"correct-token": func(t *testing.T) {
			md, ok := metadata.FromContext(ctx)
			if !ok {
				t.Fatal("metadata not present in context")
			}
			tok := &wrappers.StringValue{Value: md[mnemosyne.AccessTokenMetadataKey][0]}

			res, err := suite.charon.auth.Actor(context.Background(), tok)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Username != "test" {
				t.Errorf("wrong username, expected %s but got %s", "test", res.Username)
			}
		},
		"incorrect-token": func(t *testing.T) {
			_, err := suite.charon.auth.Actor(context.Background(), &wrappers.StringValue{Value: "incorrect-token"})
			if grpc.Code(err) != codes.Unauthenticated {
				t.Fatalf("wrong status code, expected %s but got %s", codes.Unauthenticated.String(), grpc.Code(err).String())
			}
		},
		"context": func(t *testing.T) {
			res, err := suite.charon.auth.Actor(timeout(ctx), &wrappers.StringValue{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Username != "test" {
				t.Errorf("wrong username, expected %s but got %s", "test", res.Username)
			}
		},
		"nothing": func(t *testing.T) {
			_, err := suite.charon.auth.Actor(context.Background(), &wrappers.StringValue{})
			if grpc.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}
