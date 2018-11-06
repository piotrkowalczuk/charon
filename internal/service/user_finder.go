package service

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon/internal/password"

	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"google.golang.org/grpc/codes"
)

type UserFinder interface {
	FindUser(context.Context) (*model.UserEntity, error)
}

type byUsernameAndPasswordUserFinder struct {
	username, password string
	userRepository     model.UserProvider
	hasher             password.Hasher
}

var _ UserFinder = &byUsernameAndPasswordUserFinder{}

func (f *byUsernameAndPasswordUserFinder) FindUser(ctx context.Context) (*model.UserEntity, error) {
	if f.username == "" {
		return nil, grpcerr.E(codes.InvalidArgument, "empty username")
	}
	if len(f.password) == 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "empty password")
	}

	usr, err := f.userRepository.FindOneByUsername(ctx, f.username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.Unauthenticated, "user with such username or password does not exists")
		}
		return nil, err
	}

	if bytes.Equal(usr.Password, model.ExternalPassword) {
		return nil, grpcerr.E(codes.FailedPrecondition, "authentication failure, external password manager not implemented")
	}
	if matches := f.hasher.Compare(usr.Password, []byte(f.password)); !matches {
		return nil, grpcerr.E(codes.Unauthenticated, "user with such username or password does not exists")
	}

	return usr, nil
}

type byRefreshTokenUserFinder struct {
	refreshToken           string
	userRepository         model.UserProvider
	refreshTokenRepository model.RefreshTokenProvider
}

func (f *byRefreshTokenUserFinder) FindUser(ctx context.Context) (*model.UserEntity, error) {
	if f.refreshToken == "" {
		return nil, grpcerr.E(codes.InvalidArgument, "empty refresh token")
	}

	refreshToken, err := f.refreshTokenRepository.FindOneByToken(ctx, f.refreshToken)
	if err != nil {
		return nil, err
	}
	user, err := f.userRepository.FindOneByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserFinderFactory struct {
	UserRepository         model.UserProvider
	RefreshTokenRepository model.RefreshTokenProvider
	Hasher                 password.Hasher
}

func (f *UserFinderFactory) ByUsernameAndPassword(username, password string) UserFinder {
	return &byUsernameAndPasswordUserFinder{
		username:       username,
		password:       password,
		userRepository: f.UserRepository,
		hasher:         f.Hasher,
	}
}

func (f *UserFinderFactory) ByRefreshToken(refreshToken string) UserFinder {
	return &byRefreshTokenUserFinder{
		refreshToken:           refreshToken,
		userRepository:         f.UserRepository,
		refreshTokenRepository: f.RefreshTokenRepository,
	}
}
