package session

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

type Actor struct {
	User        *model.UserEntity
	Session     *mnemosynerpc.Session
	Permissions charon.Permissions
	IsLocal     bool
}
