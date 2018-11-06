package mapping

import (
	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/ntypes"
)

func ReverseRefreshToken(ent *model.RefreshTokenEntity) (*charonrpc.RefreshToken, error) {
	var (
		err                                        error
		expireAt, lastUsedAt, createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ent.CreatedAt); err != nil {
		return nil, err
	}
	if ent.UpdatedAt.Valid {
		if updatedAt, err = ptypes.TimestampProto(ent.UpdatedAt.Time); err != nil {
			return nil, err
		}
	}
	if ent.ExpireAt.Valid {
		if expireAt, err = ptypes.TimestampProto(ent.ExpireAt.Time); err != nil {
			return nil, err
		}
	}
	if ent.LastUsedAt.Valid {
		if lastUsedAt, err = ptypes.TimestampProto(ent.LastUsedAt.Time); err != nil {
			return nil, err
		}
	}

	return &charonrpc.RefreshToken{
		Token:      ent.Token,
		Notes:      &ent.Notes,
		Revoked:    ent.Revoked,
		ExpireAt:   expireAt,
		LastUsedAt: lastUsedAt,
		UserId:     ent.UserID,
		CreatedAt:  createdAt,
		CreatedBy:  &ent.CreatedBy,
		UpdatedAt:  updatedAt,
		UpdatedBy:  &ent.UpdatedBy,
	}, nil
}

func ReverseRefreshTokens(in []*model.RefreshTokenEntity) ([]*charonrpc.RefreshToken, error) {
	res := make([]*charonrpc.RefreshToken, 0, len(in))
	for _, ent := range in {
		msg, err := ReverseRefreshToken(ent)
		if err != nil {
			return nil, err
		}
		res = append(res, msg)
	}

	return res, nil
}

func RefreshTokenQuery(q *charonrpc.RefreshTokenQuery) *model.RefreshTokenCriteria {
	var revoked ntypes.Bool
	if q.GetRevoked() != nil {
		revoked = *q.GetRevoked()
	}

	return &model.RefreshTokenCriteria{
		UserID:     q.GetUserId(),
		Notes:      q.GetNotes(),
		Revoked:    revoked, // TODO: pointer?
		ExpireAt:   q.GetExpireAt(),
		LastUsedAt: q.GetLastUsedAt(),
		CreatedAt:  q.GetCreatedAt(),
		UpdatedAt:  q.GetUpdatedAt(),
	}
}
