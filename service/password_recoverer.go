package service

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-soa/charon/lib"
	"github.com/go-soa/charon/lib/security"
)

var PasswordRecoverer lib.PasswordRecoverer

// InitPasswordRecoverer ...
func InitPasswordRecoverer(
	logger logrus.StdLogger,
	passwordHasher security.PasswordHasher,
	userRepository lib.UserRepository,
	passwordRecoveryRepository lib.PasswordRecoveryRepository,
	mailer lib.Sender,
) {
	PasswordRecoverer = lib.NewPasswordRecoverer(
		logger,
		passwordHasher,
		userRepository,
		passwordRecoveryRepository,
		mailer,
	)
}
