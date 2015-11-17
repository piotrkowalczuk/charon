package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
)

type rpcServer struct {
	logger         log.Logger
	monitor        *monitoring
	mnemosyne      mnemosyne.Mnemosyne
	passwordHasher charon.PasswordHasher
	userRepository UserRepository
}
