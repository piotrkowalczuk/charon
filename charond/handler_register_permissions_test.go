package charond

import (
	"testing"

	"context"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegisterPermissionsHandler_Register_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]struct {
		permissions []string
		assert      func(*testing.T, error)
	}{
		"simple": {
			permissions: []string{
				"a:b:c",
				"a:bb:cc",
			},
			assert: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			},
		},
		"inconsistent-subsystem": {
			permissions: []string{
				"a:b:c",
				"aa:bb:cc",
			},
			assert: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("missing error")
				}
				if st, ok := status.FromError(err); ok {
					if st.Code() != codes.InvalidArgument {
						t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
					}
				}
			},
		},
		"empty-subsystem": {
			permissions: []string{
				":b:c",
			},
			assert: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("missing error")
				}
				if st, ok := status.FromError(err); ok {
					if st.Code() != codes.InvalidArgument {
						t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
					}
				}
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			_, err := suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
				Permissions: c.permissions,
			})
			c.assert(t, err)
		})
	}
}
func TestRegisterPermissionsHandler_Register_Unit(t *testing.T) {
	registryMock := &modelmock.PermissionRegistry{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.RegisterPermissionsRequest
		err  error
	}{
		"ok": {
			init: func(t *testing.T) {
				registryMock.On("Register", mock.Anything, mock.Anything).
					Return(int64(1), int64(2), int64(3), nil).
					Once()
			},
			req: charonrpc.RegisterPermissionsRequest{Permissions: []string{
				"a:b:c",
				"a:bb:cc",
			}},
		},
		"empty-subsystem": {
			init: func(t *testing.T) {
				zero := int64(0)
				registryMock.On("Register", mock.Anything, mock.Anything).
					Return(zero, zero, zero, model.ErrEmptySubsystem).
					Once()
			},
			req: charonrpc.RegisterPermissionsRequest{Permissions: []string{
				"a:b:c",
				"a:bb:cc",
			}},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"request-cancel": {
			init: func(t *testing.T) {
				zero := int64(0)
				registryMock.On("Register", mock.Anything, mock.Anything).
					Return(zero, zero, zero, context.Canceled).
					Once()
			},
			req: charonrpc.RegisterPermissionsRequest{Permissions: []string{
				"a:b:c",
				"a:bb:cc",
			}},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := registerPermissionsHandler{
		registry: registryMock,
	}
	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			registryMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Register(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, registryMock)
		})
	}
}
