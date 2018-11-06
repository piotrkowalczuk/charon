package model

import (
	"context"
	"database/sql"
)

// RefreshTokenProvider ...
type RefreshTokenProvider interface {
	// Find ...
	Find(context.Context, *RefreshTokenFindExpr) ([]*RefreshTokenEntity, error)
	// FindOneByTokenAndUserID ...
	FindOneByTokenAndUserID(context.Context, string, int64) (*RefreshTokenEntity, error)
	// Create ...
	Create(context.Context, *RefreshTokenEntity) (*RefreshTokenEntity, error)
	// UpdateOneByTokenAndUserID ...
	UpdateOneByTokenAndUserID(context.Context, string, int64, *RefreshTokenPatch) (*RefreshTokenEntity, error)
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
