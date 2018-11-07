package charond

import (
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleMnemosyneError(err error) error {
	if sts, ok := status.FromError(err); ok && sts.Code() == codes.NotFound {
		return grpcerr.E(codes.Unauthenticated, "session not found")
	}

	return grpcerr.E(codes.Internal, "session fetch failure", err)
}
