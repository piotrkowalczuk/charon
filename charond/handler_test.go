package charond

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestHandler(t *testing.T) {
	var (
		id  int64
		ctx context.Context
		err error
		act *actor
		tkn string
	)

	Convey("retrieveActor", t, func() {
		userRepositoryMock := &model.MockUserProvider{}
		permissionRepositoryMock := &model.MockPermissionProvider{}
		sessionMock := &mnemosynetest.SessionManagerClient{}
		h := &handler{
			session: sessionMock,
		}
		h.repository.user = userRepositoryMock
		h.repository.permission = permissionRepositoryMock

		Convey("As unauthenticated user", func() {
			ctx = context.Background()
			sessionMock.On("Context", mock.Anything, none(), mock.Anything).
				Return(nil, errors.New("mnemosyned: test error")).
				Once()

			Convey("Should return an error", func() {
				act, err = h.retrieveActor(ctx)

				So(err, ShouldNotBeNil)
				So(act, ShouldBeNil)
			})
		})
		Convey("As authenticated user", func() {
			id = 7856282
			tkn = "0000000001hash"
			ctx = mnemosyne.NewAccessTokenContext(context.Background(), tkn)
			sessionMock.On("Context", ctx, none(), mock.Anything).
				Return(&mnemosynerpc.ContextResponse{
					Session: &mnemosynerpc.Session{
						AccessToken: tkn,
						SubjectId:   session.ActorIDFromInt64(id).String(),
					},
				}, nil).
				Once()

			Convey("When user exists", func() {
				userRepositoryMock.On("FindOneByID", id).
					Return(&model.UserEntity{ID: id}, nil).
					Once()

				Convey("And it has some permissions", func() {
					permissionRepositoryMock.On("FindByUserID", id).
						Return([]*model.PermissionEntity{
							{
								Subsystem: charon.PermissionCanRetrieve.Subsystem(),
								Module:    charon.PermissionCanRetrieve.Module(),
								Action:    charon.PermissionCanRetrieve.Action(),
							},
							{
								Subsystem: charon.UserCanRetrieveAsOwner.Subsystem(),
								Module:    charon.UserCanRetrieveAsOwner.Module(),
								Action:    charon.UserCanRetrieveAsOwner.Action(),
							},
						}, nil).
						Once()

					Convey("Then it should be retrieved without any error", func() {
						act, err = h.retrieveActor(ctx)

						So(err, ShouldBeNil)
						So(act, ShouldNotBeNil)
						So(act.user.ID, ShouldEqual, id)
						So(act.permissions, ShouldHaveLength, 2)
					})
				})
				Convey("And it has no permissions", func() {
					permissionRepositoryMock.On("FindByUserID", id).
						Return(nil, sql.ErrNoRows).
						Once()

					Convey("Then it should be retrieved without any error", func() {
						act, err = h.retrieveActor(ctx)

						So(err, ShouldBeNil)
						So(act, ShouldNotBeNil)
						So(act.user.ID, ShouldEqual, id)
						So(act.permissions, ShouldHaveLength, 0)
					})
				})
			})
		})
	})
}
