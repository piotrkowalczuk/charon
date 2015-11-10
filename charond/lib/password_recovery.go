package lib

import "time"

const (
	PasswordRecoveryStatusNew = iota
	PasswordRecoveryStatusRecovered
	PasswordRecoveryStatusAbandoned
)

type PasswordRecovery struct {
	ID                int64      `json:"-"`
	UserID            int64      `json:"-"`
	Status            int64      `json:"-"`
	ConfirmationToken string     `json:"-"`
	CreatedAt         *time.Time `json:"-"`
	RecoveredAt       *time.Time `json:"-"`
}
