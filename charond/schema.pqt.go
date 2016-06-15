package charond

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/pqcomp"
	"github.com/piotrkowalczuk/pqt/pqtgo"
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
	Author            *userEntity
	Modifier          *userEntity
	Permission        []*permissionEntity
	Group             []*groupEntity
}

func (e *userEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserColumnConfirmationToken:
		return &e.ConfirmationToken, true
	case tableUserColumnCreatedAt:
		return &e.CreatedAt, true
	case tableUserColumnCreatedBy:
		return &e.CreatedBy, true
	case tableUserColumnFirstName:
		return &e.FirstName, true
	case tableUserColumnID:
		return &e.ID, true
	case tableUserColumnIsActive:
		return &e.IsActive, true
	case tableUserColumnIsConfirmed:
		return &e.IsConfirmed, true
	case tableUserColumnIsStaff:
		return &e.IsStaff, true
	case tableUserColumnIsSuperuser:
		return &e.IsSuperuser, true
	case tableUserColumnLastLoginAt:
		return &e.LastLoginAt, true
	case tableUserColumnLastName:
		return &e.LastName, true
	case tableUserColumnPassword:
		return &e.Password, true
	case tableUserColumnUpdatedAt:
		return &e.UpdatedAt, true
	case tableUserColumnUpdatedBy:
		return &e.UpdatedBy, true
	case tableUserColumnUsername:
		return &e.Username, true
	default:
		return nil, false
	}
}
func (e *userEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// userIterator is not thread safe.
type userIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *userIterator) Next() bool {
	return i.rows.Next()
}

func (i *userIterator) Close() error {
	return i.rows.Close()
}

func (i *userIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *userIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround user method that makes iterator more generic.
func (i *userIterator) Ent() (interface{}, error) {
	return i.User()
}

func (i *userIterator) User() (*userEntity, error) {
	var ent userEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type userCriteria struct {
	offset, limit     int64
	sort              map[string]bool
	confirmationToken []byte
	createdAt         *qtypes.Timestamp
	createdBy         *qtypes.Int64
	firstName         *qtypes.String
	id                *qtypes.Int64
	isActive          *ntypes.Bool
	isConfirmed       *ntypes.Bool
	isStaff           *ntypes.Bool
	isSuperuser       *ntypes.Bool
	lastLoginAt       *qtypes.Timestamp
	lastName          *qtypes.String
	password          []byte
	updatedAt         *qtypes.Timestamp
	updatedBy         *qtypes.Int64
	username          *qtypes.String
}

func (c *userCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)
	if c.confirmationToken != nil {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		if wrt, err = wbuf.WriteString(tableUserColumnConfirmationToken); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt, err = wbuf.WriteString("="); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = pw.WriteTo(wbuf); err != nil {
			return
		}
		wr += wrt64

		args.Add(c.confirmationToken)
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.createdBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.createdBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnCreatedBy)
				if c.createdBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.createdBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[1])
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
			if c.firstName.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.firstName.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			if c.firstName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.firstName.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			if c.firstName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.firstName.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnFirstName)
			if c.firstName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.firstName.Value()))
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
			if c.id.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.id.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnID)
				if c.id.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.id.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[1])
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
		args.Add(c.isActive)
	}
	if c.isConfirmed != nil && c.isConfirmed.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsConfirmed)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args.Add(c.isConfirmed)
	}
	if c.isStaff != nil && c.isStaff.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsStaff)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args.Add(c.isStaff)
	}
	if c.isSuperuser != nil && c.isSuperuser.Valid {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		wbuf.WriteString(tableUserColumnIsSuperuser)
		wbuf.WriteString("=")
		pw.WriteTo(wbuf)
		args.Add(c.isSuperuser)
	}

	if c.lastLoginAt != nil && c.lastLoginAt.Valid {
		lastLoginAtt1 := c.lastLoginAt.Value()
		if lastLoginAtt1 != nil {
			lastLoginAt1, err := ptypes.Timestamp(lastLoginAtt1)
			if err != nil {
				return wr, err
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
				if c.lastLoginAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.lastLoginAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.lastLoginAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.lastLoginAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.lastLoginAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnLastLoginAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.lastLoginAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.lastLoginAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserColumnLastLoginAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.lastLoginAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				lastLoginAtt2 := c.lastLoginAt.Values[1]
				if lastLoginAtt2 != nil {
					lastLoginAt2, err := ptypes.Timestamp(lastLoginAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserColumnLastLoginAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(lastLoginAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserColumnLastLoginAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(lastLoginAt2)
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
			if c.lastName.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.lastName.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			if c.lastName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.lastName.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			if c.lastName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.lastName.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnLastName)
			if c.lastName.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.lastName.Value()))
		}
	}

	if c.password != nil {
		if dirty {
			wbuf.WriteString(" AND ")
		}
		dirty = true
		if wrt, err = wbuf.WriteString(tableUserColumnPassword); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt, err = wbuf.WriteString("="); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = pw.WriteTo(wbuf); err != nil {
			return
		}
		wr += wrt64

		args.Add(c.password)
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
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
			if c.updatedBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.updatedBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserColumnUpdatedBy)
				if c.updatedBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.updatedBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[1])
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
			if c.username.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.username.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			if c.username.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.username.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			if c.username.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.username.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserColumnUsername)
			if c.username.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.username.Value()))
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type userPatch struct {
	id                int64
	confirmationToken []byte
	createdAt         *time.Time
	createdBy         *ntypes.Int64
	firstName         *ntypes.String
	isActive          *ntypes.Bool
	isConfirmed       *ntypes.Bool
	isStaff           *ntypes.Bool
	isSuperuser       *ntypes.Bool
	lastLoginAt       *time.Time
	lastName          *ntypes.String
	password          []byte
	updatedAt         *time.Time
	updatedBy         *ntypes.Int64
	username          *ntypes.String
}

type userRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserRows(rows *sql.Rows) ([]*userEntity, error) {
	var (
		entities []*userEntity
		err      error
	)
	for rows.Next() {
		var ent userEntity
		err = rows.Scan(
			&ent.ConfirmationToken,
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.FirstName,
			&ent.ID,
			&ent.IsActive,
			&ent.IsConfirmed,
			&ent.IsStaff,
			&ent.IsSuperuser,
			&ent.LastLoginAt,
			&ent.LastName,
			&ent.Password,
			&ent.UpdatedAt,
			&ent.UpdatedBy,
			&ent.Username,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *userRepositoryBase) Count(c *userCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *userRepositoryBase) Find(c *userCriteria) ([]*userEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserRows(rows)
}
func (r *userRepositoryBase) FindIter(c *userCriteria) (*userIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &userIterator{rows: rows}, nil
}
func (r *userRepositoryBase) FindOneByID(id int64) (*userEntity, error) {
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
func (r *userRepositoryBase) Insert(e *userEntity) (*userEntity, error) {
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
func (r *userRepositoryBase) UpdateByID(patch *userPatch) (*userEntity, error) {
	update := pqcomp.New(0, 15)
	update.AddExpr(tableUserColumnID, pqcomp.Equal, patch.id)
	update.AddExpr(tableUserColumnConfirmationToken, pqcomp.Equal, patch.confirmationToken)
	if patch.createdAt != nil {
		update.AddExpr(tableUserColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tableUserColumnCreatedBy, pqcomp.Equal, patch.createdBy)
	update.AddExpr(tableUserColumnFirstName, pqcomp.Equal, patch.firstName)
	update.AddExpr(tableUserColumnIsActive, pqcomp.Equal, patch.isActive)
	update.AddExpr(tableUserColumnIsConfirmed, pqcomp.Equal, patch.isConfirmed)
	update.AddExpr(tableUserColumnIsStaff, pqcomp.Equal, patch.isStaff)
	update.AddExpr(tableUserColumnIsSuperuser, pqcomp.Equal, patch.isSuperuser)
	update.AddExpr(tableUserColumnLastLoginAt, pqcomp.Equal, patch.lastLoginAt)
	update.AddExpr(tableUserColumnLastName, pqcomp.Equal, patch.lastName)
	update.AddExpr(tableUserColumnPassword, pqcomp.Equal, patch.password)
	if patch.updatedAt != nil {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableUserColumnUpdatedBy, pqcomp.Equal, patch.updatedBy)
	update.AddExpr(tableUserColumnUsername, pqcomp.Equal, patch.username)

	if update.Len() == 0 {
		return nil, errors.New("user update failure, nothing to update")
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
func (r *userRepositoryBase) DeleteByID(id int64) (int64, error) {
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
	Author      *userEntity
	Modifier    *userEntity
	Permission  []*permissionEntity
	Users       []*userEntity
}

func (e *groupEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tableGroupColumnCreatedAt:
		return &e.CreatedAt, true
	case tableGroupColumnCreatedBy:
		return &e.CreatedBy, true
	case tableGroupColumnDescription:
		return &e.Description, true
	case tableGroupColumnID:
		return &e.ID, true
	case tableGroupColumnName:
		return &e.Name, true
	case tableGroupColumnUpdatedAt:
		return &e.UpdatedAt, true
	case tableGroupColumnUpdatedBy:
		return &e.UpdatedBy, true
	default:
		return nil, false
	}
}
func (e *groupEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// groupIterator is not thread safe.
type groupIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *groupIterator) Next() bool {
	return i.rows.Next()
}

func (i *groupIterator) Close() error {
	return i.rows.Close()
}

func (i *groupIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *groupIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround group method that makes iterator more generic.
func (i *groupIterator) Ent() (interface{}, error) {
	return i.Group()
}

func (i *groupIterator) Group() (*groupEntity, error) {
	var ent groupEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type groupCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp
	createdBy     *qtypes.Int64
	description   *qtypes.String
	id            *qtypes.Int64
	name          *qtypes.String
	updatedAt     *qtypes.Timestamp
	updatedBy     *qtypes.Int64
}

func (c *groupCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableGroupColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableGroupColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableGroupColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.createdBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.createdBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnCreatedBy)
				if c.createdBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.createdBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[1])
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
			if c.description.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.description.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			if c.description.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.description.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			if c.description.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.description.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnDescription)
			if c.description.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.description.Value()))
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
			if c.id.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.id.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnID)
				if c.id.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.id.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[1])
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
			if c.name.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.name.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			if c.name.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.name.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			if c.name.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.name.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnName)
			if c.name.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.name.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableGroupColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableGroupColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableGroupColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
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
			if c.updatedBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.updatedBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupColumnUpdatedBy)
				if c.updatedBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.updatedBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[1])
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type groupPatch struct {
	id          int64
	createdAt   *time.Time
	createdBy   *ntypes.Int64
	description *ntypes.String
	name        *ntypes.String
	updatedAt   *time.Time
	updatedBy   *ntypes.Int64
}

type groupRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanGroupRows(rows *sql.Rows) ([]*groupEntity, error) {
	var (
		entities []*groupEntity
		err      error
	)
	for rows.Next() {
		var ent groupEntity
		err = rows.Scan(
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.Description,
			&ent.ID,
			&ent.Name,
			&ent.UpdatedAt,
			&ent.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *groupRepositoryBase) Count(c *groupCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *groupRepositoryBase) Find(c *groupCriteria) ([]*groupEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanGroupRows(rows)
}
func (r *groupRepositoryBase) FindIter(c *groupCriteria) (*groupIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &groupIterator{rows: rows}, nil
}
func (r *groupRepositoryBase) FindOneByID(id int64) (*groupEntity, error) {
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
func (r *groupRepositoryBase) Insert(e *groupEntity) (*groupEntity, error) {
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
func (r *groupRepositoryBase) UpdateByID(patch *groupPatch) (*groupEntity, error) {
	update := pqcomp.New(0, 7)
	update.AddExpr(tableGroupColumnID, pqcomp.Equal, patch.id)
	if patch.createdAt != nil {
		update.AddExpr(tableGroupColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tableGroupColumnCreatedBy, pqcomp.Equal, patch.createdBy)
	update.AddExpr(tableGroupColumnDescription, pqcomp.Equal, patch.description)
	update.AddExpr(tableGroupColumnName, pqcomp.Equal, patch.name)
	if patch.updatedAt != nil {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableGroupColumnUpdatedBy, pqcomp.Equal, patch.updatedBy)

	if update.Len() == 0 {
		return nil, errors.New("charond: group update failure, nothing to update")
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
func (r *groupRepositoryBase) DeleteByID(id int64) (int64, error) {
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

func (e *permissionEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tablePermissionColumnAction:
		return &e.Action, true
	case tablePermissionColumnCreatedAt:
		return &e.CreatedAt, true
	case tablePermissionColumnID:
		return &e.ID, true
	case tablePermissionColumnModule:
		return &e.Module, true
	case tablePermissionColumnSubsystem:
		return &e.Subsystem, true
	case tablePermissionColumnUpdatedAt:
		return &e.UpdatedAt, true
	default:
		return nil, false
	}
}
func (e *permissionEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// permissionIterator is not thread safe.
type permissionIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *permissionIterator) Next() bool {
	return i.rows.Next()
}

func (i *permissionIterator) Close() error {
	return i.rows.Close()
}

func (i *permissionIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *permissionIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround permission method that makes iterator more generic.
func (i *permissionIterator) Ent() (interface{}, error) {
	return i.Permission()
}

func (i *permissionIterator) Permission() (*permissionEntity, error) {
	var ent permissionEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type permissionCriteria struct {
	offset, limit int64
	sort          map[string]bool
	action        *qtypes.String
	createdAt     *qtypes.Timestamp
	id            *qtypes.Int64
	module        *qtypes.String
	subsystem     *qtypes.String
	updatedAt     *qtypes.Timestamp
}

func (c *permissionCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)

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
			if c.action.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.action.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			if c.action.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.action.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			if c.action.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.action.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnAction)
			if c.action.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.action.Value()))
		}
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tablePermissionColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tablePermissionColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tablePermissionColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.id.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.id.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnID)
				if c.id.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.id.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tablePermissionColumnID)
			if c.id.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.id.Values[1])
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
			if c.module.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.module.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			if c.module.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.module.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			if c.module.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.module.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnModule)
			if c.module.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.module.Value()))
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
			if c.subsystem.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.subsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			if c.subsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.subsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			if c.subsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.subsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tablePermissionColumnSubsystem)
			if c.subsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.subsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tablePermissionColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tablePermissionColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tablePermissionColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tablePermissionColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
				}
			}
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type permissionPatch struct {
	id        int64
	action    *ntypes.String
	createdAt *time.Time
	module    *ntypes.String
	subsystem *ntypes.String
	updatedAt *time.Time
}

type permissionRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanPermissionRows(rows *sql.Rows) ([]*permissionEntity, error) {
	var (
		entities []*permissionEntity
		err      error
	)
	for rows.Next() {
		var ent permissionEntity
		err = rows.Scan(
			&ent.Action,
			&ent.CreatedAt,
			&ent.ID,
			&ent.Module,
			&ent.Subsystem,
			&ent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *permissionRepositoryBase) Count(c *permissionCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *permissionRepositoryBase) Find(c *permissionCriteria) ([]*permissionEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanPermissionRows(rows)
}
func (r *permissionRepositoryBase) FindIter(c *permissionCriteria) (*permissionIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &permissionIterator{rows: rows}, nil
}
func (r *permissionRepositoryBase) FindOneByID(id int64) (*permissionEntity, error) {
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
func (r *permissionRepositoryBase) Insert(e *permissionEntity) (*permissionEntity, error) {
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
func (r *permissionRepositoryBase) UpdateByID(patch *permissionPatch) (*permissionEntity, error) {
	update := pqcomp.New(0, 6)
	update.AddExpr(tablePermissionColumnID, pqcomp.Equal, patch.id)
	update.AddExpr(tablePermissionColumnAction, pqcomp.Equal, patch.action)
	if patch.createdAt != nil {
		update.AddExpr(tablePermissionColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tablePermissionColumnModule, pqcomp.Equal, patch.module)
	update.AddExpr(tablePermissionColumnSubsystem, pqcomp.Equal, patch.subsystem)
	if patch.updatedAt != nil {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}

	if update.Len() == 0 {
		return nil, errors.New("permission update failure, nothing to update")
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
func (r *permissionRepositoryBase) DeleteByID(id int64) (int64, error) {
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
	Author    *userEntity
	Modifier  *userEntity
}

func (e *userGroupsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserGroupsColumnCreatedAt:
		return &e.CreatedAt, true
	case tableUserGroupsColumnCreatedBy:
		return &e.CreatedBy, true
	case tableUserGroupsColumnGroupID:
		return &e.GroupID, true
	case tableUserGroupsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case tableUserGroupsColumnUpdatedBy:
		return &e.UpdatedBy, true
	case tableUserGroupsColumnUserID:
		return &e.UserID, true
	default:
		return nil, false
	}
}
func (e *userGroupsEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// userGroupsIterator is not thread safe.
type userGroupsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *userGroupsIterator) Next() bool {
	return i.rows.Next()
}

func (i *userGroupsIterator) Close() error {
	return i.rows.Close()
}

func (i *userGroupsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *userGroupsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround userGroups method that makes iterator more generic.
func (i *userGroupsIterator) Ent() (interface{}, error) {
	return i.UserGroups()
}

func (i *userGroupsIterator) UserGroups() (*userGroupsEntity, error) {
	var ent userGroupsEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type userGroupsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *qtypes.Timestamp
	createdBy     *qtypes.Int64
	groupID       *qtypes.Int64
	updatedAt     *qtypes.Timestamp
	updatedBy     *qtypes.Int64
	userID        *qtypes.Int64
}

func (c *userGroupsCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserGroupsColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserGroupsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserGroupsColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.createdBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.createdBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnCreatedBy)
				if c.createdBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.createdBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserGroupsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[1])
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
			if c.groupID.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.groupID.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnGroupID)
				if c.groupID.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.groupID.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserGroupsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Values[1])
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserGroupsColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
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
			if c.updatedBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.updatedBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
				if c.updatedBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.updatedBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserGroupsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[1])
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
			if c.userID.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.userID.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserGroupsColumnUserID)
				if c.userID.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.userID.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserGroupsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Values[1])
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type userGroupsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserGroupsRows(rows *sql.Rows) ([]*userGroupsEntity, error) {
	var (
		entities []*userGroupsEntity
		err      error
	)
	for rows.Next() {
		var ent userGroupsEntity
		err = rows.Scan(
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.GroupID,
			&ent.UpdatedAt,
			&ent.UpdatedBy,
			&ent.UserID,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *userGroupsRepositoryBase) Count(c *userGroupsCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *userGroupsRepositoryBase) Find(c *userGroupsCriteria) ([]*userGroupsEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserGroupsRows(rows)
}
func (r *userGroupsRepositoryBase) FindIter(c *userGroupsCriteria) (*userGroupsIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &userGroupsIterator{rows: rows}, nil
}
func (r *userGroupsRepositoryBase) Insert(e *userGroupsEntity) (*userGroupsEntity, error) {
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
	Author              *userEntity
	Modifier            *userEntity
}

func (e *groupPermissionsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tableGroupPermissionsColumnCreatedAt:
		return &e.CreatedAt, true
	case tableGroupPermissionsColumnCreatedBy:
		return &e.CreatedBy, true
	case tableGroupPermissionsColumnGroupID:
		return &e.GroupID, true
	case tableGroupPermissionsColumnPermissionAction:
		return &e.PermissionAction, true
	case tableGroupPermissionsColumnPermissionModule:
		return &e.PermissionModule, true
	case tableGroupPermissionsColumnPermissionSubsystem:
		return &e.PermissionSubsystem, true
	case tableGroupPermissionsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case tableGroupPermissionsColumnUpdatedBy:
		return &e.UpdatedBy, true
	default:
		return nil, false
	}
}
func (e *groupPermissionsEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// groupPermissionsIterator is not thread safe.
type groupPermissionsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *groupPermissionsIterator) Next() bool {
	return i.rows.Next()
}

func (i *groupPermissionsIterator) Close() error {
	return i.rows.Close()
}

func (i *groupPermissionsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *groupPermissionsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround groupPermissions method that makes iterator more generic.
func (i *groupPermissionsIterator) Ent() (interface{}, error) {
	return i.GroupPermissions()
}

func (i *groupPermissionsIterator) GroupPermissions() (*groupPermissionsEntity, error) {
	var ent groupPermissionsEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type groupPermissionsCriteria struct {
	offset, limit       int64
	sort                map[string]bool
	createdAt           *qtypes.Timestamp
	createdBy           *qtypes.Int64
	groupID             *qtypes.Int64
	permissionAction    *qtypes.String
	permissionModule    *qtypes.String
	permissionSubsystem *qtypes.String
	updatedAt           *qtypes.Timestamp
	updatedBy           *qtypes.Int64
}

func (c *groupPermissionsCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableGroupPermissionsColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.createdBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.createdBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
				if c.createdBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.createdBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[1])
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
			if c.groupID.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.groupID.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnGroupID)
				if c.groupID.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.groupID.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupPermissionsColumnGroupID)
			if c.groupID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.groupID.Values[1])
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
			if c.permissionAction.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionAction.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionAction.Value()))
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
			if c.permissionModule.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionModule.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionModule.Value()))
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
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionSubsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionSubsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableGroupPermissionsColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
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
			if c.updatedBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.updatedBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
				if c.updatedBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.updatedBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableGroupPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[1])
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type groupPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanGroupPermissionsRows(rows *sql.Rows) ([]*groupPermissionsEntity, error) {
	var (
		entities []*groupPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent groupPermissionsEntity
		err = rows.Scan(
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.GroupID,
			&ent.PermissionAction,
			&ent.PermissionModule,
			&ent.PermissionSubsystem,
			&ent.UpdatedAt,
			&ent.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *groupPermissionsRepositoryBase) Count(c *groupPermissionsCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *groupPermissionsRepositoryBase) Find(c *groupPermissionsCriteria) ([]*groupPermissionsEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanGroupPermissionsRows(rows)
}
func (r *groupPermissionsRepositoryBase) FindIter(c *groupPermissionsCriteria) (*groupPermissionsIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &groupPermissionsIterator{rows: rows}, nil
}
func (r *groupPermissionsRepositoryBase) Insert(e *groupPermissionsEntity) (*groupPermissionsEntity, error) {
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
	Author              *userEntity
	Modifier            *userEntity
}

func (e *userPermissionsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserPermissionsColumnCreatedAt:
		return &e.CreatedAt, true
	case tableUserPermissionsColumnCreatedBy:
		return &e.CreatedBy, true
	case tableUserPermissionsColumnPermissionAction:
		return &e.PermissionAction, true
	case tableUserPermissionsColumnPermissionModule:
		return &e.PermissionModule, true
	case tableUserPermissionsColumnPermissionSubsystem:
		return &e.PermissionSubsystem, true
	case tableUserPermissionsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case tableUserPermissionsColumnUpdatedBy:
		return &e.UpdatedBy, true
	case tableUserPermissionsColumnUserID:
		return &e.UserID, true
	default:
		return nil, false
	}
}
func (e *userPermissionsEntity) Props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.Prop(cn); ok {
			res = append(res, prop)
		} else {
			return nil, fmt.Errorf("unexpected column provided: %s", cn)
		}
	}
	return res, nil
}

// userPermissionsIterator is not thread safe.
type userPermissionsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *userPermissionsIterator) Next() bool {
	return i.rows.Next()
}

func (i *userPermissionsIterator) Close() error {
	return i.rows.Close()
}

func (i *userPermissionsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *userPermissionsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper arround userPermissions method that makes iterator more generic.
func (i *userPermissionsIterator) Ent() (interface{}, error) {
	return i.UserPermissions()
}

func (i *userPermissionsIterator) UserPermissions() (*userPermissionsEntity, error) {
	var ent userPermissionsEntity
	cols, err := i.rows.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type userPermissionsCriteria struct {
	offset, limit       int64
	sort                map[string]bool
	createdAt           *qtypes.Timestamp
	createdBy           *qtypes.Int64
	permissionAction    *qtypes.String
	permissionModule    *qtypes.String
	permissionSubsystem *qtypes.String
	updatedAt           *qtypes.Timestamp
	updatedBy           *qtypes.Int64
	userID              *qtypes.Int64
}

func (c *userPermissionsCriteria) WriteSQL(b *bytes.Buffer, pw *pqtgo.PlaceholderWriter, args *pqtgo.Arguments) (wr int64, err error) {
	var (
		wrt   int
		wrt64 int64
		dirty bool
	)

	wbuf := bytes.NewBuffer(nil)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return wr, err
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
				if c.createdAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.createdAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(createdAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserPermissionsColumnCreatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(createdAt2)
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
			if c.createdBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.createdBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
				if c.createdBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.createdBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserPermissionsColumnCreatedBy)
			if c.createdBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.createdBy.Values[1])
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
			if c.permissionAction.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionAction.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionAction.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionAction)
			if c.permissionAction.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionAction.Value()))
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
			if c.permissionModule.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionModule.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionModule.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionModule)
			if c.permissionModule.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionModule.Value()))
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
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString(" = ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.permissionSubsystem.Value())
		case qtypes.TextQueryType_SUBSTRING:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_PREFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%s%%", c.permissionSubsystem.Value()))
		case qtypes.TextQueryType_HAS_SUFFIX:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnPermissionSubsystem)
			if c.permissionSubsystem.Negation {
				wbuf.WriteString(" NOT LIKE ")
			} else {
				wbuf.WriteString(" LIKE ")
			}
			pw.WriteTo(wbuf)
			args.Add(fmt.Sprintf("%%%s", c.permissionSubsystem.Value()))
		}
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return wr, err
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
				if c.updatedAt.Negation {
					wbuf.WriteString("<>")
				} else {
					wbuf.WriteString("=")
				}
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString(">")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_GREATER_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString(">=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("<")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_LESS_EQUAL:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
				wbuf.WriteString("<=")
				pw.WriteTo(wbuf)
				args.Add(c.updatedAt.Value())
			case qtypes.NumericQueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if dirty {
						wbuf.WriteString(" AND ")
					}
					dirty = true

					wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
					wbuf.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							wbuf.WriteString(",")
						}
						pw.WriteTo(wbuf)
						args.Add(v)
					}
					wbuf.WriteString(") ")
				}
			case qtypes.NumericQueryType_BETWEEN:
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return wr, err
					}
					wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
					wbuf.WriteString(" > ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt1)
					wbuf.WriteString(" AND ")
					wbuf.WriteString(tableUserPermissionsColumnUpdatedAt)
					wbuf.WriteString(" < ")
					pw.WriteTo(wbuf)
					args.Add(updatedAt2)
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
			if c.updatedBy.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.updatedBy.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
				if c.updatedBy.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.updatedBy.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserPermissionsColumnUpdatedBy)
			if c.updatedBy.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.updatedBy.Values[1])
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
			if c.userID.Negation {
				wbuf.WriteString(" <> ")
			} else {
				wbuf.WriteString("=")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_GREATER:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_GREATER_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" < ")
			} else {
				wbuf.WriteString(" >= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_LESS:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID)
		case qtypes.NumericQueryType_LESS_EQUAL:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" > ")
			} else {
				wbuf.WriteString(" <= ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Value())
		case qtypes.NumericQueryType_IN:
			if len(c.userID.Values) > 0 {
				if dirty {
					wbuf.WriteString(" AND ")
				}
				dirty = true

				wbuf.WriteString(tableUserPermissionsColumnUserID)
				if c.userID.Negation {
					wbuf.WriteString(" NOT IN (")
				} else {
					wbuf.WriteString(" IN (")
				}
				for i, v := range c.userID.Values {
					if i != 0 {
						wbuf.WriteString(",")
					}
					pw.WriteTo(wbuf)
					args.Add(v)
				}
				wbuf.WriteString(") ")
			}
		case qtypes.NumericQueryType_BETWEEN:
			if dirty {
				wbuf.WriteString(" AND ")
			}
			dirty = true

			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" <= ")
			} else {
				wbuf.WriteString(" > ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Values[0])
			wbuf.WriteString(" AND ")
			wbuf.WriteString(tableUserPermissionsColumnUserID)
			if c.userID.Negation {
				wbuf.WriteString(" >= ")
			} else {
				wbuf.WriteString(" < ")
			}
			pw.WriteTo(wbuf)
			args.Add(c.userID.Values[1])
		}
	}

	if dirty {
		if wrt, err = b.WriteString(" WHERE "); err != nil {
			return
		}
		wr += int64(wrt)
		if wrt64, err = wbuf.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
	}

	if c.offset > 0 {
		b.WriteString(" OFFSET ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.offset)
	}
	if c.limit > 0 {
		b.WriteString(" LIMIT ")
		if wrt64, err = pw.WriteTo(b); err != nil {
			return
		}
		wr += wrt64
		args.Add(c.limit)
	}

	return
}

type userPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserPermissionsRows(rows *sql.Rows) ([]*userPermissionsEntity, error) {
	var (
		entities []*userPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent userPermissionsEntity
		err = rows.Scan(
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.PermissionAction,
			&ent.PermissionModule,
			&ent.PermissionSubsystem,
			&ent.UpdatedAt,
			&ent.UpdatedBy,
			&ent.UserID,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

func (r *userPermissionsRepositoryBase) Count(c *userPermissionsCriteria) (int64, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT COUNT(*) FROM ")
	qbuf.WriteString(r.table)
	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return 0, err
	}
	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	err := r.db.QueryRow(qbuf.String(), args.Slice()...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *userPermissionsRepositoryBase) Find(c *userPermissionsCriteria) ([]*userPermissionsEntity, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserPermissionsRows(rows)
}
func (r *userPermissionsRepositoryBase) FindIter(c *userPermissionsCriteria) (*userPermissionsIterator, error) {

	qbuf := bytes.NewBuffer(nil)
	qbuf.WriteString("SELECT ")
	qbuf.WriteString(strings.Join(r.columns, ", "))
	qbuf.WriteString(" FROM ")
	qbuf.WriteString(r.table)

	pw := pqtgo.NewPlaceholderWriter()
	args := pqtgo.NewArguments(0)

	if _, err := c.WriteSQL(qbuf, pw, args); err != nil {
		return nil, err
	}

	if r.dbg {
		if err := r.log.Log("msg", qbuf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(qbuf.String(), args.Slice()...)
	if err != nil {
		return nil, err
	}

	return &userPermissionsIterator{rows: rows}, nil
}
func (r *userPermissionsRepositoryBase) Insert(e *userPermissionsEntity) (*userPermissionsEntity, error) {
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
