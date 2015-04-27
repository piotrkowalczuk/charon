package service

import (
	"errors"

	"github.com/go-soa/charon/lib/security"
)

var (
	// PasswordHasher ...
	PasswordHasher security.PasswordHasher
	// ErrPasswordHasherStrategyNotSupported ...
	ErrPasswordHasherStrategyNotSupported = errors.New("service: password hasher strategy not supported")
	// ErrPasswordHasherStrategyOptionsMissing ...
	ErrPasswordHasherStrategyOptionsMissing = errors.New("service: password hasher strategy options missing")
)

// PasswordHasherConfig ...
type PasswordHasherConfig struct {
	StrategyName       string                      `xml:"strategy"`
	BcryptStrategyOpts *BcryptPasswordHasherConfig `xml:"bcrypt-options"`
}

// BcryptPasswordHasherConfig ...
type BcryptPasswordHasherConfig struct {
	Cost int `xml:"cost"`
}

// InitPasswordHasher ...
func InitPasswordHasher(config PasswordHasherConfig) {
	var ph security.PasswordHasher

	switch config.StrategyName {
	case "bcrypt":
		if config.BcryptStrategyOpts == nil {
			Logger.Fatal(ErrPasswordHasherStrategyOptionsMissing)
		}

		bh, err := security.NewBcryptPasswordHasher(config.BcryptStrategyOpts.Cost, Logger)
		if err != nil {
			Logger.Fatal(err)
		}

		ph = security.PasswordHasher(bh)
	default:
		Logger.Fatal(ErrPasswordHasherStrategyNotSupported)
	}

	PasswordHasher = ph
}
