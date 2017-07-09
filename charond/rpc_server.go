package charond

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc/metadata"
)

type rpcServer struct {
	opts               DaemonOpts
	meta               metadata.MD
	logger             log.Logger
	ldap               *sync.Pool
	session            mnemosynerpc.SessionManagerClient
	passwordHasher     password.Hasher
	permissionRegistry model.PermissionRegistry
	repository         repositories
}
