package charon

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/piotrkowalczuk/mnemosyne"
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
		tkn mnemosyne.AccessToken
	)

	Convey("retrieveActor", t, func() {
		userRepositoryMock := &mockUserProvider{}
		permissionRepositoryMock := &mockPermissionProvider{}
		sessionMock := &mnemosynetest.Mnemosyne{}
		h := &handler{
			session: sessionMock,
		}
		h.repository.user = userRepositoryMock
		h.repository.permission = permissionRepositoryMock

		Convey("As unauthenticated user", func() {
			ctx = context.Background()
			sessionMock.On("FromContext", mock.Anything).
				Once().
				Return(nil, errors.New("mnemosyned: test error"))

			Convey("Should return an error", func() {
				act, err = h.retrieveActor(ctx)

				So(err, ShouldNotBeNil)
				So(act, ShouldBeNil)
			})
		})
		Convey("As authenticated user", func() {
			id = 7856282
			tkn = mnemosyne.NewAccessToken([]byte("0000000001"), []byte("hash"))
			ctx = mnemosyne.NewAccessTokenContext(context.Background(), tkn)
			sessionMock.On("FromContext", ctx).
				Once().
				Return(&mnemosyne.Session{
					AccessToken: &tkn,
					SubjectId:   SubjectIDFromInt64(id).String(),
				}, nil)

			Convey("When user exists", func() {
				userRepositoryMock.On("FindOneByID", id).
					Once().
					Return(&userEntity{ID: id}, nil)

				Convey("And it has some permissions", func() {
					permissionRepositoryMock.On("FindByUserID", id).
						Once().
						Return([]*permissionEntity{
							{
								Subsystem: PermissionCanRetrieve.Subsystem(),
								Module:    PermissionCanRetrieve.Module(),
								Action:    PermissionCanRetrieve.Action(),
							},
							{
								Subsystem: UserCanRetrieveAsOwner.Subsystem(),
								Module:    UserCanRetrieveAsOwner.Module(),
								Action:    UserCanRetrieveAsOwner.Action(),
							},
						}, nil)

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
						Once().
						Return(nil, sql.ErrNoRows)

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
