package charond

import (
	"context"
	"strconv"
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestIsAuthenticatedHandler_IsAuthenticated(t *testing.T) {
	cases := map[string]struct {
		req    charonrpc.IsAuthenticatedRequest
		ses    mnemosynerpc.Session
		expRes bool
		expErr error
	}{
		"basic": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			ses: mnemosynerpc.Session{
				AccessToken:   "123",
				SubjectId:     "charon:user:1",
				SubjectClient: "test",
			},
			expRes: true,
		},
		"no-access-token": {
			req:    charonrpc.IsAuthenticatedRequest{},
			expErr: grpc.Errorf(codes.InvalidArgument, "authentication status cannot be checked, missing token"),
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			sm := &mnemosynetest.SessionManagerClient{}
			up := &model.MockUserProvider{}

			h := &isAuthenticatedHandler{
				handler: &handler{
					session: sm,
					repository: repositories{
						user: up,
					},
				},
			}
			sm.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
				AccessToken: c.req.AccessToken,
			}, mock.Anything).Return(&mnemosynerpc.GetResponse{
				Session: &c.ses,
			}, nil)

			if c.expErr == nil {
				id, err := strconv.ParseInt(c.ses.SubjectId[12:], 10, 64)
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
				up.On("Exists", mock.Anything, id).Return(true, nil)
			}

			res, err := h.IsAuthenticated(context.TODO(), &c.req)
			if err != nil {
				if c.expErr == nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
				if c.expErr.Error() != err.Error() {
					t.Fatalf("wrong error: %s", err.Error())
				}
				return
			}
			if res.Value != c.expRes {
				t.Errorf("expected %t but got %t", c.expRes, res.Value)
			}
		})
	}
}
