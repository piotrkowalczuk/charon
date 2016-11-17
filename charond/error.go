package charond

import (
	"database/sql"

	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var errf = grpc.Errorf

func handleError(err error) error {
	switch {
	case err == sql.ErrNoRows:
		err = errf(codes.NotFound, "charond: entity does not exists")
	default:
		if pqerr, ok := err.(*pq.Error); ok {
			switch {
			case pqerr.Code == pq.ErrorCode("23502"):
				return errf(codes.InvalidArgument, "charond: %s cannot be empty", pqerr.Column)
			default:
				err = errf(codes.Internal, pqerr.Message)
			}
		}
	}
	return err
}
