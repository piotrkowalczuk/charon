package charonctl

import (
	"context"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/mnemosyne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ConsoleOpts struct {
	Conn     *grpc.ClientConn
	Username string
	Password string
}

type Console struct {
	Ctx context.Context
	consoleRegisterUser
	consoleObtainRefreshToken
}

func NewConsole(opts ConsoleOpts) (*Console, error) {
	auth := charonrpc.NewAuthClient(opts.Conn)
	user := charonrpc.NewUserManagerClient(opts.Conn)
	refreshToken := charonrpc.NewRefreshTokenManagerClient(opts.Conn)

	c := &Console{
		Ctx: context.Background(),
		consoleObtainRefreshToken: consoleObtainRefreshToken{
			refreshToken: refreshToken,
		},
		consoleRegisterUser: consoleRegisterUser{
			user: user,
		},
	}

	ctx := context.Background()

	if opts.Password != "" || opts.Username != "" {
		resp, err := auth.Login(context.Background(), &charonrpc.LoginRequest{
			Username: opts.Username,
			Password: opts.Password,
		})
		if err != nil {
			return nil, &Error{
				Msg: "(initial) login failure",
				Err: err,
			}
		}

		c.Ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(mnemosyne.AccessTokenMetadataKey, resp.Value))
	}

	return c, nil
}
