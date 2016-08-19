package charond

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	monitoringRPCLabels = []string{
		"method",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	enabled  bool
	rpc      monitoringRPC
	postgres monitoringPostgres
}

type monitoringRPC struct {
	enabled  bool
	requests *prometheus.CounterVec
	errors   *prometheus.CounterVec
	duration *prometheus.SummaryVec
}

type monitoringPostgres struct {
	enabled bool
	queries *prometheus.CounterVec
	errors  *prometheus.CounterVec
}

func initUnaryServerInterceptor(monitor monitoringRPC) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		res, err := handler(ctx, req)
		if err != nil {
			err = interceptError(err)
			if monitor.enabled {
				elapsed := float64(time.Since(start)) / float64(time.Microsecond)
				labels := prometheus.Labels{
					"handler": info.FullMethod,
					"code":    grpc.Code(err).String(),
				}
				monitor.duration.With(labels).Observe(elapsed)
				monitor.errors.With(labels).Add(1)
			}

			return nil, err
		}
		if monitor.enabled {
			elapsed := float64(time.Since(start)) / float64(time.Microsecond)
			labels := prometheus.Labels{
				"handler": info.FullMethod,
				"code":    grpc.Code(err).String(),
			}
			monitor.duration.With(labels).Observe(elapsed)
			monitor.requests.With(labels).Add(1)
		}
		return res, nil
	})
}

func interceptError(err error) error {
	if err == nil {
		return nil
	}

	code := grpc.Code(err)
	switch code {
	case codes.Unknown:
		return grpc.Errorf(codes.Internal, "charond: %s", grpc.ErrorDesc(err))
	default:
		return grpc.Errorf(code, "charond: %s", grpc.ErrorDesc(err))
	}
}
