package charond

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc"
)

func initPostgres(address string, test bool, logger log.Logger) (*sql.DB, error) {
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
		sklog.Info(logger, "database has been cleared upfront")
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

	sklog.Info(logger, "postgres connection has been established", "host", u.Host, "username", username)

	return db, nil
}

func initMnemosyne(address string, logger log.Logger, opts []grpc.DialOption) (mnemosynerpc.SessionManagerClient, *grpc.ClientConn) {
	if address == "" {
		sklog.Fatal(logger, errors.New("missing mnemosyne address"))

	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		sklog.Fatal(logger, err, "address", address)
	}

	sklog.Info(logger, "rpc connection to mnemosyne has been established", "address", address)

	return mnemosynerpc.NewSessionManagerClient(conn), conn
}

func initHasher(cost int, logger log.Logger) password.Hasher {
	bh, err := password.NewBCryptHasher(cost)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return bh
}

func initPermissionRegistry(r model.PermissionProvider, permissions charon.Permissions, logger log.Logger) (pr model.PermissionRegistry) {
	pr = model.NewPermissionRegistry(r)
	created, untouched, removed, err := pr.Register(context.TODO(), permissions)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	sklog.Info(logger, "charon permissions has been registered", "created", created, "untouched", untouched, "removed", removed)

	return
}
