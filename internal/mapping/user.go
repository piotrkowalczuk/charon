package mapping

import (
	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func ReverseUser(ent *model.UserEntity) (*charonrpc.User, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if !ent.CreatedAt.IsZero() {
		if createdAt, err = ptypes.TimestampProto(ent.CreatedAt); err != nil {
			return nil, err
		}
	}
	if ent.UpdatedAt.Valid {
		if updatedAt, err = ptypes.TimestampProto(ent.UpdatedAt.Time); err != nil {
			return nil, err
		}
	}

	return &charonrpc.User{
		Id:          ent.ID,
		Username:    ent.Username,
		FirstName:   ent.FirstName,
		LastName:    ent.LastName,
		IsSuperuser: ent.IsSuperuser,
		IsActive:    ent.IsActive,
		IsStaff:     ent.IsStaff,
		IsConfirmed: ent.IsConfirmed,
		CreatedAt:   createdAt,
		CreatedBy:   &ent.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   &ent.UpdatedBy,
	}, nil
}

func ReverseUsers(in []*model.UserEntity) ([]*charonrpc.User, error) {
	res := make([]*charonrpc.User, 0, len(in))
	for _, ent := range in {
		msg, err := ReverseUser(ent)
		if err != nil {
			return nil, err
		}
		res = append(res, msg)
	}

	return res, nil
}
