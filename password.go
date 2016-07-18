package charon

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrBCryptCostOutOfRange can be returned by NewBCryptPasswordHasher if provided cost is not between min and max.
	ErrBCryptCostOutOfRange = errors.New("charon: bcrypt cost out of range")
)

// PasswordHasher define set of methods that object needs to implement to be considered as a hasher.
type PasswordHasher interface {
	Hash([]byte) ([]byte, error)
	Compare([]byte, []byte) bool
}

// BCryptPasswordHasher hasher that use BCrypt algorithm to secure password.
type BCryptPasswordHasher struct {
	cost int
}

// NewBCryptPasswordHasher allocates new BCryptPasswordHasher.
// If cost is not between min and max value it returns an error.
func NewBCryptPasswordHasher(cost int) (PasswordHasher, error) {
	if bcrypt.MinCost > cost || cost > bcrypt.MaxCost {
		return nil, ErrBCryptCostOutOfRange
	}

	return &BCryptPasswordHasher{cost: cost}, nil
}

// Hash implements PasswordHasher interface.
func (bph BCryptPasswordHasher) Hash(plainPassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plainPassword, bph.cost)
}

// Compare implements PasswordHasher interface.
func (bph BCryptPasswordHasher) Compare(hashedPassword, plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	return err == nil
}
