// +build unit !postgres

package main

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/piotrkowalczuk/charon"
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
		tkn mnemosyne.Token
	)

	Convey("retrieveActor", t, func() {
		userRepositoryMock := &MockUserRepository{}
		permissionRepositoryMock := &MockPermissionRepository{}
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

			Convey("should return an error", func() {
				act, err = h.retrieveActor(ctx)

				So(err, ShouldNotBeNil)
				So(act, ShouldBeNil)
			})
		})
		Convey("As authenticated user", func() {
			id = 7856282
			tkn = mnemosyne.NewToken([]byte("0000000001"), []byte("hash"))
			ctx = mnemosyne.NewTokenContext(context.Background(), tkn)
			sessionMock.On("FromContext", ctx).
				Once().
				Return(&mnemosyne.Session{
				Token:     &tkn,
				SubjectId: charon.SubjectIDFromInt64(id).String(),
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
							Subsystem: charon.PermissionCanRetrieve.Subsystem(),
							Module:    charon.PermissionCanRetrieve.Module(),
							Action:    charon.PermissionCanRetrieve.Action(),
						},
						{
							Subsystem: charon.UserCanRetrieveAsOwner.Subsystem(),
							Module:    charon.UserCanRetrieveAsOwner.Module(),
							Action:    charon.UserCanRetrieveAsOwner.Action(),
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
