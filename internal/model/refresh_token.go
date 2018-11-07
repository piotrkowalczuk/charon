package model

import (
	"context"
	"database/sql"
)

// RefreshTokenProvider ...
type RefreshTokenProvider interface {
	// Find ...
	Find(context.Context, *RefreshTokenFindExpr) ([]*RefreshTokenEntity, error)
	// FindOneByToken ...
	FindOneByToken(context.Context, string) (*RefreshTokenEntity, error)
	// Create ...
	Create(context.Context, *RefreshTokenEntity) (*RefreshTokenEntity, error)
	// UpdateOneByToken ...
	UpdateOneByToken(context.Context, string, *RefreshTokenPatch) (*RefreshTokenEntity, error)
	// FindOneByTokenAndUserID .
	FindOneByTokenAndUserID(ctx context.Context, token string, userID int64) (*RefreshTokenEntity, error)
}

// RefreshTokenRepository extends RefreshTokenRepositoryBase
type RefreshTokenRepository struct {
	RefreshTokenRepositoryBase
}

// NewRefreshTokenRepository ...
func NewRefreshTokenRepository(dbPool *sql.DB) RefreshTokenProvider {
	return &RefreshTokenRepository{
		RefreshTokenRepositoryBase: RefreshTokenRepositoryBase{
			DB:      dbPool,
			Table:   TableRefreshToken,
			Columns: TableRefreshTokenColumns,
		},
	}
}

// Create ...
func (rtr *RefreshTokenRepository) Create(ctx context.Context, ent *RefreshTokenEntity) (*RefreshTokenEntity, error) {
	return rtr.Insert(ctx, ent)
}

// FindOneByTokenAndUserID ...
func (rtr *RefreshTokenRepository) FindOneByTokenAndUserID(ctx context.Context, token string, userID int64) (*RefreshTokenEntity, error) {
	ent, err := rtr.FindOneByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if ent.UserID != userID {
		return nil, sql.ErrNoRows
	}
	return ent, nil
}
