package charon

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrBcryptCostOutOfRange ...
	ErrBcryptCostOutOfRange = errors.New("charon: bcrypt cost out of range")
)

type Hasher interface {
	Hash([]byte) ([]byte, error)
}

type Comparator interface {
	Compare([]byte, []byte) bool
}

// PasswordHasher ...
type PasswordHasher interface {
	Hasher
	Comparator
}

// BcryptPasswordHasher ...
type BcryptPasswordHasher struct {
	logger log.Logger
	cost   int
}

// NewBcryptPasswordHasher ...
func NewBcryptPasswordHasher(cost int, logger log.Logger) (PasswordHasher, error) {
	if bcrypt.MinCost > cost || cost > bcrypt.MaxCost {
		return nil, ErrBcryptCostOutOfRange
	}

	return &BcryptPasswordHasher{
		cost:   cost,
		logger: logger,
	}, nil
}

// Hash implements PasswordHasher interface.
func (bph BcryptPasswordHasher) Hash(plainPassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plainPassword, bph.cost)
}

// Compare implements PasswordHasher interface.
func (bph BcryptPasswordHasher) Compare(hashedPassword, plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	if err != nil {
		if bph.logger != nil {
			sklog.Error(bph.logger, err)
		}
		return false
	}

	return true
}
