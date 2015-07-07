package lib

import (
	"database/sql"

	"github.com/go-soa/charon/model"
)

// PasswordRecoveryRepository ...
type PasswordRecoveryRepository interface {
	Create(int64, string) (*PasswordRecovery, error)
	Archive(int64, string) (int64, error)
	FindOneInProgress(int64, string) (*PasswordRecovery, error)
}

// passwordRecoveryRepository ...
type passwordRecoveryRepository struct {
	db *sql.DB
}

// NewPasswordRecoveryRepository ...
func NewPasswordRecoveryRepository(dbPool *sql.DB) (repository *passwordRecoveryRepository) {
	repository = &passwordRecoveryRepository{dbPool}

	return
}

// Create ...
func (prr *passwordRecoveryRepository) Create(userID int64, confirmationToken string) (*PasswordRecovery, error) {
	query := `
		INSERT INTO charon_password_recovery (
			user_id, confirmation_token, status, created_at
		)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`
	passwordRecovery := &PasswordRecovery{
		UserID:            userID,
		ConfirmationToken: confirmationToken,
		Status:            PasswordRecoveryStatusNew,
	}
	err := prr.db.QueryRow(
		query,
		passwordRecovery.UserID,
		passwordRecovery.ConfirmationToken,
		passwordRecovery.Status,
	).Scan(&passwordRecovery.ID)

	return passwordRecovery, mapKnownErrors(userKnownErrors, err)
}

// Create ...
func (prr *passwordRecoveryRepository) Archive(userID int64, confirmationToken string) (int64, error) {
	query := `
		UPDATE charon_password_recovery
		SET confirmation_token = $1, status = $2, recovered_at = NOW()
		WHERE confirmation_token = $3 AND user_id = $4
	`
	passwordRecovery := &PasswordRecovery{
		UserID:            userID,
		ConfirmationToken: confirmationToken,
		Status:            PasswordRecoveryStatusNew,
	}
	r, err := prr.db.Exec(
		query,
		model.UserConfirmationTokenUsed,
		PasswordRecoveryStatusRecovered,
		passwordRecovery.ConfirmationToken,
		passwordRecovery.UserID,
	)

	if err != nil {
		return 0, mapKnownErrors(userKnownErrors, err)
	}

	return r.RowsAffected()
}

// FindOneInProgress ...
func (prr *passwordRecoveryRepository) FindOneInProgress(userID int64, confirmationToken string) (*PasswordRecovery, error) {
	query := `
		SELECT pr.id, pr.user_id, pr.confirmation_token, pr.status, pr.created_at, pr.recovered_at
		FROM charon_password_recovery AS pr
		WHERE pr.user_id = $1
			AND pr.confirmation_token = $2
			AND pr.created_at > NOW() - INTERVAL '1 day'
			AND pr.recovered_at IS NULL
		ORDER BY pr.id DESC
		LIMIT 1
	`
	var passwordRecovery PasswordRecovery
	err := prr.db.QueryRow(
		query,
		userID,
		confirmationToken,
	).Scan(
		&passwordRecovery.ID,
		&passwordRecovery.UserID,
		&passwordRecovery.ConfirmationToken,
		&passwordRecovery.Status,
		&passwordRecovery.CreatedAt,
		&passwordRecovery.RecoveredAt,
	)

	return &passwordRecovery, mapKnownErrors(userKnownErrors, err)
}
