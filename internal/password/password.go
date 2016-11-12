package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrBCryptCostOutOfRange can be returned by NewBCryptHasher if provided cost is not between min and max.
	ErrBCryptCostOutOfRange = errors.New("password: bcrypt cost out of range")
)

// Hasher define set of methods that object needs to implement to be considered as a hasher.
type Hasher interface {
	Hash([]byte) ([]byte, error)
	Compare([]byte, []byte) bool
}

// BCryptHasher hasher that use BCrypt algorithm to secure password.
type BCryptHasher struct {
	cost int
}

// NewBCryptHasher allocates new NewBCryptHasher.
// If cost is not between min and max value it returns an error.
func NewBCryptHasher(cost int) (Hasher, error) {
	if bcrypt.MinCost > cost || cost > bcrypt.MaxCost {
		return nil, ErrBCryptCostOutOfRange
	}

	return &BCryptHasher{cost: cost}, nil
}

// Hash implements Hasher interface.
func (bh BCryptHasher) Hash(plainPassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plainPassword, bh.cost)
}

// Compare implements Hasher interface.
func (bh BCryptHasher) Compare(hashedPassword, plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	return err == nil
}
