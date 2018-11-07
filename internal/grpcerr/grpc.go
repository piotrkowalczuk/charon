package grpcerr

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		res, err = handler(ctx, req)
		return res, grpcError(err)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return grpcError(handler(srv, ss))
	}
}

func grpcError(err error) error {
	ert, ok := err.(*Error)
	if !ok {
		return err
	}

	if ert.Msg != "" {
		sts, err := status.Newf(ert.Code, "%s: %s", ert.Msg, ert.Err).WithDetails(ert.Details...)
		// error is not nil only if code was OK
		if err != nil {
			return status.Newf(ert.Code, "%s: %s", ert.Msg, ert.Err).Err()
		}

		return sts.Err()
	} else {
		sts, err := status.Newf(ert.Code, "%s", ert.Err).WithDetails(ert.Details...)
		// error is not nil only if code was OK
		if err != nil {
			return status.Newf(ert.Code, "%s", ert.Err).Err()
		}

		return sts.Err()
	}
}
