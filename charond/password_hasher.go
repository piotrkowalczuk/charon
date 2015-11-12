package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charond/lib/security"
	"github.com/piotrkowalczuk/sklog"
)

func initPasswordHasher(cost int, logger log.Logger) security.PasswordHasher {
	bh, err := charon.NewBcryptPasswordHasher(cost, logger)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return bh
}
