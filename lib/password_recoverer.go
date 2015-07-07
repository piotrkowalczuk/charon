package lib

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/go-soa/charon/lib/security"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrPasswordRecovererUserIsNotActive    = errors.New("lib: password recoverer user is not active")
	ErrPasswordRecovererUserIsNotConfirmed = errors.New("lib: password recoverer user is not confirmed")
)

type PasswordRecoverer interface {
	Start(string) error
	Finalize(int64, string, string) error
}

type passwordRecoverer struct {
	userRepository             UserRepository
	passwordRecoveryRepository PasswordRecoveryRepository
	passwordHasher             security.PasswordHasher
	mailer                     Sender
	logger                     logrus.StdLogger
}

// NewPasswordRecoverer ...
func NewPasswordRecoverer(
	logger logrus.StdLogger,
	passwordHasher security.PasswordHasher,
	userRepository UserRepository,
	passwordRecoveryRepository PasswordRecoveryRepository,
	mailer Sender,
) *passwordRecoverer {
	return &passwordRecoverer{
		logger:                     logger,
		userRepository:             userRepository,
		passwordRecoveryRepository: passwordRecoveryRepository,
		passwordHasher:             passwordHasher,
		mailer:                     mailer,
	}
}

// StartRecovery ...
func (pr *passwordRecoverer) Start(email string) error {
	user, err := pr.userRepository.FindOneByUsername(email)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return ErrPasswordRecovererUserIsNotActive
	}

	if !user.IsConfirmed {
		return ErrPasswordRecovererUserIsNotConfirmed
	}

	passwordRecovery, err := pr.passwordRecoveryRepository.Create(user.ID, uuid.NewV4().String())
	if err != nil {
		return err
	}

	err = pr.mailer.Send(user.Username, map[string]interface{}{
		"user_id":            passwordRecovery.UserID,
		"confirmation_token": passwordRecovery.ConfirmationToken,
	})
	if err != nil {
		return err
	}

	return nil
}

// Finalize ...
func (pr *passwordRecoverer) Finalize(userID int64, confirmationToken, plainPassword string) error {
	user, err := pr.userRepository.FindOneByID(userID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return ErrPasswordRecovererUserIsNotActive
	}

	if !user.IsConfirmed {
		return ErrPasswordRecovererUserIsNotConfirmed
	}

	passwordRecovery, err := pr.passwordRecoveryRepository.FindOneInProgress(userID, confirmationToken)
	if err != nil {
		return err
	}

	hashedPassword, err := pr.passwordHasher.Hash(plainPassword)
	if err != nil {
		return err
	}

	if err := pr.userRepository.ChangePassword(passwordRecovery.UserID, hashedPassword); err != nil {
		return err
	}

	_, err = pr.passwordRecoveryRepository.Archive(userID, confirmationToken)
	if err != nil {
		return err
	}

	// TODO: implement email notification

	return nil
}
