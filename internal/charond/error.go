package charond

import (
	"database/sql"

	"github.com/lib/pq"
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

func handleError(err error) error {
	switch {
	case err == sql.ErrNoRows:
		err = grpcerr.E(codes.NotFound, "charond: entity does not exists")
	default:
		if pqerr, ok := err.(*pq.Error); ok {
			switch {
			case pqerr.Code == pq.ErrorCode("23502"):
				return grpcerr.E(codes.InvalidArgument, "charond: %s cannot be empty", pqerr.Column)
			default:
				err = grpcerr.E(codes.Internal, pqerr.Message)
			}
		}
	}
	return err
}
