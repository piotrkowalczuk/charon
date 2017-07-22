package charond

import (
	"context"
	"database/sql"
	"net"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

func TestDeleteUserHandler_Delete(t *testing.T) {
	upm := &model.MockUserProvider{}
	ppm := &model.MockPermissionProvider{}
	sm := &mnemosynetest.SessionManagerClient{}

	sessionOnContext := func(t *testing.T, id int64) {
		sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(&mnemosynerpc.ContextResponse{
			Session: &mnemosynerpc.Session{
				SubjectId: session.ActorIDFromInt64(id).String(),
				AccessToken: func() string {
					tkn, err := mnemosyne.RandomAccessToken()
					if err != nil {
						t.Fatalf("token generation error: %s", err.Error())
					}
					return tkn
				}(),
			},
		}, nil).Once()
	}

	h := deleteUserHandler{
		handler: &handler{
			session: sm,
			repository: repositories{
				user:       upm,
				permission: ppm,
			},
		},
	}

	cases := map[string]func(t *testing.T){
		"invalid-id": func(t *testing.T) {
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: -1})
			assertErrorCode(t, err, codes.InvalidArgument, "user cannot be deleted, invalid id: -1")
		},
		"cannot-remove-if-does-not-exists": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(nil, sql.ErrNoRows).
				Times(2)
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 10})
			assertErrorCode(t, err, codes.PermissionDenied, "actor does not exists")
		},
		"cannot-remove-from-localhost": func(t *testing.T) {
			ctx := peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{
					IP: net.IPv4(127, 0, 0, 1),
				},
			})
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, errf(codes.Unknown, "random mnemosyne error")).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(nil, sql.ErrNoRows).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(ctx, &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "user cannot be removed from localhost")
		},
		"cannot-if-session-does-not-exists": func(t *testing.T) {
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, errf(codes.NotFound, "session does not exists")).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(nil, sql.ErrNoRows).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.Unauthenticated, "session does not exists")
		},
		"cannot-remove-permission-query-timeout": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return(nil, context.DeadlineExceeded).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			if err == nil {
				t.Fatal("expected error")
			}
			if err != context.DeadlineExceeded {
				t.Fatalf("wrong error, expected '%s' but got '%s'", context.DeadlineExceeded.Error(), err.Error())
			}
		},
		"cannot-remove-itself": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Times(2)
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 10})
			assertErrorCode(t, err, codes.PermissionDenied, "user is not permitted to remove himself")
		},
		"can-remove-as-superuser": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10, IsSuperuser: true}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
				Once()
			upm.On("DeleteOneByID", mock.Anything, int64(11)).Return(int64(1), nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"cannot-remove-as-stranger": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteAsOwner.Subsystem(),
					Module:    charon.UserCanDeleteAsOwner.Module(),
					Action:    charon.UserCanDeleteAsOwner.Action(),
				},
			}, nil).Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "user cannot be removed by stranger, missing permission")
		},
		"cannot-remove-as-owner": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
				ID:        11,
				CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
			}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteAsStranger.Subsystem(),
					Module:    charon.UserCanDeleteAsStranger.Module(),
					Action:    charon.UserCanDeleteAsStranger.Action(),
				},
			}, nil).Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "user cannot be removed by owner, missing permission")
		},
		"can-delete-as-stranger-but-does-not-exists": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(nil, sql.ErrNoRows).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteAsStranger.Subsystem(),
					Module:    charon.UserCanDeleteAsStranger.Module(),
					Action:    charon.UserCanDeleteAsStranger.Action(),
				},
			}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.NotFound, "user does not exists")
		},
		"can-delete-as-stranger-but-not-a-superuser": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
				ID:          11,
				IsSuperuser: true,
			}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteAsStranger.Subsystem(),
					Module:    charon.UserCanDeleteAsStranger.Module(),
					Action:    charon.UserCanDeleteAsStranger.Action(),
				},
			}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "only superuser can remove other superuser")
		},
		"can-delete-as-stranger-but-not-a-staff-member": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
				ID:      11,
				IsStaff: true,
			}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteAsStranger.Subsystem(),
					Module:    charon.UserCanDeleteAsStranger.Module(),
					Action:    charon.UserCanDeleteAsStranger.Action(),
				},
			}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "staff user cannot be removed by stranger, missing permission")
		},
		"can-delete-staff-member-but-not-as-a-stranger": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
				ID:      11,
				IsStaff: true,
			}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteStaffAsOwner.Subsystem(),
					Module:    charon.UserCanDeleteStaffAsOwner.Module(),
					Action:    charon.UserCanDeleteStaffAsOwner.Action(),
				},
			}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "staff user cannot be removed by stranger, missing permission")
		},
		"can-delete-staff-member-but-not-as-a-owner": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
				ID:        11,
				IsStaff:   true,
				CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
			}, nil).
				Once()
			ppm.On("FindByUserID", mock.Anything, int64(10)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserCanDeleteStaffAsStranger.Subsystem(),
					Module:    charon.UserCanDeleteStaffAsStranger.Module(),
					Action:    charon.UserCanDeleteStaffAsStranger.Action(),
				},
			}, nil).
				Once()
			sessionOnContext(t, 10)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteUserRequest{Id: 11})
			assertErrorCode(t, err, codes.PermissionDenied, "staff user cannot be removed by owner, missing permission")
		},
	}

	for hint, c := range cases {
		// reset mocks between cases
		sm.ExpectedCalls = []*mock.Call{}
		ppm.ExpectedCalls = []*mock.Call{}
		upm.ExpectedCalls = []*mock.Call{}

		t.Run(hint, c)
	}
}

func TestDeleteUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteUserRequest
		act session.Actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserCanDeleteAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: model.UserEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: model.UserEntity{
				ID: 2,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID: 1,
				},
				Permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID: 1,
				},
				Permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
			},
		},
	}

	h := &deleteUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err != nil {
			t.Errorf("unexpected error for %d: %s", i, err.Error())
		}
	}
}

func TestDeleteUserHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteUserRequest
		act session.Actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{},
			},
			ent: model.UserEntity{},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
			ent: model.UserEntity{
				ID: 2,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID: 1,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1, IsSuperuser: true},
				Permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			},
		},
	}

	h := &deleteUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err == nil {
			t.Errorf("expected error for %d, got nil", i)
		}
	}
}
