package mapping_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/ntypes"
)

func TestReverseUser(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		given    model.UserEntity
		expected charonrpc.User
	}{
		"empty": {
			given: model.UserEntity{},
			expected: charonrpc.User{
				UpdatedBy: &ntypes.Int64{},
				CreatedBy: &ntypes.Int64{},
			},
		},
		"simple": {
			given: model.UserEntity{
				ConfirmationToken: []byte("confirmation-token"),
				CreatedAt:         now,
				CreatedBy:         ntypes.Int64{Int64: 1, Valid: true},
				UpdatedBy:         ntypes.Int64{Int64: 2, Valid: true},
				FirstName:         "firstname",
				ID:                1,
				LastLoginAt:       pq.NullTime{Time: now, Valid: true},
				LastName:          "lastname",
				Password:          []byte("password"),
				UpdatedAt:         pq.NullTime{Time: now.Add(1 * time.Hour), Valid: true},
			},
			expected: charonrpc.User{
				Id:        1,
				FirstName: "firstname",
				LastName:  "lastname",
				UpdatedBy: &ntypes.Int64{Int64: 2, Valid: true},
				CreatedBy: &ntypes.Int64{Int64: 1, Valid: true},
				CreatedAt: func() *timestamp.Timestamp {
					ts, _ := ptypes.TimestampProto(now)
					return ts
				}(),
				UpdatedAt: func() *timestamp.Timestamp {
					ts, _ := ptypes.TimestampProto(now.Add(1 * time.Hour))
					return ts
				}(),
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			got, err := mapping.ReverseUser(&c.given)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !reflect.DeepEqual(*got, c.expected) {
				t.Errorf("wrong output, expected:\n	%v\nbut got:\n	%v", c.expected, *got)
			}
		})
	}
}

func TestReverseUsers_valid(t *testing.T) {
	now := time.Now()
	entities := []*model.UserEntity{
		{
			ConfirmationToken: []byte("confirmation-token-1"),
			CreatedAt:         now,
			CreatedBy:         ntypes.Int64{Int64: 1, Valid: true},
			UpdatedBy:         ntypes.Int64{Int64: 2, Valid: true},
			FirstName:         "firstname-1",
			ID:                1,
			LastLoginAt:       pq.NullTime{Time: now, Valid: true},
			LastName:          "lastname-1",
			Password:          []byte("password-1"),
			UpdatedAt:         pq.NullTime{Time: now.Add(1 * time.Hour), Valid: true},
		},
		{
			ConfirmationToken: []byte("confirmation-token-2"),
			CreatedAt:         now,
			CreatedBy:         ntypes.Int64{Int64: 1, Valid: true},
			UpdatedBy:         ntypes.Int64{Int64: 2, Valid: true},
			FirstName:         "firstname-2",
			ID:                1,
			LastLoginAt:       pq.NullTime{Time: now, Valid: true},
			LastName:          "lastname-2",
			Password:          []byte("password-2"),
			UpdatedAt:         pq.NullTime{Time: now.Add(1 * time.Hour), Valid: true},
		},
	}
	users, err := mapping.ReverseUsers(entities)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != len(entities) {
		t.Error("wrong output length")
	}
}

func TestReverseUsers_invalid(t *testing.T) {
	const maxValidSeconds = 253402300800
	now := time.Now()
	cases := [][]*model.UserEntity{
		{
			{
				ConfirmationToken: []byte("confirmation-token-3"),
				CreatedAt:         time.Unix(maxValidSeconds, 0).UTC(),
				CreatedBy:         ntypes.Int64{Int64: 1, Valid: true},
				UpdatedBy:         ntypes.Int64{Int64: 2, Valid: true},
				FirstName:         "firstname-3",
				ID:                1,
				LastLoginAt:       pq.NullTime{Time: now, Valid: true},
				LastName:          "lastname-3",
				Password:          []byte("password-3"),
				UpdatedAt:         pq.NullTime{Valid: true},
			},
		},
		{
			{
				ConfirmationToken: []byte("confirmation-token-3"),
				CreatedBy:         ntypes.Int64{Int64: 1, Valid: true},
				UpdatedBy:         ntypes.Int64{Int64: 2, Valid: true},
				FirstName:         "firstname-3",
				ID:                1,
				LastLoginAt:       pq.NullTime{Time: now, Valid: true},
				LastName:          "lastname-3",
				Password:          []byte("password-3"),
				UpdatedAt:         pq.NullTime{Time: time.Unix(maxValidSeconds, 0).UTC(), Valid: true},
			},
		},
	}
	for _, entities := range cases {
		t.Run("", func(t *testing.T) {
			if _, err := mapping.ReverseUsers(entities); err == nil {
				t.Error("error expected, got nil")
			}
		})
	}
}
