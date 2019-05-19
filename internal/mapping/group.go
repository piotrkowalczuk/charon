package mapping

import (
	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/model"
)

// ReverseGroup maps internal entity struct into protobuf message used by a client.
func ReverseGroup(ent *model.GroupEntity) (*charonrpc.Group, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ent.CreatedAt); err != nil {
		return nil, err
	}
	if ent.UpdatedAt.Valid {
		if updatedAt, err = ptypes.TimestampProto(ent.UpdatedAt.Time); err != nil {
			return nil, err
		}
	}

	return &charonrpc.Group{
		Id:          ent.ID,
		Name:        ent.Name,
		Description: ent.Description.Chars,
		CreatedAt:   createdAt,
		CreatedBy:   &ent.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   &ent.UpdatedBy,
	}, nil
}

// ReverseGroups does same thing like ReverseGroup but operate on slices.
func ReverseGroups(in []*model.GroupEntity) ([]*charonrpc.Group, error) {
	res := make([]*charonrpc.Group, 0, len(in))
	for _, ent := range in {
		msg, err := ReverseGroup(ent)
		if err != nil {
			return nil, err
		}
		res = append(res, msg)
	}

	return res, nil
}
