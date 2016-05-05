package charon

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/pqcomp"
	"github.com/piotrkowalczuk/pqt"
	"github.com/piotrkowalczuk/qtypes"
)

const (
	tableUser                              = "charon.user"
	tableUserColumnConfirmationToken       = "confirmation_token"
	tableUserColumnCreatedAt               = "created_at"
	tableUserColumnCreatedBy               = "created_by"
	tableUserColumnFirstName               = "first_name"
	tableUserColumnID                      = "id"
	tableUserColumnIsActive                = "is_active"
	tableUserColumnIsConfirmed             = "is_confirmed"
	tableUserColumnIsStaff                 = "is_staff"
	tableUserColumnIsSuperuser             = "is_superuser"
	tableUserColumnLastLoginAt             = "last_login_at"
	tableUserColumnLastName                = "last_name"
	tableUserColumnPassword                = "password"
	tableUserColumnUpdatedAt               = "updated_at"
	tableUserColumnUpdatedBy               = "updated_by"
	tableUserColumnUsername                = "username"
	tableUserConstraintCreatedByForeignKey = "charon.user_created_by_fkey"
	tableUserConstraintPrimaryKey          = "charon.user_id_pkey"
	tableUserConstraintUpdatedByForeignKey = "charon.user_updated_by_fkey"
	tableUserConstraintUsernameUnique      = "charon.user_username_key"
)

var (
	tableUserColumns = []string{
		tableUserColumnConfirmationToken,
		tableUserColumnCreatedAt,
		tableUserColumnCreatedBy,
		tableUserColumnFirstName,
		tableUserColumnID,
		tableUserColumnIsActive,
		tableUserColumnIsConfirmed,
		tableUserColumnIsStaff,
		tableUserColumnIsSuperuser,
		tableUserColumnLastLoginAt,
		tableUserColumnLastName,
		tableUserColumnPassword,
		tableUserColumnUpdatedAt,
		tableUserColumnUpdatedBy,
		tableUserColumnUsername,
	}
)

type userEntity struct {
	ConfirmationToken []byte
	CreatedAt         time.Time
	CreatedBy         *ntypes.Int64
	FirstName         string
	ID                int64
	IsActive          bool
	IsConfirmed       bool
	IsStaff           bool
	IsSuperuser       bool
	LastLoginAt       *time.Time
	LastName          string
	Password          []byte
	UpdatedAt         *time.Time
	UpdatedBy         *ntypes.Int64
	Username          string
	Author            []*userEntity
	Modifier          []*userEntity
	Permission        []*permissionEntity
	Group             []*groupEntity
}
type userCriteria struct {
	offset, limit     int64
	sort              map[string]bool
	confirmationToken []byte

	createdAt *qtypes.Timestamp

	createdBy *qtypes.Int64

	firstName *qtypes.String

	id *qtypes.Int64

	isActive *ntypes.Bool

	isConfirmed *ntypes.Bool

	isStaff *ntypes.Bool

	isSuperuser *ntypes.Bool

	lastLoginAt *qtypes.Timestamp

	lastName *qtypes.String

	password []byte

	updatedAt *qtypes.Timestamp

	updatedBy *qtypes.Int64

	username *qtypes.String
}

type userRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userRepository) Find(c *userCriteria) ([]*userEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false
	if c.confirmationToken != nil {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnConfirmationToken)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.confirmationToken)
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tableUserColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.createdBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[0])

			wbuf.WriteString(tableUserColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[1])
		}
	}

	if c.firstName != nil && c.firstName.Valid {
		switch c.firstName.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			if c.firstName.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.firstName.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.firstName.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.firstName.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.firstName.Value()))
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" IN ")
			for _, v := range c.id.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[0])

			wbuf.WriteString(tableUserColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[1])
		}
	}

	if c.isActive != nil && c.isActive.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsActive)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.isActive)
	}
	if c.isConfirmed != nil && c.isConfirmed.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsConfirmed)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.isConfirmed)
	}
	if c.isStaff != nil && c.isStaff.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsStaff)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.isStaff)
	}
	if c.isSuperuser != nil && c.isSuperuser.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsSuperuser)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.isSuperuser)
	}

	if c.lastLoginAt != nil && c.lastLoginAt.Valid {
		lastLoginAtt1 := c.lastLoginAt.Value()
		if lastLoginAtt1 != nil {
			lastLoginAt1, err := ptypes.Timestamp(lastLoginAtt1)
			if err != nil {
				return nil, err
			}

			switch c.lastLoginAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				if c.lastLoginAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.lastLoginAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				lastLoginAtt2 := c.lastLoginAt.Values[1]
				if lastLoginAtt2 != nil {
					lastLoginAt2, err := ptypes.Timestamp(lastLoginAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserColumnLastLoginAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, lastLoginAt1)

					wbuf.WriteString(tableUserColumnLastLoginAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, lastLoginAt2)
				}
			}
		}
	}

	if c.lastName != nil && c.lastName.Valid {
		switch c.lastName.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			if c.lastName.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.lastName.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.lastName.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.lastName.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.lastName.Value()))
		}
	}

	if c.password != nil {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnPassword)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args = append(args, c.password)
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tableUserColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.updatedBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[0])

			wbuf.WriteString(tableUserColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[1])
		}
	}

	if c.username != nil && c.username.Valid {
		switch c.username.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			if c.username.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.username.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.username.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.username.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.username.Value()))
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*userEntity
	for rows.Next() {
		var entity userEntity
		err = rows.Scan(
			&entity.ConfirmationToken,
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.FirstName,
			&entity.ID,
			&entity.IsActive,
			&entity.IsConfirmed,
			&entity.IsStaff,
			&entity.IsSuperuser,
			&entity.LastLoginAt,
			&entity.LastName,
			&entity.Password,
			&entity.UpdatedAt,
			&entity.UpdatedBy,
			&entity.Username,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *userRepository) FindOneByID(id int64) (*userEntity, error) {
	var (
		query  string
		entity userEntity
	)
	query = `SELECT confirmation_token,
created_at,
created_by,
first_name,
id,
is_active,
is_confirmed,
is_staff,
is_superuser,
last_login_at,
last_name,
password,
updated_at,
updated_by,
username
 FROM charon.user WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&entity.ConfirmationToken,
		&entity.CreatedAt,
		&entity.CreatedBy,
		&entity.FirstName,
		&entity.ID,
		&entity.IsActive,
		&entity.IsConfirmed,
		&entity.IsStaff,
		&entity.IsSuperuser,
		&entity.LastLoginAt,
		&entity.LastName,
		&entity.Password,
		&entity.UpdatedAt,
		&entity.UpdatedBy,
		&entity.Username,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *userRepository) Insert(e *userEntity) (*userEntity, error) {
	insert := pqcomp.New(0, 15)
	insert.AddExpr(tableUserColumnConfirmationToken, "", e.ConfirmationToken)
	insert.AddExpr(tableUserColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tableUserColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableUserColumnFirstName, "", e.FirstName)
	insert.AddExpr(tableUserColumnIsActive, "", e.IsActive)
	insert.AddExpr(tableUserColumnIsConfirmed, "", e.IsConfirmed)
	insert.AddExpr(tableUserColumnIsStaff, "", e.IsStaff)
	insert.AddExpr(tableUserColumnIsSuperuser, "", e.IsSuperuser)
	insert.AddExpr(tableUserColumnLastLoginAt, "", e.LastLoginAt)
	insert.AddExpr(tableUserColumnLastName, "", e.LastName)
	insert.AddExpr(tableUserColumnPassword, "", e.Password)
	insert.AddExpr(tableUserColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableUserColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(tableUserColumnUsername, "", e.Username)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.ConfirmationToken,
		&e.CreatedAt,
		&e.CreatedBy,
		&e.FirstName,
		&e.ID,
		&e.IsActive,
		&e.IsConfirmed,
		&e.IsStaff,
		&e.IsSuperuser,
		&e.LastLoginAt,
		&e.LastName,
		&e.Password,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.Username,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userRepository) UpdateByID(
	id int64,
	confirmationToken []byte,
	createdAt *time.Time,
	createdBy *ntypes.Int64,
	firstName *ntypes.String,
	isActive *ntypes.Bool,
	isConfirmed *ntypes.Bool,
	isStaff *ntypes.Bool,
	isSuperuser *ntypes.Bool,
	lastLoginAt *time.Time,
	lastName *ntypes.String,
	password []byte,
	updatedAt *time.Time,
	updatedBy *ntypes.Int64,
	username *ntypes.String,
) (*userEntity, error) {
	update := pqcomp.New(0, 15)
	update.AddExpr(tableUserColumnID, pqcomp.Equal, id)
	update.AddExpr(tableUserColumnConfirmationToken, pqcomp.Equal, confirmationToken)
	if createdAt != nil {
		update.AddExpr(tableUserColumnCreatedAt, pqcomp.Equal, createdAt)

	}
	update.AddExpr(tableUserColumnCreatedBy, pqcomp.Equal, createdBy)
	update.AddExpr(tableUserColumnFirstName, pqcomp.Equal, firstName)

	update.AddExpr(tableUserColumnIsActive, pqcomp.Equal, isActive)

	update.AddExpr(tableUserColumnIsConfirmed, pqcomp.Equal, isConfirmed)

	update.AddExpr(tableUserColumnIsStaff, pqcomp.Equal, isStaff)

	update.AddExpr(tableUserColumnIsSuperuser, pqcomp.Equal, isSuperuser)
	update.AddExpr(tableUserColumnLastLoginAt, pqcomp.Equal, lastLoginAt)
	update.AddExpr(tableUserColumnLastName, pqcomp.Equal, lastName)
	update.AddExpr(tableUserColumnPassword, pqcomp.Equal, password)
	if updatedAt != nil {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.Equal, updatedAt)
	} else {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableUserColumnUpdatedBy, pqcomp.Equal, updatedBy)
	update.AddExpr(tableUserColumnUsername, pqcomp.Equal, username)

	if update.Len() == 0 {
		return nil, errors.New("charon: user update failure, nothing to update")
	}
	query := "UPDATE charon.user SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e userEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.ConfirmationToken,
		&e.CreatedAt,
		&e.CreatedBy,
		&e.FirstName,
		&e.ID,
		&e.IsActive,
		&e.IsConfirmed,
		&e.IsStaff,
		&e.IsSuperuser,
		&e.LastLoginAt,
		&e.LastName,
		&e.Password,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.Username,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *userRepository) DeleteByID(id int64) (int64, error) {
	query := "DELETE FROM charon.user WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	tableGroup                              = "charon.group"
	tableGroupColumnCreatedAt               = "created_at"
	tableGroupColumnCreatedBy               = "created_by"
	tableGroupColumnDescription             = "description"
	tableGroupColumnID                      = "id"
	tableGroupColumnName                    = "name"
	tableGroupColumnUpdatedAt               = "updated_at"
	tableGroupColumnUpdatedBy               = "updated_by"
	tableGroupConstraintCreatedByForeignKey = "charon.group_created_by_fkey"
	tableGroupConstraintPrimaryKey          = "charon.group_id_pkey"
	tableGroupConstraintNameUnique          = "charon.group_name_key"
	tableGroupConstraintUpdatedByForeignKey = "charon.group_updated_by_fkey"
)

var (
	tableGroupColumns = []string{
		tableGroupColumnCreatedAt,
		tableGroupColumnCreatedBy,
		tableGroupColumnDescription,
		tableGroupColumnID,
		tableGroupColumnName,
		tableGroupColumnUpdatedAt,
		tableGroupColumnUpdatedBy,
	}
)

type groupEntity struct {
	CreatedAt   time.Time
	CreatedBy   *ntypes.Int64
	Description *ntypes.String
	ID          int64
	Name        string
	UpdatedAt   *time.Time
	UpdatedBy   *ntypes.Int64
	Author      []*userEntity
	Modifier    []*userEntity
	Permission  []*permissionEntity
	Users       []*userEntity
}
type groupCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp

	createdBy *qtypes.Int64

	description *qtypes.String

	id *qtypes.Int64

	name *qtypes.String

	updatedAt *qtypes.Timestamp

	updatedBy *qtypes.Int64
}

type groupRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *groupRepository) Find(c *groupCriteria) ([]*groupEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableGroupColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tableGroupColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.createdBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[0])

			wbuf.WriteString(tableGroupColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[1])
		}
	}

	if c.description != nil && c.description.Valid {
		switch c.description.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			if c.description.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.description.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.description.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.description.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.description.Value()))
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" IN ")
			for _, v := range c.id.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[0])

			wbuf.WriteString(tableGroupColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[1])
		}
	}

	if c.name != nil && c.name.Valid {
		switch c.name.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			if c.name.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.name.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.name.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.name.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.name.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableGroupColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tableGroupColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.updatedBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[0])

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[1])
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*groupEntity
	for rows.Next() {
		var entity groupEntity
		err = rows.Scan(
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.Description,
			&entity.ID,
			&entity.Name,
			&entity.UpdatedAt,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *groupRepository) FindOneByID(id int64) (*groupEntity, error) {
	var (
		query  string
		entity groupEntity
	)
	query = `SELECT created_at,
created_by,
description,
id,
name,
updated_at,
updated_by
 FROM charon.group WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&entity.CreatedAt,
		&entity.CreatedBy,
		&entity.Description,
		&entity.ID,
		&entity.Name,
		&entity.UpdatedAt,
		&entity.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *groupRepository) Insert(e *groupEntity) (*groupEntity, error) {
	insert := pqcomp.New(0, 7)
	insert.AddExpr(tableGroupColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tableGroupColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableGroupColumnDescription, "", e.Description)
	insert.AddExpr(tableGroupColumnName, "", e.Name)
	insert.AddExpr(tableGroupColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableGroupColumnUpdatedBy, "", e.UpdatedBy)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *groupRepository) UpdateByID(
	id int64,
	createdAt *time.Time,
	createdBy *ntypes.Int64,
	description *ntypes.String,
	name *ntypes.String,
	updatedAt *time.Time,
	updatedBy *ntypes.Int64,
) (*groupEntity, error) {
	update := pqcomp.New(0, 7)
	update.AddExpr(tableGroupColumnID, pqcomp.Equal, id)
	if createdAt != nil {
		update.AddExpr(tableGroupColumnCreatedAt, pqcomp.Equal, createdAt)

	}
	update.AddExpr(tableGroupColumnCreatedBy, pqcomp.Equal, createdBy)
	update.AddExpr(tableGroupColumnDescription, pqcomp.Equal, description)
	update.AddExpr(tableGroupColumnName, pqcomp.Equal, name)
	if updatedAt != nil {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.Equal, updatedAt)
	} else {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableGroupColumnUpdatedBy, pqcomp.Equal, updatedBy)

	if update.Len() == 0 {
		return nil, errors.New("charon: group update failure, nothing to update")
	}
	query := "UPDATE charon.group SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e groupEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *groupRepository) DeleteByID(id int64) (int64, error) {
	query := "DELETE FROM charon.group WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	tablePermission                                      = "charon.permission"
	tablePermissionColumnAction                          = "action"
	tablePermissionColumnCreatedAt                       = "created_at"
	tablePermissionColumnID                              = "id"
	tablePermissionColumnModule                          = "module"
	tablePermissionColumnSubsystem                       = "subsystem"
	tablePermissionColumnUpdatedAt                       = "updated_at"
	tablePermissionConstraintPrimaryKey                  = "charon.permission_id_pkey"
	tablePermissionConstraintSubsystemModuleActionUnique = "charon.permission_subsystem_module_action_key"
)

var (
	tablePermissionColumns = []string{
		tablePermissionColumnAction,
		tablePermissionColumnCreatedAt,
		tablePermissionColumnID,
		tablePermissionColumnModule,
		tablePermissionColumnSubsystem,
		tablePermissionColumnUpdatedAt,
	}
)

type permissionEntity struct {
	Action    string
	CreatedAt time.Time
	ID        int64
	Module    string
	Subsystem string
	UpdatedAt *time.Time
	Groups    []*groupEntity
	Users     []*userEntity
}
type permissionCriteria struct {
	offset, limit int64
	sort          map[string]bool
	action        *qtypes.String

	createdAt *qtypes.Timestamp

	id *qtypes.Int64

	module *qtypes.String

	subsystem *qtypes.String

	updatedAt *qtypes.Timestamp
}

type permissionRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *permissionRepository) Find(c *permissionCriteria) ([]*permissionEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false

	if c.action != nil && c.action.Valid {
		switch c.action.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			if c.action.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.action.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.action.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.action.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.action.Value()))
		}
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tablePermissionColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tablePermissionColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.id)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" IN ")
			for _, v := range c.id.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[0])

			wbuf.WriteString(tablePermissionColumnID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.id.Values[1])
		}
	}

	if c.module != nil && c.module.Valid {
		switch c.module.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			if c.module.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.module.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.module.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.module.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.module.Value()))
		}
	}

	if c.subsystem != nil && c.subsystem.Valid {
		switch c.subsystem.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			if c.subsystem.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.subsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.subsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.subsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.subsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tablePermissionColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tablePermissionColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*permissionEntity
	for rows.Next() {
		var entity permissionEntity
		err = rows.Scan(
			&entity.Action,
			&entity.CreatedAt,
			&entity.ID,
			&entity.Module,
			&entity.Subsystem,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *permissionRepository) FindOneByID(id int64) (*permissionEntity, error) {
	var (
		query  string
		entity permissionEntity
	)
	query = `SELECT action,
created_at,
id,
module,
subsystem,
updated_at
 FROM charon.permission WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&entity.Action,
		&entity.CreatedAt,
		&entity.ID,
		&entity.Module,
		&entity.Subsystem,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *permissionRepository) Insert(e *permissionEntity) (*permissionEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tablePermissionColumnAction, "", e.Action)
	insert.AddExpr(tablePermissionColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tablePermissionColumnModule, "", e.Module)
	insert.AddExpr(tablePermissionColumnSubsystem, "", e.Subsystem)
	insert.AddExpr(tablePermissionColumnUpdatedAt, "", e.UpdatedAt)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *permissionRepository) UpdateByID(
	id int64,
	action *ntypes.String,
	createdAt *time.Time,
	module *ntypes.String,
	subsystem *ntypes.String,
	updatedAt *time.Time,
) (*permissionEntity, error) {
	update := pqcomp.New(0, 6)
	update.AddExpr(tablePermissionColumnID, pqcomp.Equal, id)
	update.AddExpr(tablePermissionColumnAction, pqcomp.Equal, action)
	if createdAt != nil {
		update.AddExpr(tablePermissionColumnCreatedAt, pqcomp.Equal, createdAt)

	}
	update.AddExpr(tablePermissionColumnModule, pqcomp.Equal, module)
	update.AddExpr(tablePermissionColumnSubsystem, pqcomp.Equal, subsystem)
	if updatedAt != nil {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.Equal, updatedAt)
	} else {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}

	if update.Len() == 0 {
		return nil, errors.New("charon: permission update failure, nothing to update")
	}
	query := "UPDATE charon.permission SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e permissionEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *permissionRepository) DeleteByID(id int64) (int64, error) {
	query := "DELETE FROM charon.permission WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	tableUserGroups                              = "charon.user_groups"
	tableUserGroupsColumnCreatedAt               = "created_at"
	tableUserGroupsColumnCreatedBy               = "created_by"
	tableUserGroupsColumnGroupID                 = "group_id"
	tableUserGroupsColumnUpdatedAt               = "updated_at"
	tableUserGroupsColumnUpdatedBy               = "updated_by"
	tableUserGroupsColumnUserID                  = "user_id"
	tableUserGroupsConstraintCreatedByForeignKey = "charon.user_groups_created_by_fkey"
	tableUserGroupsConstraintUpdatedByForeignKey = "charon.user_groups_updated_by_fkey"
	tableUserGroupsConstraintUserIDForeignKey    = "charon.user_groups_user_id_fkey"
	tableUserGroupsConstraintGroupIDForeignKey   = "charon.user_groups_group_id_fkey"
	tableUserGroupsConstraintUserIDGroupIDUnique = "charon.user_groups_user_id_group_id_key"
)

var (
	tableUserGroupsColumns = []string{
		tableUserGroupsColumnCreatedAt,
		tableUserGroupsColumnCreatedBy,
		tableUserGroupsColumnGroupID,
		tableUserGroupsColumnUpdatedAt,
		tableUserGroupsColumnUpdatedBy,
		tableUserGroupsColumnUserID,
	}
)

type userGroupsEntity struct {
	CreatedAt time.Time
	CreatedBy *ntypes.Int64
	GroupID   int64
	UpdatedAt *time.Time
	UpdatedBy *ntypes.Int64
	UserID    int64
	User      *userEntity
	Group     *groupEntity
	Author    []*userEntity
	Modifier  []*userEntity
}
type userGroupsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp

	createdBy *qtypes.Int64

	groupID *qtypes.Int64

	updatedAt *qtypes.Timestamp

	updatedBy *qtypes.Int64

	userID *qtypes.Int64
}

type userGroupsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userGroupsRepository) Find(c *userGroupsCriteria) ([]*userGroupsEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserGroupsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tableUserGroupsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.createdBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[0])

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" IN ")
			for _, v := range c.groupID.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID.Values[0])

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID.Values[1])
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.updatedBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[0])

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" IN ")
			for _, v := range c.userID.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID.Values[0])

			wbuf.WriteString(tableUserGroupsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID.Values[1])
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*userGroupsEntity
	for rows.Next() {
		var entity userGroupsEntity
		err = rows.Scan(
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.GroupID,
			&entity.UpdatedAt,
			&entity.UpdatedBy,
			&entity.UserID,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *userGroupsRepository) Insert(e *userGroupsEntity) (*userGroupsEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tableUserGroupsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tableUserGroupsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableUserGroupsColumnGroupID, "", e.GroupID)
	insert.AddExpr(tableUserGroupsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableUserGroupsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(tableUserGroupsColumnUserID, "", e.UserID)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

const (
	tableGroupPermissions                                                                           = "charon.group_permissions"
	tableGroupPermissionsColumnCreatedAt                                                            = "created_at"
	tableGroupPermissionsColumnCreatedBy                                                            = "created_by"
	tableGroupPermissionsColumnGroupID                                                              = "group_id"
	tableGroupPermissionsColumnPermissionAction                                                     = "permission_action"
	tableGroupPermissionsColumnPermissionModule                                                     = "permission_module"
	tableGroupPermissionsColumnPermissionSubsystem                                                  = "permission_subsystem"
	tableGroupPermissionsColumnUpdatedAt                                                            = "updated_at"
	tableGroupPermissionsColumnUpdatedBy                                                            = "updated_by"
	tableGroupPermissionsConstraintCreatedByForeignKey                                              = "charon.group_permissions_created_by_fkey"
	tableGroupPermissionsConstraintUpdatedByForeignKey                                              = "charon.group_permissions_updated_by_fkey"
	tableGroupPermissionsConstraintGroupIDForeignKey                                                = "charon.group_permissions_group_id_fkey"
	tableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey    = "charon.group_permissions_subsystem_module_action_fkey"
	tableGroupPermissionsConstraintGroupIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.group_permissions_group_id_subsystem_module_action_key"
)

var (
	tableGroupPermissionsColumns = []string{
		tableGroupPermissionsColumnCreatedAt,
		tableGroupPermissionsColumnCreatedBy,
		tableGroupPermissionsColumnGroupID,
		tableGroupPermissionsColumnPermissionAction,
		tableGroupPermissionsColumnPermissionModule,
		tableGroupPermissionsColumnPermissionSubsystem,
		tableGroupPermissionsColumnUpdatedAt,
		tableGroupPermissionsColumnUpdatedBy,
	}
)

type groupPermissionsEntity struct {
	CreatedAt           time.Time
	CreatedBy           *ntypes.Int64
	GroupID             int64
	PermissionAction    string
	PermissionModule    string
	PermissionSubsystem string
	UpdatedAt           *time.Time
	UpdatedBy           *ntypes.Int64
	Group               *groupEntity
	Permission          *permissionEntity
	Author              []*userEntity
	Modifier            []*userEntity
}
type groupPermissionsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp

	createdBy *qtypes.Int64

	groupID *qtypes.Int64

	permissionAction *qtypes.String

	permissionModule *qtypes.String

	permissionSubsystem *qtypes.String

	updatedAt *qtypes.Timestamp

	updatedBy *qtypes.Int64
}

type groupPermissionsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *groupPermissionsRepository) Find(c *groupPermissionsCriteria) ([]*groupPermissionsEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.createdBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[0])

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" IN ")
			for _, v := range c.groupID.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID.Values[0])

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.groupID.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionAction.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionAction.Value()))
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionModule.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionModule.Value()))
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionSubsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionSubsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.updatedBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[0])

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[1])
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*groupPermissionsEntity
	for rows.Next() {
		var entity groupPermissionsEntity
		err = rows.Scan(
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.GroupID,
			&entity.PermissionAction,
			&entity.PermissionModule,
			&entity.PermissionSubsystem,
			&entity.UpdatedAt,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *groupPermissionsRepository) Insert(e *groupPermissionsEntity) (*groupPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(tableGroupPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tableGroupPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableGroupPermissionsColumnGroupID, "", e.GroupID)
	insert.AddExpr(tableGroupPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(tableGroupPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedBy, "", e.UpdatedBy)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

const (
	tableUserPermissions                                                                          = "charon.user_permissions"
	tableUserPermissionsColumnCreatedAt                                                           = "created_at"
	tableUserPermissionsColumnCreatedBy                                                           = "created_by"
	tableUserPermissionsColumnPermissionAction                                                    = "permission_action"
	tableUserPermissionsColumnPermissionModule                                                    = "permission_module"
	tableUserPermissionsColumnPermissionSubsystem                                                 = "permission_subsystem"
	tableUserPermissionsColumnUpdatedAt                                                           = "updated_at"
	tableUserPermissionsColumnUpdatedBy                                                           = "updated_by"
	tableUserPermissionsColumnUserID                                                              = "user_id"
	tableUserPermissionsConstraintCreatedByForeignKey                                             = "charon.user_permissions_created_by_fkey"
	tableUserPermissionsConstraintUpdatedByForeignKey                                             = "charon.user_permissions_updated_by_fkey"
	tableUserPermissionsConstraintUserIDForeignKey                                                = "charon.user_permissions_user_id_fkey"
	tableUserPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey   = "charon.user_permissions_subsystem_module_action_fkey"
	tableUserPermissionsConstraintUserIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.user_permissions_user_id_subsystem_module_action_key"
)

var (
	tableUserPermissionsColumns = []string{
		tableUserPermissionsColumnCreatedAt,
		tableUserPermissionsColumnCreatedBy,
		tableUserPermissionsColumnPermissionAction,
		tableUserPermissionsColumnPermissionModule,
		tableUserPermissionsColumnPermissionSubsystem,
		tableUserPermissionsColumnUpdatedAt,
		tableUserPermissionsColumnUpdatedBy,
		tableUserPermissionsColumnUserID,
	}
)

type userPermissionsEntity struct {
	CreatedAt           time.Time
	CreatedBy           *ntypes.Int64
	PermissionAction    string
	PermissionModule    string
	PermissionSubsystem string
	UpdatedAt           *time.Time
	UpdatedBy           *ntypes.Int64
	UserID              int64
	User                *userEntity
	Permission          *permissionEntity
	Author              []*userEntity
	Modifier            []*userEntity
}
type userPermissionsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp

	createdBy *qtypes.Int64

	permissionAction *qtypes.String

	permissionModule *qtypes.String

	permissionSubsystem *qtypes.String

	updatedAt *qtypes.Timestamp

	updatedBy *qtypes.Int64

	userID *qtypes.Int64
}

type userPermissionsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userPermissionsRepository) Find(c *userPermissionsCriteria) ([]*userPermissionsEntity, error) {
	wbuf := bytes.NewBuffer(nil)
	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqt.NewPlaceholderWriter()
	args := make([]interface{}, 0)
	dirty := false

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.createdAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt1)

					wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.createdBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[0])

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.createdBy.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionAction.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionAction.Value()))
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionModule.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionModule.Value()))
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case qtypes.TextQueryType_NOT_A_TEXT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.TextQueryType_EXACT:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.permissionSubsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			wbuf.WriteString(" LIKE ")
			pw.WriteTo(wbuf)
			args = append(args, fmt.Sprintf("%%%s", c.permissionSubsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case qtypes.NumericQueryType_NOT_A_NUMBER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					wbuf.WriteString(" IS NOT NULL ")
				} else {
					wbuf.WriteString(" IS NULL ")
				}
			case qtypes.NumericQueryType_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_NOT_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("!=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_IN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString(" IN ")
				pw.WriteTo(wbuf)
				args = append(args, c.updatedAt)
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt1)

					wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args = append(args, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" IN ")
			for _, v := range c.updatedBy.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[0])

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case qtypes.NumericQueryType_NOT_A_NUMBER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" IS NOT NULL ")
			} else {
				wbuf.WriteString(" IS NULL ")
			}
		case qtypes.NumericQueryType_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString("=")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_NOT_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" <> ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" < ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" >= ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID)
		case qtypes.NumericQueryType_IN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" IN ")
			for _, v := range c.userID.Values {
				pw.WriteTo(wbuf)
				args = append(args, v)
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID.Values[0])

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			wbuf.WriteString(" > ")
			pw.WriteTo(wbuf)
			args = append(args, c.userID.Values[1])
		}
	}

	fmt.Println("is dirty", dirty)
	if dirty {
		if _, err := qbuf.WriteString(" WHERE "); err != nil {
			return nil, err
		}
		if _, err := wbuf.WriteTo(qbuf); err != nil {
			return nil, err
		}
	}

	qbuf.WriteString(" OFFSET ")
	pw.WriteTo(qbuf)
	args = append(args, c.offset)
	qbuf.WriteString(" LIMIT ")
	pw.WriteTo(qbuf)
	args = append(args, c.limit)

	rows, err := r.db.Query(qbuf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*userPermissionsEntity
	for rows.Next() {
		var entity userPermissionsEntity
		err = rows.Scan(
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.PermissionAction,
			&entity.PermissionModule,
			&entity.PermissionSubsystem,
			&entity.UpdatedAt,
			&entity.UpdatedBy,
			&entity.UserID,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}
func (r *userPermissionsRepository) Insert(e *userPermissionsEntity) (*userPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(tableUserPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(tableUserPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableUserPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(tableUserPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(tableUserPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(tableUserPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableUserPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(tableUserPermissionsColumnUserID, "", e.UserID)

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() != 0 {
		b.WriteString(" (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(") VALUES (")
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(")")
		if len(r.columns) > 0 {
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(r.columns, ","))
		}
	}

	err := r.db.QueryRow(b.String(), insert.Args()...).Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

const schemaSQL = `
-- do not modify, generated by pqt

CREATE SCHEMA IF NOT EXISTS charon; 

CREATE TABLE IF NOT EXISTS charon.user (
	confirmation_token BYTEA,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	first_name TEXT NOT NULL,
	id BIGSERIAL,
	is_active BOOL DEFAULT FALSE NOT NULL,
	is_confirmed BOOL DEFAULT FALSE NOT NULL,
	is_staff BOOL DEFAULT FALSE NOT NULL,
	is_superuser BOOL DEFAULT FALSE NOT NULL,
	last_login_at TIMESTAMPTZ,
	last_name TEXT NOT NULL,
	password BYTEA NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	username TEXT NOT NULL,

	CONSTRAINT "charon.user_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.user_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_username_key" UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS charon.group (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	description TEXT,
	id BIGSERIAL,
	name TEXT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,

	CONSTRAINT "charon.group_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.group_name_key" UNIQUE (name),
	CONSTRAINT "charon.group_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.permission (
	action TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	module TEXT NOT NULL,
	subsystem TEXT NOT NULL,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "charon.permission_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.permission_subsystem_module_action_key" UNIQUE (subsystem, module, action)
);

CREATE TABLE IF NOT EXISTS charon.user_groups (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	group_id BIGINT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	user_id BIGINT NOT NULL,

	CONSTRAINT "charon.user_groups_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_groups_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_groups_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_groups_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
	CONSTRAINT "charon.user_groups_user_id_group_id_key" UNIQUE (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS charon.group_permissions (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	group_id BIGINT NOT NULL,
	permission_action TEXT NOT NULL,
	permission_module TEXT NOT NULL,
	permission_subsystem TEXT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,

	CONSTRAINT "charon.group_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_permissions_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
	CONSTRAINT "charon.group_permissions_subsystem_module_action_fkey" FOREIGN KEY (permission_subsystem, permission_module, permission_action) REFERENCES charon.permission (subsystem, module, action),
	CONSTRAINT "charon.group_permissions_group_id_subsystem_module_action_key" UNIQUE (group_id, permission_subsystem, permission_module, permission_action)
);

CREATE TABLE IF NOT EXISTS charon.user_permissions (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	permission_action TEXT NOT NULL,
	permission_module TEXT NOT NULL,
	permission_subsystem TEXT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	user_id BIGINT NOT NULL,

	CONSTRAINT "charon.user_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_subsystem_module_action_fkey" FOREIGN KEY (permission_subsystem, permission_module, permission_action) REFERENCES charon.permission (subsystem, module, action),
	CONSTRAINT "charon.user_permissions_user_id_subsystem_module_action_key" UNIQUE (user_id, permission_subsystem, permission_module, permission_action)
);

`
