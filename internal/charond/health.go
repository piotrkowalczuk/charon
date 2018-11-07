package charond

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type healthHandler struct {
	logger   *zap.Logger
	postgres *sql.DB
}

func (hh *healthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if hh.postgres != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := hh.postgres.PingContext(ctx); err != nil {
			hh.logger.Debug("health check failure due to postgres connection")
			http.Error(rw, "postgres ping failure", http.StatusServiceUnavailable)
			return
		}
	}

	hh.logger.Debug("successful health check")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("1"))
}
