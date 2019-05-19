package charonctl

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/ntypes"
)

type ObtainRefreshTokenArg struct {
	Notes       string
	ExpireAfter time.Duration
}

type consoleObtainRefreshToken struct {
	refreshToken charonrpc.RefreshTokenManagerClient
}

func (cort *consoleObtainRefreshToken) ObtainRefreshToken(ctx context.Context, arg *ObtainRefreshTokenArg) error {
	var (
		expireAt *timestamp.Timestamp
		err      error
	)
	if arg.ExpireAfter != time.Duration(0) {
		if expireAt, err = ptypes.TimestampProto(time.Now().Add(arg.ExpireAfter)); err != nil {
			return err
		}
	}

	_, err = cort.refreshToken.Create(ctx, &charonrpc.CreateRefreshTokenRequest{
		Notes:    &ntypes.String{Chars: arg.Notes, Valid: arg.Notes != ""},
		ExpireAt: expireAt,
	})
	if err != nil {
		return err
	}

	res, err := cort.refreshToken.List(ctx, &charonrpc.ListRefreshTokensRequest{})
	if err != nil {
		return err
	}
	for _, rt := range res.RefreshTokens {
		eat, err := ptypes.Timestamp(rt.ExpireAt)
		if err != nil {
			fmt.Printf("%-36s ", "never")
		} else {
			fmt.Printf("%-36s ", eat.String())
		}

		fmt.Printf("%s", rt.Token)
		if rt.Notes != nil && rt.Notes.Valid {
			fmt.Printf(" - %s", rt.Notes.StringOr(""))
		}
		fmt.Print("\n")
	}
	return nil
}
