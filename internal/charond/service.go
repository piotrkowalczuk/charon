package charond

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func initPostgres(address string, test bool, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("connection failure: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		cancel := time.NewTimer(10 * time.Second)
		attempts := 1
	PingLoop:
		for {
			select {
			case <-time.After(1 * time.Second):
				if err := db.Ping(); err != nil {
					attempts++
					continue PingLoop
				}
				break PingLoop
			case <-cancel.C:
				return nil, fmt.Errorf("postgres connection failed after %d attempts", attempts)
			}
		}
	}

	if test {
		if err = teardownDatabase(db); err != nil {
			return nil, err
		}
		logger.Info("database has been cleared upfront")
	}
	if err = setupDatabase(db); err != nil {
		return nil, err
	}

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	username := ""
	if u.User != nil {
		username = u.User.Username()
	}

	logger.Info("postgres connection has been established", zap.String("host", u.Host), zap.String("username", username))

	return db, nil
}

func initMnemosyne(address string, logger *zap.Logger, opts []grpc.DialOption) (mnemosynerpc.SessionManagerClient, *grpc.ClientConn) {
	if address == "" {
		logger.Error("missing mnemosyne address")
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Error("mnemosyne dial falilure", zap.Error(err), zap.String("address", address))
	}

	logger.Info("rpc connection to mnemosyne has been established", zap.String("address", address))

	return mnemosynerpc.NewSessionManagerClient(conn), conn
}

func initHasher(cost int, logger *zap.Logger) password.Hasher {
	bh, err := password.NewBCryptHasher(cost)
	if err != nil {
		logger.Fatal("hasher initialization failure", zap.Error(err))
	}

	return bh
}

func initPermissionRegistry(r model.PermissionProvider, permissions charon.Permissions, logger *zap.Logger) (pr model.PermissionRegistry) {
	pr = model.NewPermissionRegistry(r)
	created, untouched, removed, err := pr.Register(context.TODO(), permissions)
	if err != nil {
		logger.Fatal("permission registry initialization failure", zap.Error(err))
	}

	logger.Info("charon permissions has been registered", zap.Int64("created", created), zap.Int64("untouched", untouched), zap.Int64("removed", removed))

	return
}
