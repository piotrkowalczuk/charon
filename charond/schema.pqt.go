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

func (c *userCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {
	if c.confirmationToken != nil {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		if _, err = com.WriteString(tableUserColumnConfirmationToken); err != nil {
			return
		}
		if _, err = com.WriteString(" = "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}

		if com.Dirty {
			if opt.Cast != "" {
				if _, err = com.WriteString(opt.Cast); err != nil {
					return
				}
			} else {
				if _, err = com.WriteString(" "); err != nil {
					return
				}
			}
		}

		com.Add(c.confirmationToken)
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.createdBy, tableUserColumnCreatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.firstName, tableUserColumnFirstName, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.id, tableUserColumnID, com, pqtgo.And); err != nil {
		return
	}
	if c.isActive != nil && c.isActive.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(tableUserColumnIsActive)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.isActive)
	}
	if c.isConfirmed != nil && c.isConfirmed.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(tableUserColumnIsConfirmed)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.isConfirmed)
	}
	if c.isStaff != nil && c.isStaff.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(tableUserColumnIsStaff)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.isStaff)
	}
	if c.isSuperuser != nil && c.isSuperuser.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(tableUserColumnIsSuperuser)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.isSuperuser)
	}

	if c.lastLoginAt != nil && c.lastLoginAt.Valid {
		lastLoginAtt1 := c.lastLoginAt.Value()
		if lastLoginAtt1 != nil {
			lastLoginAt1, err := ptypes.Timestamp(lastLoginAtt1)
			if err != nil {
				return err
			}
			switch c.lastLoginAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				if c.lastLoginAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				if c.lastLoginAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.lastLoginAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.lastLoginAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.lastLoginAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.lastLoginAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnLastLoginAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.lastLoginAt.Value())
			case qtypes.QueryType_IN:
				if len(c.lastLoginAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserColumnLastLoginAt)
					com.WriteString(" IN (")
					for i, v := range c.lastLoginAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				lastLoginAtt2 := c.lastLoginAt.Values[1]
				if lastLoginAtt2 != nil {
					lastLoginAt2, err := ptypes.Timestamp(lastLoginAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserColumnLastLoginAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(lastLoginAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserColumnLastLoginAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(lastLoginAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryString(c.lastName, tableUserColumnLastName, com, pqtgo.And); err != nil {
		return
	}
	if c.password != nil {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		if _, err = com.WriteString(tableUserColumnPassword); err != nil {
			return
		}
		if _, err = com.WriteString(" = "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}

		if com.Dirty {
			if opt.Cast != "" {
				if _, err = com.WriteString(opt.Cast); err != nil {
					return
				}
			} else {
				if _, err = com.WriteString(" "); err != nil {
					return
				}
			}
		}

		com.Add(c.password)
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.updatedBy, tableUserColumnUpdatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.username, tableUserColumnUsername, com, pqtgo.And); err != nil {
		return
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
	}

	return
}

type userPatch struct {
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

	com := pqtgo.NewComposer(15)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepositoryBase) Find(c *userCriteria) ([]*userEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserRows(rows)
}
func (r *userRepositoryBase) FindIter(c *userCriteria) (*userIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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
func (r *userRepositoryBase) UpdateOneByID(id int64, patch *userPatch) (*userEntity, error) {
	update := pqcomp.New(0, 15)
	update.AddExpr(tableUserColumnID, pqcomp.Equal, id)
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
func (r *userRepositoryBase) DeleteOneByID(id int64) (int64, error) {
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

func (c *groupCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableGroupColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableGroupColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tableGroupColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.createdBy, tableGroupColumnCreatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.description, tableGroupColumnDescription, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.id, tableGroupColumnID, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.name, tableGroupColumnName, com, pqtgo.And); err != nil {
		return
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableGroupColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableGroupColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tableGroupColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.updatedBy, tableGroupColumnUpdatedBy, com, pqtgo.And); err != nil {
		return
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
	}

	return
}

type groupPatch struct {
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

	com := pqtgo.NewComposer(7)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *groupRepositoryBase) Find(c *groupCriteria) ([]*groupEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanGroupRows(rows)
}
func (r *groupRepositoryBase) FindIter(c *groupCriteria) (*groupIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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
func (r *groupRepositoryBase) UpdateOneByID(id int64, patch *groupPatch) (*groupEntity, error) {
	update := pqcomp.New(0, 7)
	update.AddExpr(tableGroupColumnID, pqcomp.Equal, id)
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
		return nil, errors.New("group update failure, nothing to update")
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
func (r *groupRepositoryBase) DeleteOneByID(id int64) (int64, error) {
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

func (c *permissionCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if err = pqtgo.WriteCompositionQueryString(c.action, tablePermissionColumnAction, com, pqtgo.And); err != nil {
		return
	}

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tablePermissionColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tablePermissionColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tablePermissionColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.id, tablePermissionColumnID, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.module, tablePermissionColumnModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.subsystem, tablePermissionColumnSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tablePermissionColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tablePermissionColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tablePermissionColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tablePermissionColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
	}

	return
}

type permissionPatch struct {
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

	com := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *permissionRepositoryBase) Find(c *permissionCriteria) ([]*permissionEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanPermissionRows(rows)
}
func (r *permissionRepositoryBase) FindIter(c *permissionCriteria) (*permissionIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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
func (r *permissionRepositoryBase) UpdateOneByID(id int64, patch *permissionPatch) (*permissionEntity, error) {
	update := pqcomp.New(0, 6)
	update.AddExpr(tablePermissionColumnID, pqcomp.Equal, id)
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
func (r *permissionRepositoryBase) DeleteOneByID(id int64) (int64, error) {
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

func (c *userGroupsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserGroupsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserGroupsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserGroupsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.createdBy, tableUserGroupsColumnCreatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.groupID, tableUserGroupsColumnGroupID, com, pqtgo.And); err != nil {
		return
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserGroupsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserGroupsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserGroupsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserGroupsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.updatedBy, tableUserGroupsColumnUpdatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.userID, tableUserGroupsColumnUserID, com, pqtgo.And); err != nil {
		return
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
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

	com := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userGroupsRepositoryBase) Find(c *userGroupsCriteria) ([]*userGroupsEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserGroupsRows(rows)
}
func (r *userGroupsRepositoryBase) FindIter(c *userGroupsCriteria) (*userGroupsIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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

func (c *groupPermissionsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableGroupPermissionsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableGroupPermissionsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tableGroupPermissionsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.createdBy, tableGroupPermissionsColumnCreatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.groupID, tableGroupPermissionsColumnGroupID, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionAction, tableGroupPermissionsColumnPermissionAction, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionModule, tableGroupPermissionsColumnPermissionModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionSubsystem, tableGroupPermissionsColumnPermissionSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableGroupPermissionsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.updatedBy, tableGroupPermissionsColumnUpdatedBy, com, pqtgo.And); err != nil {
		return
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
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

	com := pqtgo.NewComposer(8)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *groupPermissionsRepositoryBase) Find(c *groupPermissionsCriteria) ([]*groupPermissionsEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanGroupPermissionsRows(rows)
}
func (r *groupPermissionsRepositoryBase) FindIter(c *groupPermissionsCriteria) (*groupPermissionsIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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

func (c *userPermissionsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return err
			}
			switch c.createdAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				if c.createdAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.createdAt.Value())
			case qtypes.QueryType_IN:
				if len(c.createdAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserPermissionsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.createdAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserPermissionsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(createdAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserPermissionsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(createdAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.createdBy, tableUserPermissionsColumnCreatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionAction, tableUserPermissionsColumnPermissionAction, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionModule, tableUserPermissionsColumnPermissionModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.permissionSubsystem, tableUserPermissionsColumnPermissionSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return err
			}
			switch c.updatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				if c.updatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(tableUserPermissionsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.updatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.updatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(tableUserPermissionsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.updatedAt.Values {
						if i != 0 {
							com.WriteString(",")
						}
						com.WritePlaceholder()
						com.Add(v)
					}
					com.WriteString(") ")
				}
			case qtypes.QueryType_BETWEEN:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(tableUserPermissionsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(updatedAt1)
					com.WriteString(" AND ")
					com.WriteString(tableUserPermissionsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(updatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.updatedBy, tableUserPermissionsColumnUpdatedBy, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.userID, tableUserPermissionsColumnUserID, com, pqtgo.And); err != nil {
		return
	}

	if c.offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.offset)
	}
	if c.limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.limit)
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

	com := pqtgo.NewComposer(8)
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM ")
	buf.WriteString(r.table)

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return 0, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Count"); err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.QueryRow(buf.String(), com.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userPermissionsRepositoryBase) Find(c *userPermissionsCriteria) ([]*userPermissionsEntity, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return ScanUserPermissionsRows(rows)
}
func (r *userPermissionsRepositoryBase) FindIter(c *userPermissionsCriteria) (*userPermissionsIterator, error) {

	com := pqtgo.NewComposer(1)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(r.columns, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.table)
	buf.WriteString(" ")

	if err := c.WriteComposition("", com, pqtgo.And); err != nil {
		return nil, err
	}
	if com.Dirty {
		buf.WriteString(" WHERE ")
	}
	if com.Len() > 0 {
		buf.ReadFrom(com)
	}

	if r.dbg {
		if err := r.log.Log("msg", buf.String(), "function", "Find"); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query(buf.String(), com.Args()...)
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
