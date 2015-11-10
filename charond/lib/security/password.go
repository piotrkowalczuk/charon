package security

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrBcryptHasherCostOutOfRange ...
	ErrBcryptHasherCostOutOfRange = errors.New("security: bcrypt cost out of range")
)

// PasswordHasher ...
type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) bool
}

// BcryptPasswordHasher ...
type BcryptPasswordHasher struct {
	logger *logrus.Logger
	cost   int
}

// NewBcryptPasswordHasher ...
func NewBcryptPasswordHasher(cost int, logger *logrus.Logger) (*BcryptPasswordHasher, error) {
	if bcrypt.MinCost > cost || cost > bcrypt.MaxCost {
		return nil, ErrBcryptHasherCostOutOfRange
	}

	return &BcryptPasswordHasher{
		cost:   cost,
		logger: logger,
	}, nil
}

// Hash implements PasswordHasher interface.
func (bph BcryptPasswordHasher) Hash(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bph.cost)

	return string(hashedPassword), err
}

// Compare implements PasswordHasher interface.
func (bph BcryptPasswordHasher) Compare(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if bph.logger != nil {
			bph.logger.Debug(err)
		}
		return false
	}

	return true
}
