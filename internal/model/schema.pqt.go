package model

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
	TableUser                              = "charon.user"
	TableUserColumnConfirmationToken       = "confirmation_token"
	TableUserColumnCreatedAt               = "created_at"
	TableUserColumnCreatedBy               = "created_by"
	TableUserColumnFirstName               = "first_name"
	TableUserColumnID                      = "id"
	TableUserColumnIsActive                = "is_active"
	TableUserColumnIsConfirmed             = "is_confirmed"
	TableUserColumnIsStaff                 = "is_staff"
	TableUserColumnIsSuperuser             = "is_superuser"
	TableUserColumnLastLoginAt             = "last_login_at"
	TableUserColumnLastName                = "last_name"
	TableUserColumnPassword                = "password"
	TableUserColumnUpdatedAt               = "updated_at"
	TableUserColumnUpdatedBy               = "updated_by"
	TableUserColumnUsername                = "username"
	TableUserConstraintCreatedByForeignKey = "charon.user_created_by_fkey"
	TableUserConstraintPrimaryKey          = "charon.user_id_pkey"
	TableUserConstraintUpdatedByForeignKey = "charon.user_updated_by_fkey"
	TableUserConstraintUsernameUnique      = "charon.user_username_key"
)

var (
	TableUserColumns = []string{
		TableUserColumnConfirmationToken,
		TableUserColumnCreatedAt,
		TableUserColumnCreatedBy,
		TableUserColumnFirstName,
		TableUserColumnID,
		TableUserColumnIsActive,
		TableUserColumnIsConfirmed,
		TableUserColumnIsStaff,
		TableUserColumnIsSuperuser,
		TableUserColumnLastLoginAt,
		TableUserColumnLastName,
		TableUserColumnPassword,
		TableUserColumnUpdatedAt,
		TableUserColumnUpdatedBy,
		TableUserColumnUsername,
	}
)

type UserEntity struct {
	// ConfirmationToken ...
	ConfirmationToken []byte
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy *ntypes.Int64
	// FirstName ...
	FirstName string
	// ID ...
	ID int64
	// IsActive ...
	IsActive bool
	// IsConfirmed ...
	IsConfirmed bool
	// IsStaff ...
	IsStaff bool
	// IsSuperuser ...
	IsSuperuser bool
	// LastLoginAt ...
	LastLoginAt *time.Time
	// LastName ...
	LastName string
	// Password ...
	Password []byte
	// UpdatedAt ...
	UpdatedAt *time.Time
	// UpdatedBy ...
	UpdatedBy *ntypes.Int64
	// Username ...
	Username string
	// Author ...
	Author *UserEntity
	// Modifier ...
	Modifier *UserEntity
	// Permissions ...
	Permissions []*PermissionEntity
	// Groups ...
	Groups []*GroupEntity
}

func (e *UserEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TableUserColumnConfirmationToken:
		return &e.ConfirmationToken, true
	case TableUserColumnCreatedAt:
		return &e.CreatedAt, true
	case TableUserColumnCreatedBy:
		return &e.CreatedBy, true
	case TableUserColumnFirstName:
		return &e.FirstName, true
	case TableUserColumnID:
		return &e.ID, true
	case TableUserColumnIsActive:
		return &e.IsActive, true
	case TableUserColumnIsConfirmed:
		return &e.IsConfirmed, true
	case TableUserColumnIsStaff:
		return &e.IsStaff, true
	case TableUserColumnIsSuperuser:
		return &e.IsSuperuser, true
	case TableUserColumnLastLoginAt:
		return &e.LastLoginAt, true
	case TableUserColumnLastName:
		return &e.LastName, true
	case TableUserColumnPassword:
		return &e.Password, true
	case TableUserColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableUserColumnUpdatedBy:
		return &e.UpdatedBy, true
	case TableUserColumnUsername:
		return &e.Username, true
	default:
		return nil, false
	}
}
func (e *UserEntity) Props(cns ...string) ([]interface{}, error) {

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

// UserIterator is not thread safe.
type UserIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *UserIterator) Next() bool {
	return i.rows.Next()
}

func (i *UserIterator) Close() error {
	return i.rows.Close()
}

func (i *UserIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *UserIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around User method that makes iterator more generic.
func (i *UserIterator) Ent() (interface{}, error) {
	return i.User()
}

func (i *UserIterator) User() (*UserEntity, error) {
	var ent UserEntity
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

type UserCriteria struct {
	Offset, Limit     int64
	Sort              map[string]bool
	ConfirmationToken []byte
	CreatedAt         *qtypes.Timestamp
	CreatedBy         *qtypes.Int64
	FirstName         *qtypes.String
	ID                *qtypes.Int64
	IsActive          *ntypes.Bool
	IsConfirmed       *ntypes.Bool
	IsStaff           *ntypes.Bool
	IsSuperuser       *ntypes.Bool
	LastLoginAt       *qtypes.Timestamp
	LastName          *qtypes.String
	Password          []byte
	UpdatedAt         *qtypes.Timestamp
	UpdatedBy         *qtypes.Int64
	Username          *qtypes.String
}

func (c *UserCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {
	if c.ConfirmationToken != nil {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		if _, err = com.WriteString(TableUserColumnConfirmationToken); err != nil {
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

		com.Add(c.ConfirmationToken)
	}

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.CreatedBy, TableUserColumnCreatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.FirstName, TableUserColumnFirstName, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.ID, TableUserColumnID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}
	if c.IsActive != nil && c.IsActive.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(TableUserColumnIsActive)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.IsActive)
	}
	if c.IsConfirmed != nil && c.IsConfirmed.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(TableUserColumnIsConfirmed)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.IsConfirmed)
	}
	if c.IsStaff != nil && c.IsStaff.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(TableUserColumnIsStaff)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.IsStaff)
	}
	if c.IsSuperuser != nil && c.IsSuperuser.Valid {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		com.WriteString(TableUserColumnIsSuperuser)
		com.WriteString(" = ")
		com.WritePlaceholder()
		com.Add(c.IsSuperuser)
	}

	if c.LastLoginAt != nil && c.LastLoginAt.Valid {
		LastLoginAtt1 := c.LastLoginAt.Value()
		if LastLoginAtt1 != nil {
			LastLoginAt1, err := ptypes.Timestamp(LastLoginAtt1)
			if err != nil {
				return err
			}
			switch c.LastLoginAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				if c.LastLoginAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				if c.LastLoginAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.LastLoginAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.LastLoginAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.LastLoginAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.LastLoginAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnLastLoginAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.LastLoginAt.Value())
			case qtypes.QueryType_IN:
				if len(c.LastLoginAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserColumnLastLoginAt)
					com.WriteString(" IN (")
					for i, v := range c.LastLoginAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				LastLoginAtt2 := c.LastLoginAt.Values[1]
				if LastLoginAtt2 != nil {
					LastLoginAt2, err := ptypes.Timestamp(LastLoginAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserColumnLastLoginAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(LastLoginAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserColumnLastLoginAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(LastLoginAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryString(c.LastName, TableUserColumnLastName, com, pqtgo.And); err != nil {
		return
	}
	if c.Password != nil {
		if com.Dirty {
			com.WriteString(" AND ")
		}
		com.Dirty = true
		if _, err = com.WriteString(TableUserColumnPassword); err != nil {
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

		com.Add(c.Password)
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UpdatedBy, TableUserColumnUpdatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.Username, TableUserColumnUsername, com, pqtgo.And); err != nil {
		return
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TableUserColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type UserPatch struct {
	ConfirmationToken []byte
	CreatedAt         *time.Time
	CreatedBy         *ntypes.Int64
	FirstName         *ntypes.String
	IsActive          *ntypes.Bool
	IsConfirmed       *ntypes.Bool
	IsStaff           *ntypes.Bool
	IsSuperuser       *ntypes.Bool
	LastLoginAt       *time.Time
	LastName          *ntypes.String
	Password          []byte
	UpdatedAt         *time.Time
	UpdatedBy         *ntypes.Int64
	Username          *ntypes.String
}

type UserRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserRows(rows *sql.Rows) ([]*UserEntity, error) {
	var (
		entities []*UserEntity
		err      error
	)
	for rows.Next() {
		var ent UserEntity
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

func (r *UserRepositoryBase) Count(c *UserCriteria) (int64, error) {

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

func (r *UserRepositoryBase) Find(c *UserCriteria) ([]*UserEntity, error) {

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
func (r *UserRepositoryBase) FindIter(c *UserCriteria) (*UserIterator, error) {

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

	return &UserIterator{rows: rows}, nil
}
func (r *UserRepositoryBase) FindOneByID(id int64) (*UserEntity, error) {
	var (
		ent UserEntity
	)
	query := `SELECT confirmation_token,
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

	return &ent, nil
}
func (r *UserRepositoryBase) FindOneByUsername(userUsername string) (*UserEntity, error) {
	var (
		ent UserEntity
	)
	query := `SELECT confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username FROM charon.user WHERE username = $1`
	err := r.db.QueryRow(query, userUsername).Scan(
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

	return &ent, nil
}
func (r *UserRepositoryBase) Insert(e *UserEntity) (*UserEntity, error) {
	insert := pqcomp.New(0, 15)
	insert.AddExpr(TableUserColumnConfirmationToken, "", e.ConfirmationToken)
	insert.AddExpr(TableUserColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserColumnFirstName, "", e.FirstName)
	insert.AddExpr(TableUserColumnIsActive, "", e.IsActive)
	insert.AddExpr(TableUserColumnIsConfirmed, "", e.IsConfirmed)
	insert.AddExpr(TableUserColumnIsStaff, "", e.IsStaff)
	insert.AddExpr(TableUserColumnIsSuperuser, "", e.IsSuperuser)
	insert.AddExpr(TableUserColumnLastLoginAt, "", e.LastLoginAt)
	insert.AddExpr(TableUserColumnLastName, "", e.LastName)
	insert.AddExpr(TableUserColumnPassword, "", e.Password)
	insert.AddExpr(TableUserColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserColumnUsername, "", e.Username)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *UserRepositoryBase) Upsert(e *UserEntity, p *UserPatch, inf ...string) (*UserEntity, error) {
	insert := pqcomp.New(0, 15)
	update := insert.Compose(15)
	insert.AddExpr(TableUserColumnConfirmationToken, "", e.ConfirmationToken)
	insert.AddExpr(TableUserColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserColumnFirstName, "", e.FirstName)
	insert.AddExpr(TableUserColumnIsActive, "", e.IsActive)
	insert.AddExpr(TableUserColumnIsConfirmed, "", e.IsConfirmed)
	insert.AddExpr(TableUserColumnIsStaff, "", e.IsStaff)
	insert.AddExpr(TableUserColumnIsSuperuser, "", e.IsSuperuser)
	insert.AddExpr(TableUserColumnLastLoginAt, "", e.LastLoginAt)
	insert.AddExpr(TableUserColumnLastName, "", e.LastName)
	insert.AddExpr(TableUserColumnPassword, "", e.Password)
	insert.AddExpr(TableUserColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserColumnUsername, "", e.Username)
	if len(inf) > 0 {
		update.AddExpr(TableUserColumnConfirmationToken, "=", p.ConfirmationToken)
		update.AddExpr(TableUserColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TableUserColumnCreatedBy, "=", p.CreatedBy)
		update.AddExpr(TableUserColumnFirstName, "=", p.FirstName)
		update.AddExpr(TableUserColumnIsActive, "=", p.IsActive)
		update.AddExpr(TableUserColumnIsConfirmed, "=", p.IsConfirmed)
		update.AddExpr(TableUserColumnIsStaff, "=", p.IsStaff)
		update.AddExpr(TableUserColumnIsSuperuser, "=", p.IsSuperuser)
		update.AddExpr(TableUserColumnLastLoginAt, "=", p.LastLoginAt)
		update.AddExpr(TableUserColumnLastName, "=", p.LastName)
		update.AddExpr(TableUserColumnPassword, "=", p.Password)
		update.AddExpr(TableUserColumnUpdatedAt, "=", p.UpdatedAt)
		update.AddExpr(TableUserColumnUpdatedBy, "=", p.UpdatedBy)
		update.AddExpr(TableUserColumnUsername, "=", p.Username)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *UserRepositoryBase) UpdateOneByID(id int64, patch *UserPatch) (*UserEntity, error) {
	update := pqcomp.New(1, 15)
	update.AddArg(id)

	update.AddExpr(TableUserColumnConfirmationToken, pqcomp.Equal, patch.ConfirmationToken)
	if patch.CreatedAt != nil {
		update.AddExpr(TableUserColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableUserColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableUserColumnFirstName, pqcomp.Equal, patch.FirstName)
	update.AddExpr(TableUserColumnIsActive, pqcomp.Equal, patch.IsActive)
	update.AddExpr(TableUserColumnIsConfirmed, pqcomp.Equal, patch.IsConfirmed)
	update.AddExpr(TableUserColumnIsStaff, pqcomp.Equal, patch.IsStaff)
	update.AddExpr(TableUserColumnIsSuperuser, pqcomp.Equal, patch.IsSuperuser)
	update.AddExpr(TableUserColumnLastLoginAt, pqcomp.Equal, patch.LastLoginAt)
	update.AddExpr(TableUserColumnLastName, pqcomp.Equal, patch.LastName)
	update.AddExpr(TableUserColumnPassword, pqcomp.Equal, patch.Password)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableUserColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableUserColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableUserColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)
	update.AddExpr(TableUserColumnUsername, pqcomp.Equal, patch.Username)

	if update.Len() == 0 {
		return nil, errors.New("User update failure, nothing to update")
	}
	query := "UPDATE charon.user SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e UserEntity
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
func (r *UserRepositoryBase) UpdateOneByUsername(userUsername string, patch *UserPatch) (*UserEntity, error) {
	update := pqcomp.New(1, 15)
	update.AddArg(userUsername)
	update.AddExpr(TableUserColumnConfirmationToken, pqcomp.Equal, patch.ConfirmationToken)
	if patch.CreatedAt != nil {
		update.AddExpr(TableUserColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableUserColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableUserColumnFirstName, pqcomp.Equal, patch.FirstName)
	update.AddExpr(TableUserColumnIsActive, pqcomp.Equal, patch.IsActive)
	update.AddExpr(TableUserColumnIsConfirmed, pqcomp.Equal, patch.IsConfirmed)
	update.AddExpr(TableUserColumnIsStaff, pqcomp.Equal, patch.IsStaff)
	update.AddExpr(TableUserColumnIsSuperuser, pqcomp.Equal, patch.IsSuperuser)
	update.AddExpr(TableUserColumnLastLoginAt, pqcomp.Equal, patch.LastLoginAt)
	update.AddExpr(TableUserColumnLastName, pqcomp.Equal, patch.LastName)
	update.AddExpr(TableUserColumnPassword, pqcomp.Equal, patch.Password)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableUserColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableUserColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableUserColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)
	update.AddExpr(TableUserColumnUsername, pqcomp.Equal, patch.Username)

	if update.Len() == 0 {
		return nil, errors.New("User update failure, nothing to update")
	}
	query := "UPDATE charon.user SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE username = $1 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneByUsername"); err != nil {
			return nil, err
		}
	}
	var e UserEntity
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

func (r *UserRepositoryBase) DeleteOneByID(id int64) (int64, error) {
	query := "DELETE FROM charon.user WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	TableGroup                              = "charon.group"
	TableGroupColumnCreatedAt               = "created_at"
	TableGroupColumnCreatedBy               = "created_by"
	TableGroupColumnDescription             = "description"
	TableGroupColumnID                      = "id"
	TableGroupColumnName                    = "name"
	TableGroupColumnUpdatedAt               = "updated_at"
	TableGroupColumnUpdatedBy               = "updated_by"
	TableGroupConstraintCreatedByForeignKey = "charon.group_created_by_fkey"
	TableGroupConstraintPrimaryKey          = "charon.group_id_pkey"
	TableGroupConstraintNameUnique          = "charon.group_name_key"
	TableGroupConstraintUpdatedByForeignKey = "charon.group_updated_by_fkey"
)

var (
	TableGroupColumns = []string{
		TableGroupColumnCreatedAt,
		TableGroupColumnCreatedBy,
		TableGroupColumnDescription,
		TableGroupColumnID,
		TableGroupColumnName,
		TableGroupColumnUpdatedAt,
		TableGroupColumnUpdatedBy,
	}
)

type GroupEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy *ntypes.Int64
	// Description ...
	Description *ntypes.String
	// ID ...
	ID int64
	// Name ...
	Name string
	// UpdatedAt ...
	UpdatedAt *time.Time
	// UpdatedBy ...
	UpdatedBy *ntypes.Int64
	// Author ...
	Author *UserEntity
	// Modifier ...
	Modifier *UserEntity
	// Permissions ...
	Permissions []*PermissionEntity
	// Users ...
	Users []*UserEntity
}

func (e *GroupEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TableGroupColumnCreatedAt:
		return &e.CreatedAt, true
	case TableGroupColumnCreatedBy:
		return &e.CreatedBy, true
	case TableGroupColumnDescription:
		return &e.Description, true
	case TableGroupColumnID:
		return &e.ID, true
	case TableGroupColumnName:
		return &e.Name, true
	case TableGroupColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableGroupColumnUpdatedBy:
		return &e.UpdatedBy, true
	default:
		return nil, false
	}
}
func (e *GroupEntity) Props(cns ...string) ([]interface{}, error) {

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

// GroupIterator is not thread safe.
type GroupIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *GroupIterator) Next() bool {
	return i.rows.Next()
}

func (i *GroupIterator) Close() error {
	return i.rows.Close()
}

func (i *GroupIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *GroupIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around Group method that makes iterator more generic.
func (i *GroupIterator) Ent() (interface{}, error) {
	return i.Group()
}

func (i *GroupIterator) Group() (*GroupEntity, error) {
	var ent GroupEntity
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

type GroupCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	CreatedAt     *qtypes.Timestamp
	CreatedBy     *qtypes.Int64
	Description   *qtypes.String
	ID            *qtypes.Int64
	Name          *qtypes.String
	UpdatedAt     *qtypes.Timestamp
	UpdatedBy     *qtypes.Int64
}

func (c *GroupCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableGroupColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableGroupColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableGroupColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.CreatedBy, TableGroupColumnCreatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.Description, TableGroupColumnDescription, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.ID, TableGroupColumnID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.Name, TableGroupColumnName, com, pqtgo.And); err != nil {
		return
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableGroupColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableGroupColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableGroupColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UpdatedBy, TableGroupColumnUpdatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TableGroupColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type GroupPatch struct {
	CreatedAt   *time.Time
	CreatedBy   *ntypes.Int64
	Description *ntypes.String
	Name        *ntypes.String
	UpdatedAt   *time.Time
	UpdatedBy   *ntypes.Int64
}

type GroupRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanGroupRows(rows *sql.Rows) ([]*GroupEntity, error) {
	var (
		entities []*GroupEntity
		err      error
	)
	for rows.Next() {
		var ent GroupEntity
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

func (r *GroupRepositoryBase) Count(c *GroupCriteria) (int64, error) {

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

func (r *GroupRepositoryBase) Find(c *GroupCriteria) ([]*GroupEntity, error) {

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
func (r *GroupRepositoryBase) FindIter(c *GroupCriteria) (*GroupIterator, error) {

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

	return &GroupIterator{rows: rows}, nil
}
func (r *GroupRepositoryBase) FindOneByID(id int64) (*GroupEntity, error) {
	var (
		ent GroupEntity
	)
	query := `SELECT created_at,
created_by,
description,
id,
name,
updated_at,
updated_by
 FROM charon.group WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
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

	return &ent, nil
}
func (r *GroupRepositoryBase) FindOneByName(groupName string) (*GroupEntity, error) {
	var (
		ent GroupEntity
	)
	query := `SELECT created_at, created_by, description, id, name, updated_at, updated_by FROM charon.group WHERE name = $1`
	err := r.db.QueryRow(query, groupName).Scan(
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

	return &ent, nil
}
func (r *GroupRepositoryBase) Insert(e *GroupEntity) (*GroupEntity, error) {
	insert := pqcomp.New(0, 7)
	insert.AddExpr(TableGroupColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableGroupColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableGroupColumnDescription, "", e.Description)
	insert.AddExpr(TableGroupColumnName, "", e.Name)
	insert.AddExpr(TableGroupColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableGroupColumnUpdatedBy, "", e.UpdatedBy)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *GroupRepositoryBase) Upsert(e *GroupEntity, p *GroupPatch, inf ...string) (*GroupEntity, error) {
	insert := pqcomp.New(0, 7)
	update := insert.Compose(7)
	insert.AddExpr(TableGroupColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableGroupColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableGroupColumnDescription, "", e.Description)
	insert.AddExpr(TableGroupColumnName, "", e.Name)
	insert.AddExpr(TableGroupColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableGroupColumnUpdatedBy, "", e.UpdatedBy)
	if len(inf) > 0 {
		update.AddExpr(TableGroupColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TableGroupColumnCreatedBy, "=", p.CreatedBy)
		update.AddExpr(TableGroupColumnDescription, "=", p.Description)
		update.AddExpr(TableGroupColumnName, "=", p.Name)
		update.AddExpr(TableGroupColumnUpdatedAt, "=", p.UpdatedAt)
		update.AddExpr(TableGroupColumnUpdatedBy, "=", p.UpdatedBy)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *GroupRepositoryBase) UpdateOneByID(id int64, patch *GroupPatch) (*GroupEntity, error) {
	update := pqcomp.New(1, 7)
	update.AddArg(id)

	if patch.CreatedAt != nil {
		update.AddExpr(TableGroupColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableGroupColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableGroupColumnDescription, pqcomp.Equal, patch.Description)
	update.AddExpr(TableGroupColumnName, pqcomp.Equal, patch.Name)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableGroupColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableGroupColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableGroupColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)

	if update.Len() == 0 {
		return nil, errors.New("Group update failure, nothing to update")
	}
	query := "UPDATE charon.group SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e GroupEntity
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
func (r *GroupRepositoryBase) UpdateOneByName(groupName string, patch *GroupPatch) (*GroupEntity, error) {
	update := pqcomp.New(1, 7)
	update.AddArg(groupName)
	if patch.CreatedAt != nil {
		update.AddExpr(TableGroupColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableGroupColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableGroupColumnDescription, pqcomp.Equal, patch.Description)
	update.AddExpr(TableGroupColumnName, pqcomp.Equal, patch.Name)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableGroupColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableGroupColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableGroupColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)

	if update.Len() == 0 {
		return nil, errors.New("Group update failure, nothing to update")
	}
	query := "UPDATE charon.group SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE name = $1 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneByName"); err != nil {
			return nil, err
		}
	}
	var e GroupEntity
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

func (r *GroupRepositoryBase) DeleteOneByID(id int64) (int64, error) {
	query := "DELETE FROM charon.group WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	TablePermission                                      = "charon.permission"
	TablePermissionColumnAction                          = "action"
	TablePermissionColumnCreatedAt                       = "created_at"
	TablePermissionColumnID                              = "id"
	TablePermissionColumnModule                          = "module"
	TablePermissionColumnSubsystem                       = "subsystem"
	TablePermissionColumnUpdatedAt                       = "updated_at"
	TablePermissionConstraintPrimaryKey                  = "charon.permission_id_pkey"
	TablePermissionConstraintSubsystemModuleActionUnique = "charon.permission_subsystem_module_action_key"
)

var (
	TablePermissionColumns = []string{
		TablePermissionColumnAction,
		TablePermissionColumnCreatedAt,
		TablePermissionColumnID,
		TablePermissionColumnModule,
		TablePermissionColumnSubsystem,
		TablePermissionColumnUpdatedAt,
	}
)

type PermissionEntity struct {
	// Action ...
	Action string
	// CreatedAt ...
	CreatedAt time.Time
	// ID ...
	ID int64
	// Module ...
	Module string
	// Subsystem ...
	Subsystem string
	// UpdatedAt ...
	UpdatedAt *time.Time
	// Groups ...
	Groups []*GroupEntity
	// Users ...
	Users []*UserEntity
}

func (e *PermissionEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TablePermissionColumnAction:
		return &e.Action, true
	case TablePermissionColumnCreatedAt:
		return &e.CreatedAt, true
	case TablePermissionColumnID:
		return &e.ID, true
	case TablePermissionColumnModule:
		return &e.Module, true
	case TablePermissionColumnSubsystem:
		return &e.Subsystem, true
	case TablePermissionColumnUpdatedAt:
		return &e.UpdatedAt, true
	default:
		return nil, false
	}
}
func (e *PermissionEntity) Props(cns ...string) ([]interface{}, error) {

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

// PermissionIterator is not thread safe.
type PermissionIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *PermissionIterator) Next() bool {
	return i.rows.Next()
}

func (i *PermissionIterator) Close() error {
	return i.rows.Close()
}

func (i *PermissionIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *PermissionIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around Permission method that makes iterator more generic.
func (i *PermissionIterator) Ent() (interface{}, error) {
	return i.Permission()
}

func (i *PermissionIterator) Permission() (*PermissionEntity, error) {
	var ent PermissionEntity
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

type PermissionCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	Action        *qtypes.String
	CreatedAt     *qtypes.Timestamp
	ID            *qtypes.Int64
	Module        *qtypes.String
	Subsystem     *qtypes.String
	UpdatedAt     *qtypes.Timestamp
}

func (c *PermissionCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if err = pqtgo.WriteCompositionQueryString(c.Action, TablePermissionColumnAction, com, pqtgo.And); err != nil {
		return
	}

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TablePermissionColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TablePermissionColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TablePermissionColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.ID, TablePermissionColumnID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.Module, TablePermissionColumnModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.Subsystem, TablePermissionColumnSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TablePermissionColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TablePermissionColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TablePermissionColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TablePermissionColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TablePermissionColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type PermissionPatch struct {
	Action    *ntypes.String
	CreatedAt *time.Time
	Module    *ntypes.String
	Subsystem *ntypes.String
	UpdatedAt *time.Time
}

type PermissionRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanPermissionRows(rows *sql.Rows) ([]*PermissionEntity, error) {
	var (
		entities []*PermissionEntity
		err      error
	)
	for rows.Next() {
		var ent PermissionEntity
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

func (r *PermissionRepositoryBase) Count(c *PermissionCriteria) (int64, error) {

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

func (r *PermissionRepositoryBase) Find(c *PermissionCriteria) ([]*PermissionEntity, error) {

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
func (r *PermissionRepositoryBase) FindIter(c *PermissionCriteria) (*PermissionIterator, error) {

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

	return &PermissionIterator{rows: rows}, nil
}
func (r *PermissionRepositoryBase) FindOneByID(id int64) (*PermissionEntity, error) {
	var (
		ent PermissionEntity
	)
	query := `SELECT action,
created_at,
id,
module,
subsystem,
updated_at
 FROM charon.permission WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
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

	return &ent, nil
}
func (r *PermissionRepositoryBase) FindOneBySubsystemAndModuleAndAction(permissionSubsystem string, permissionModule string, permissionAction string) (*PermissionEntity, error) {
	var (
		ent PermissionEntity
	)
	query := `SELECT action, created_at, id, module, subsystem, updated_at FROM charon.permission WHERE subsystem = $1 AND module = $2 AND action = $3`
	err := r.db.QueryRow(query, permissionSubsystem, permissionModule, permissionAction).Scan(
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

	return &ent, nil
}
func (r *PermissionRepositoryBase) Insert(e *PermissionEntity) (*PermissionEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(TablePermissionColumnAction, "", e.Action)
	insert.AddExpr(TablePermissionColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TablePermissionColumnModule, "", e.Module)
	insert.AddExpr(TablePermissionColumnSubsystem, "", e.Subsystem)
	insert.AddExpr(TablePermissionColumnUpdatedAt, "", e.UpdatedAt)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *PermissionRepositoryBase) Upsert(e *PermissionEntity, p *PermissionPatch, inf ...string) (*PermissionEntity, error) {
	insert := pqcomp.New(0, 6)
	update := insert.Compose(6)
	insert.AddExpr(TablePermissionColumnAction, "", e.Action)
	insert.AddExpr(TablePermissionColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TablePermissionColumnModule, "", e.Module)
	insert.AddExpr(TablePermissionColumnSubsystem, "", e.Subsystem)
	insert.AddExpr(TablePermissionColumnUpdatedAt, "", e.UpdatedAt)
	if len(inf) > 0 {
		update.AddExpr(TablePermissionColumnAction, "=", p.Action)
		update.AddExpr(TablePermissionColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TablePermissionColumnModule, "=", p.Module)
		update.AddExpr(TablePermissionColumnSubsystem, "=", p.Subsystem)
		update.AddExpr(TablePermissionColumnUpdatedAt, "=", p.UpdatedAt)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *PermissionRepositoryBase) UpdateOneByID(id int64, patch *PermissionPatch) (*PermissionEntity, error) {
	update := pqcomp.New(1, 6)
	update.AddArg(id)

	update.AddExpr(TablePermissionColumnAction, pqcomp.Equal, patch.Action)
	if patch.CreatedAt != nil {
		update.AddExpr(TablePermissionColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TablePermissionColumnModule, pqcomp.Equal, patch.Module)
	update.AddExpr(TablePermissionColumnSubsystem, pqcomp.Equal, patch.Subsystem)
	if patch.UpdatedAt != nil {
		update.AddExpr(TablePermissionColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TablePermissionColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}

	if update.Len() == 0 {
		return nil, errors.New("Permission update failure, nothing to update")
	}
	query := "UPDATE charon.permission SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE id = $1 RETURNING " + strings.Join(r.columns, ", ")
	var e PermissionEntity
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
func (r *PermissionRepositoryBase) UpdateOneBySubsystemAndModuleAndAction(permissionSubsystem string, permissionModule string, permissionAction string, patch *PermissionPatch) (*PermissionEntity, error) {
	update := pqcomp.New(3, 6)
	update.AddArg(permissionSubsystem)
	update.AddArg(permissionModule)
	update.AddArg(permissionAction)
	update.AddExpr(TablePermissionColumnAction, pqcomp.Equal, patch.Action)
	if patch.CreatedAt != nil {
		update.AddExpr(TablePermissionColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TablePermissionColumnModule, pqcomp.Equal, patch.Module)
	update.AddExpr(TablePermissionColumnSubsystem, pqcomp.Equal, patch.Subsystem)
	if patch.UpdatedAt != nil {
		update.AddExpr(TablePermissionColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TablePermissionColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}

	if update.Len() == 0 {
		return nil, errors.New("Permission update failure, nothing to update")
	}
	query := "UPDATE charon.permission SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE subsystem = $1 AND module = $2 AND action = $3 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneBySubsystemAndModuleAndAction"); err != nil {
			return nil, err
		}
	}
	var e PermissionEntity
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

func (r *PermissionRepositoryBase) DeleteOneByID(id int64) (int64, error) {
	query := "DELETE FROM charon.permission WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	TableUserGroups                              = "charon.user_groups"
	TableUserGroupsColumnCreatedAt               = "created_at"
	TableUserGroupsColumnCreatedBy               = "created_by"
	TableUserGroupsColumnGroupID                 = "group_id"
	TableUserGroupsColumnUpdatedAt               = "updated_at"
	TableUserGroupsColumnUpdatedBy               = "updated_by"
	TableUserGroupsColumnUserID                  = "user_id"
	TableUserGroupsConstraintCreatedByForeignKey = "charon.user_groups_created_by_fkey"
	TableUserGroupsConstraintUpdatedByForeignKey = "charon.user_groups_updated_by_fkey"
	TableUserGroupsConstraintUserIDForeignKey    = "charon.user_groups_user_id_fkey"
	TableUserGroupsConstraintGroupIDForeignKey   = "charon.user_groups_group_id_fkey"
	TableUserGroupsConstraintUserIDGroupIDUnique = "charon.user_groups_user_id_group_id_key"
)

var (
	TableUserGroupsColumns = []string{
		TableUserGroupsColumnCreatedAt,
		TableUserGroupsColumnCreatedBy,
		TableUserGroupsColumnGroupID,
		TableUserGroupsColumnUpdatedAt,
		TableUserGroupsColumnUpdatedBy,
		TableUserGroupsColumnUserID,
	}
)

type UserGroupsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy *ntypes.Int64
	// GroupID ...
	GroupID int64
	// UpdatedAt ...
	UpdatedAt *time.Time
	// UpdatedBy ...
	UpdatedBy *ntypes.Int64
	// UserID ...
	UserID int64
	// User ...
	User *UserEntity
	// Group ...
	Group *GroupEntity
	// Author ...
	Author *UserEntity
	// Modifier ...
	Modifier *UserEntity
}

func (e *UserGroupsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TableUserGroupsColumnCreatedAt:
		return &e.CreatedAt, true
	case TableUserGroupsColumnCreatedBy:
		return &e.CreatedBy, true
	case TableUserGroupsColumnGroupID:
		return &e.GroupID, true
	case TableUserGroupsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableUserGroupsColumnUpdatedBy:
		return &e.UpdatedBy, true
	case TableUserGroupsColumnUserID:
		return &e.UserID, true
	default:
		return nil, false
	}
}
func (e *UserGroupsEntity) Props(cns ...string) ([]interface{}, error) {

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

// UserGroupsIterator is not thread safe.
type UserGroupsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *UserGroupsIterator) Next() bool {
	return i.rows.Next()
}

func (i *UserGroupsIterator) Close() error {
	return i.rows.Close()
}

func (i *UserGroupsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *UserGroupsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around UserGroups method that makes iterator more generic.
func (i *UserGroupsIterator) Ent() (interface{}, error) {
	return i.UserGroups()
}

func (i *UserGroupsIterator) UserGroups() (*UserGroupsEntity, error) {
	var ent UserGroupsEntity
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

type UserGroupsCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	CreatedAt     *qtypes.Timestamp
	CreatedBy     *qtypes.Int64
	GroupID       *qtypes.Int64
	UpdatedAt     *qtypes.Timestamp
	UpdatedBy     *qtypes.Int64
	UserID        *qtypes.Int64
}

func (c *UserGroupsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserGroupsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserGroupsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserGroupsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.CreatedBy, TableUserGroupsColumnCreatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.GroupID, TableUserGroupsColumnGroupID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserGroupsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserGroupsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserGroupsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserGroupsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UpdatedBy, TableUserGroupsColumnUpdatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UserID, TableUserGroupsColumnUserID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TableUserGroupsColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type UserGroupsPatch struct {
	CreatedAt *time.Time
	CreatedBy *ntypes.Int64
	GroupID   *ntypes.Int64
	UpdatedAt *time.Time
	UpdatedBy *ntypes.Int64
	UserID    *ntypes.Int64
}

type UserGroupsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserGroupsRows(rows *sql.Rows) ([]*UserGroupsEntity, error) {
	var (
		entities []*UserGroupsEntity
		err      error
	)
	for rows.Next() {
		var ent UserGroupsEntity
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

func (r *UserGroupsRepositoryBase) Count(c *UserGroupsCriteria) (int64, error) {

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

func (r *UserGroupsRepositoryBase) Find(c *UserGroupsCriteria) ([]*UserGroupsEntity, error) {

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
func (r *UserGroupsRepositoryBase) FindIter(c *UserGroupsCriteria) (*UserGroupsIterator, error) {

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

	return &UserGroupsIterator{rows: rows}, nil
}
func (r *UserGroupsRepositoryBase) FindOneByUserIDAndGroupID(userGroupsUserID int64, userGroupsGroupID int64) (*UserGroupsEntity, error) {
	var (
		ent UserGroupsEntity
	)
	query := `SELECT created_at, created_by, group_id, updated_at, updated_by, user_id FROM charon.user_groups WHERE user_id = $1 AND group_id = $2`
	err := r.db.QueryRow(query, userGroupsUserID, userGroupsGroupID).Scan(
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

	return &ent, nil
}
func (r *UserGroupsRepositoryBase) Insert(e *UserGroupsEntity) (*UserGroupsEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(TableUserGroupsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserGroupsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserGroupsColumnGroupID, "", e.GroupID)
	insert.AddExpr(TableUserGroupsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserGroupsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserGroupsColumnUserID, "", e.UserID)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *UserGroupsRepositoryBase) Upsert(e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (*UserGroupsEntity, error) {
	insert := pqcomp.New(0, 6)
	update := insert.Compose(6)
	insert.AddExpr(TableUserGroupsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserGroupsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserGroupsColumnGroupID, "", e.GroupID)
	insert.AddExpr(TableUserGroupsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserGroupsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserGroupsColumnUserID, "", e.UserID)
	if len(inf) > 0 {
		update.AddExpr(TableUserGroupsColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TableUserGroupsColumnCreatedBy, "=", p.CreatedBy)
		update.AddExpr(TableUserGroupsColumnGroupID, "=", p.GroupID)
		update.AddExpr(TableUserGroupsColumnUpdatedAt, "=", p.UpdatedAt)
		update.AddExpr(TableUserGroupsColumnUpdatedBy, "=", p.UpdatedBy)
		update.AddExpr(TableUserGroupsColumnUserID, "=", p.UserID)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *UserGroupsRepositoryBase) UpdateOneByUserIDAndGroupID(userGroupsUserID int64, userGroupsGroupID int64, patch *UserGroupsPatch) (*UserGroupsEntity, error) {
	update := pqcomp.New(2, 6)
	update.AddArg(userGroupsUserID)
	update.AddArg(userGroupsGroupID)
	if patch.CreatedAt != nil {
		update.AddExpr(TableUserGroupsColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableUserGroupsColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableUserGroupsColumnGroupID, pqcomp.Equal, patch.GroupID)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableUserGroupsColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableUserGroupsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableUserGroupsColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)
	update.AddExpr(TableUserGroupsColumnUserID, pqcomp.Equal, patch.UserID)

	if update.Len() == 0 {
		return nil, errors.New("UserGroups update failure, nothing to update")
	}
	query := "UPDATE charon.user_groups SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE user_id = $1 AND group_id = $2 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneByUserIDAndGroupID"); err != nil {
			return nil, err
		}
	}
	var e UserGroupsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
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

	return &e, nil
}

const (
	TableGroupPermissions                                                                           = "charon.group_permissions"
	TableGroupPermissionsColumnCreatedAt                                                            = "created_at"
	TableGroupPermissionsColumnCreatedBy                                                            = "created_by"
	TableGroupPermissionsColumnGroupID                                                              = "group_id"
	TableGroupPermissionsColumnPermissionAction                                                     = "permission_action"
	TableGroupPermissionsColumnPermissionModule                                                     = "permission_module"
	TableGroupPermissionsColumnPermissionSubsystem                                                  = "permission_subsystem"
	TableGroupPermissionsColumnUpdatedAt                                                            = "updated_at"
	TableGroupPermissionsColumnUpdatedBy                                                            = "updated_by"
	TableGroupPermissionsConstraintCreatedByForeignKey                                              = "charon.group_permissions_created_by_fkey"
	TableGroupPermissionsConstraintUpdatedByForeignKey                                              = "charon.group_permissions_updated_by_fkey"
	TableGroupPermissionsConstraintGroupIDForeignKey                                                = "charon.group_permissions_group_id_fkey"
	TableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey    = "charon.group_permissions_subsystem_module_action_fkey"
	TableGroupPermissionsConstraintGroupIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.group_permissions_group_id_subsystem_module_action_key"
)

var (
	TableGroupPermissionsColumns = []string{
		TableGroupPermissionsColumnCreatedAt,
		TableGroupPermissionsColumnCreatedBy,
		TableGroupPermissionsColumnGroupID,
		TableGroupPermissionsColumnPermissionAction,
		TableGroupPermissionsColumnPermissionModule,
		TableGroupPermissionsColumnPermissionSubsystem,
		TableGroupPermissionsColumnUpdatedAt,
		TableGroupPermissionsColumnUpdatedBy,
	}
)

type GroupPermissionsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy *ntypes.Int64
	// GroupID ...
	GroupID int64
	// PermissionAction ...
	PermissionAction string
	// PermissionModule ...
	PermissionModule string
	// PermissionSubsystem ...
	PermissionSubsystem string
	// UpdatedAt ...
	UpdatedAt *time.Time
	// UpdatedBy ...
	UpdatedBy *ntypes.Int64
	// Group ...
	Group *GroupEntity
	// Author ...
	Author *UserEntity
	// Modifier ...
	Modifier *UserEntity
}

func (e *GroupPermissionsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TableGroupPermissionsColumnCreatedAt:
		return &e.CreatedAt, true
	case TableGroupPermissionsColumnCreatedBy:
		return &e.CreatedBy, true
	case TableGroupPermissionsColumnGroupID:
		return &e.GroupID, true
	case TableGroupPermissionsColumnPermissionAction:
		return &e.PermissionAction, true
	case TableGroupPermissionsColumnPermissionModule:
		return &e.PermissionModule, true
	case TableGroupPermissionsColumnPermissionSubsystem:
		return &e.PermissionSubsystem, true
	case TableGroupPermissionsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableGroupPermissionsColumnUpdatedBy:
		return &e.UpdatedBy, true
	default:
		return nil, false
	}
}
func (e *GroupPermissionsEntity) Props(cns ...string) ([]interface{}, error) {

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

// GroupPermissionsIterator is not thread safe.
type GroupPermissionsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *GroupPermissionsIterator) Next() bool {
	return i.rows.Next()
}

func (i *GroupPermissionsIterator) Close() error {
	return i.rows.Close()
}

func (i *GroupPermissionsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *GroupPermissionsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around GroupPermissions method that makes iterator more generic.
func (i *GroupPermissionsIterator) Ent() (interface{}, error) {
	return i.GroupPermissions()
}

func (i *GroupPermissionsIterator) GroupPermissions() (*GroupPermissionsEntity, error) {
	var ent GroupPermissionsEntity
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

type GroupPermissionsCriteria struct {
	Offset, Limit       int64
	Sort                map[string]bool
	CreatedAt           *qtypes.Timestamp
	CreatedBy           *qtypes.Int64
	GroupID             *qtypes.Int64
	PermissionAction    *qtypes.String
	PermissionModule    *qtypes.String
	PermissionSubsystem *qtypes.String
	UpdatedAt           *qtypes.Timestamp
	UpdatedBy           *qtypes.Int64
}

func (c *GroupPermissionsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableGroupPermissionsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableGroupPermissionsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableGroupPermissionsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.CreatedBy, TableGroupPermissionsColumnCreatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.GroupID, TableGroupPermissionsColumnGroupID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionAction, TableGroupPermissionsColumnPermissionAction, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionModule, TableGroupPermissionsColumnPermissionModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionSubsystem, TableGroupPermissionsColumnPermissionSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableGroupPermissionsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableGroupPermissionsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UpdatedBy, TableGroupPermissionsColumnUpdatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TableGroupPermissionsColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type GroupPermissionsPatch struct {
	CreatedAt           *time.Time
	CreatedBy           *ntypes.Int64
	GroupID             *ntypes.Int64
	PermissionAction    *ntypes.String
	PermissionModule    *ntypes.String
	PermissionSubsystem *ntypes.String
	UpdatedAt           *time.Time
	UpdatedBy           *ntypes.Int64
}

type GroupPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanGroupPermissionsRows(rows *sql.Rows) ([]*GroupPermissionsEntity, error) {
	var (
		entities []*GroupPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent GroupPermissionsEntity
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

func (r *GroupPermissionsRepositoryBase) Count(c *GroupPermissionsCriteria) (int64, error) {

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

func (r *GroupPermissionsRepositoryBase) Find(c *GroupPermissionsCriteria) ([]*GroupPermissionsEntity, error) {

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
func (r *GroupPermissionsRepositoryBase) FindIter(c *GroupPermissionsCriteria) (*GroupPermissionsIterator, error) {

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

	return &GroupPermissionsIterator{rows: rows}, nil
}
func (r *GroupPermissionsRepositoryBase) FindOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string) (*GroupPermissionsEntity, error) {
	var (
		ent GroupPermissionsEntity
	)
	query := `SELECT created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by FROM charon.group_permissions WHERE group_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4`
	err := r.db.QueryRow(query, groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction).Scan(
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

	return &ent, nil
}
func (r *GroupPermissionsRepositoryBase) Insert(e *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(TableGroupPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableGroupPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableGroupPermissionsColumnGroupID, "", e.GroupID)
	insert.AddExpr(TableGroupPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(TableGroupPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(TableGroupPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(TableGroupPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableGroupPermissionsColumnUpdatedBy, "", e.UpdatedBy)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *GroupPermissionsRepositoryBase) Upsert(e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (*GroupPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	update := insert.Compose(8)
	insert.AddExpr(TableGroupPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableGroupPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableGroupPermissionsColumnGroupID, "", e.GroupID)
	insert.AddExpr(TableGroupPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(TableGroupPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(TableGroupPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(TableGroupPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableGroupPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	if len(inf) > 0 {
		update.AddExpr(TableGroupPermissionsColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TableGroupPermissionsColumnCreatedBy, "=", p.CreatedBy)
		update.AddExpr(TableGroupPermissionsColumnGroupID, "=", p.GroupID)
		update.AddExpr(TableGroupPermissionsColumnPermissionAction, "=", p.PermissionAction)
		update.AddExpr(TableGroupPermissionsColumnPermissionModule, "=", p.PermissionModule)
		update.AddExpr(TableGroupPermissionsColumnPermissionSubsystem, "=", p.PermissionSubsystem)
		update.AddExpr(TableGroupPermissionsColumnUpdatedAt, "=", p.UpdatedAt)
		update.AddExpr(TableGroupPermissionsColumnUpdatedBy, "=", p.UpdatedBy)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *GroupPermissionsRepositoryBase) UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, patch *GroupPermissionsPatch) (*GroupPermissionsEntity, error) {
	update := pqcomp.New(4, 8)
	update.AddArg(groupPermissionsGroupID)
	update.AddArg(groupPermissionsPermissionSubsystem)
	update.AddArg(groupPermissionsPermissionModule)
	update.AddArg(groupPermissionsPermissionAction)
	if patch.CreatedAt != nil {
		update.AddExpr(TableGroupPermissionsColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableGroupPermissionsColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableGroupPermissionsColumnGroupID, pqcomp.Equal, patch.GroupID)
	update.AddExpr(TableGroupPermissionsColumnPermissionAction, pqcomp.Equal, patch.PermissionAction)
	update.AddExpr(TableGroupPermissionsColumnPermissionModule, pqcomp.Equal, patch.PermissionModule)
	update.AddExpr(TableGroupPermissionsColumnPermissionSubsystem, pqcomp.Equal, patch.PermissionSubsystem)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableGroupPermissionsColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableGroupPermissionsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableGroupPermissionsColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)

	if update.Len() == 0 {
		return nil, errors.New("GroupPermissions update failure, nothing to update")
	}
	query := "UPDATE charon.group_permissions SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE group_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction"); err != nil {
			return nil, err
		}
	}
	var e GroupPermissionsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
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

	return &e, nil
}

const (
	TableUserPermissions                                                                          = "charon.user_permissions"
	TableUserPermissionsColumnCreatedAt                                                           = "created_at"
	TableUserPermissionsColumnCreatedBy                                                           = "created_by"
	TableUserPermissionsColumnPermissionAction                                                    = "permission_action"
	TableUserPermissionsColumnPermissionModule                                                    = "permission_module"
	TableUserPermissionsColumnPermissionSubsystem                                                 = "permission_subsystem"
	TableUserPermissionsColumnUpdatedAt                                                           = "updated_at"
	TableUserPermissionsColumnUpdatedBy                                                           = "updated_by"
	TableUserPermissionsColumnUserID                                                              = "user_id"
	TableUserPermissionsConstraintCreatedByForeignKey                                             = "charon.user_permissions_created_by_fkey"
	TableUserPermissionsConstraintUpdatedByForeignKey                                             = "charon.user_permissions_updated_by_fkey"
	TableUserPermissionsConstraintUserIDForeignKey                                                = "charon.user_permissions_user_id_fkey"
	TableUserPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey   = "charon.user_permissions_subsystem_module_action_fkey"
	TableUserPermissionsConstraintUserIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.user_permissions_user_id_subsystem_module_action_key"
)

var (
	TableUserPermissionsColumns = []string{
		TableUserPermissionsColumnCreatedAt,
		TableUserPermissionsColumnCreatedBy,
		TableUserPermissionsColumnPermissionAction,
		TableUserPermissionsColumnPermissionModule,
		TableUserPermissionsColumnPermissionSubsystem,
		TableUserPermissionsColumnUpdatedAt,
		TableUserPermissionsColumnUpdatedBy,
		TableUserPermissionsColumnUserID,
	}
)

type UserPermissionsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy *ntypes.Int64
	// PermissionAction ...
	PermissionAction string
	// PermissionModule ...
	PermissionModule string
	// PermissionSubsystem ...
	PermissionSubsystem string
	// UpdatedAt ...
	UpdatedAt *time.Time
	// UpdatedBy ...
	UpdatedBy *ntypes.Int64
	// UserID ...
	UserID int64
	// User ...
	User *UserEntity
	// Author ...
	Author *UserEntity
	// Modifier ...
	Modifier *UserEntity
}

func (e *UserPermissionsEntity) Prop(cn string) (interface{}, bool) {
	switch cn {
	case TableUserPermissionsColumnCreatedAt:
		return &e.CreatedAt, true
	case TableUserPermissionsColumnCreatedBy:
		return &e.CreatedBy, true
	case TableUserPermissionsColumnPermissionAction:
		return &e.PermissionAction, true
	case TableUserPermissionsColumnPermissionModule:
		return &e.PermissionModule, true
	case TableUserPermissionsColumnPermissionSubsystem:
		return &e.PermissionSubsystem, true
	case TableUserPermissionsColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableUserPermissionsColumnUpdatedBy:
		return &e.UpdatedBy, true
	case TableUserPermissionsColumnUserID:
		return &e.UserID, true
	default:
		return nil, false
	}
}
func (e *UserPermissionsEntity) Props(cns ...string) ([]interface{}, error) {

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

// UserPermissionsIterator is not thread safe.
type UserPermissionsIterator struct {
	rows *sql.Rows
	cols []string
}

func (i *UserPermissionsIterator) Next() bool {
	return i.rows.Next()
}

func (i *UserPermissionsIterator) Close() error {
	return i.rows.Close()
}

func (i *UserPermissionsIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache outpu inside iterator.
func (i *UserPermissionsIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around UserPermissions method that makes iterator more generic.
func (i *UserPermissionsIterator) Ent() (interface{}, error) {
	return i.UserPermissions()
}

func (i *UserPermissionsIterator) UserPermissions() (*UserPermissionsEntity, error) {
	var ent UserPermissionsEntity
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

type UserPermissionsCriteria struct {
	Offset, Limit       int64
	Sort                map[string]bool
	CreatedAt           *qtypes.Timestamp
	CreatedBy           *qtypes.Int64
	PermissionAction    *qtypes.String
	PermissionModule    *qtypes.String
	PermissionSubsystem *qtypes.String
	UpdatedAt           *qtypes.Timestamp
	UpdatedBy           *qtypes.Int64
	UserID              *qtypes.Int64
}

func (c *UserPermissionsCriteria) WriteComposition(sel string, com *pqtgo.Composer, opt *pqtgo.CompositionOpts) (err error) {

	if c.CreatedAt != nil && c.CreatedAt.Valid {
		CreatedAtt1 := c.CreatedAt.Value()
		if CreatedAtt1 != nil {
			CreatedAt1, err := ptypes.Timestamp(CreatedAtt1)
			if err != nil {
				return err
			}
			switch c.CreatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				if c.CreatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnCreatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.CreatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.CreatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserPermissionsColumnCreatedAt)
					com.WriteString(" IN (")
					for i, v := range c.CreatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				CreatedAtt2 := c.CreatedAt.Values[1]
				if CreatedAtt2 != nil {
					CreatedAt2, err := ptypes.Timestamp(CreatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserPermissionsColumnCreatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(CreatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserPermissionsColumnCreatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(CreatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.CreatedBy, TableUserPermissionsColumnCreatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionAction, TableUserPermissionsColumnPermissionAction, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionModule, TableUserPermissionsColumnPermissionModule, com, pqtgo.And); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryString(c.PermissionSubsystem, TableUserPermissionsColumnPermissionSubsystem, com, pqtgo.And); err != nil {
		return
	}

	if c.UpdatedAt != nil && c.UpdatedAt.Valid {
		UpdatedAtt1 := c.UpdatedAt.Value()
		if UpdatedAtt1 != nil {
			UpdatedAt1, err := ptypes.Timestamp(UpdatedAtt1)
			if err != nil {
				return err
			}
			switch c.UpdatedAt.Type {
			case qtypes.QueryType_NULL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" IS NOT NULL ")
				} else {
					com.WriteString(" IS NULL ")
				}
			case qtypes.QueryType_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				if c.UpdatedAt.Negation {
					com.WriteString(" <> ")
				} else {
					com.WriteString(" = ")
				}
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				com.WriteString(">")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_GREATER_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				com.WriteString(">=")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				com.WriteString(" < ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_LESS_EQUAL:
				if com.Dirty {
					com.WriteString(" AND ")
				}
				com.Dirty = true

				com.WriteString(TableUserPermissionsColumnUpdatedAt)
				com.WriteString(" <= ")
				com.WritePlaceholder()
				com.Add(c.UpdatedAt.Value())
			case qtypes.QueryType_IN:
				if len(c.UpdatedAt.Values) > 0 {
					if com.Dirty {
						com.WriteString(" AND ")
					}
					com.Dirty = true

					com.WriteString(TableUserPermissionsColumnUpdatedAt)
					com.WriteString(" IN (")
					for i, v := range c.UpdatedAt.Values {
						if i != 0 {
							com.WriteString(", ")
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

				UpdatedAtt2 := c.UpdatedAt.Values[1]
				if UpdatedAtt2 != nil {
					UpdatedAt2, err := ptypes.Timestamp(UpdatedAtt2)
					if err != nil {
						return err
					}
					com.WriteString(TableUserPermissionsColumnUpdatedAt)
					com.WriteString(" > ")
					com.WritePlaceholder()
					com.Add(UpdatedAt1)
					com.WriteString(" AND ")
					com.WriteString(TableUserPermissionsColumnUpdatedAt)
					com.WriteString(" < ")
					com.WritePlaceholder()
					com.Add(UpdatedAt2)
				}
			}
		}
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UpdatedBy, TableUserPermissionsColumnUpdatedBy, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if err = pqtgo.WriteCompositionQueryInt64(c.UserID, TableUserPermissionsColumnUserID, com, &pqtgo.CompositionOpts{
		Joint:  " AND ",
		IsJSON: false,
	}); err != nil {
		return
	}

	if len(c.Sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")

		for cn, asc := range c.Sort {
			for _, tcn := range TableUserPermissionsColumns {
				if cn == tcn {
					if i > 0 {
						com.WriteString(", ")
					}
					com.WriteString(cn)
					if !asc {
						com.WriteString(" DESC ")
					}
					i++
					break
				}
			}
		}
	}
	if c.Offset > 0 {
		if _, err = com.WriteString(" OFFSET "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Offset)
	}
	if c.Limit > 0 {
		if _, err = com.WriteString(" LIMIT "); err != nil {
			return
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		if _, err = com.WriteString(" "); err != nil {
			return
		}
		com.Add(c.Limit)
	}

	return
}

type UserPermissionsPatch struct {
	CreatedAt           *time.Time
	CreatedBy           *ntypes.Int64
	PermissionAction    *ntypes.String
	PermissionModule    *ntypes.String
	PermissionSubsystem *ntypes.String
	UpdatedAt           *time.Time
	UpdatedBy           *ntypes.Int64
	UserID              *ntypes.Int64
}

type UserPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func ScanUserPermissionsRows(rows *sql.Rows) ([]*UserPermissionsEntity, error) {
	var (
		entities []*UserPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent UserPermissionsEntity
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

func (r *UserPermissionsRepositoryBase) Count(c *UserPermissionsCriteria) (int64, error) {

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

func (r *UserPermissionsRepositoryBase) Find(c *UserPermissionsCriteria) ([]*UserPermissionsEntity, error) {

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
func (r *UserPermissionsRepositoryBase) FindIter(c *UserPermissionsCriteria) (*UserPermissionsIterator, error) {

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

	return &UserPermissionsIterator{rows: rows}, nil
}
func (r *UserPermissionsRepositoryBase) FindOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string) (*UserPermissionsEntity, error) {
	var (
		ent UserPermissionsEntity
	)
	query := `SELECT created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id FROM charon.user_permissions WHERE user_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4`
	err := r.db.QueryRow(query, userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction).Scan(
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

	return &ent, nil
}
func (r *UserPermissionsRepositoryBase) Insert(e *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(TableUserPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(TableUserPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(TableUserPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(TableUserPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserPermissionsColumnUserID, "", e.UserID)

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
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Insert"); err != nil {
			return nil, err
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
func (r *UserPermissionsRepositoryBase) Upsert(e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (*UserPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	update := insert.Compose(8)
	insert.AddExpr(TableUserPermissionsColumnCreatedAt, "", e.CreatedAt)
	insert.AddExpr(TableUserPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(TableUserPermissionsColumnPermissionAction, "", e.PermissionAction)
	insert.AddExpr(TableUserPermissionsColumnPermissionModule, "", e.PermissionModule)
	insert.AddExpr(TableUserPermissionsColumnPermissionSubsystem, "", e.PermissionSubsystem)
	insert.AddExpr(TableUserPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(TableUserPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(TableUserPermissionsColumnUserID, "", e.UserID)
	if len(inf) > 0 {
		update.AddExpr(TableUserPermissionsColumnCreatedAt, "=", p.CreatedAt)
		update.AddExpr(TableUserPermissionsColumnCreatedBy, "=", p.CreatedBy)
		update.AddExpr(TableUserPermissionsColumnPermissionAction, "=", p.PermissionAction)
		update.AddExpr(TableUserPermissionsColumnPermissionModule, "=", p.PermissionModule)
		update.AddExpr(TableUserPermissionsColumnPermissionSubsystem, "=", p.PermissionSubsystem)
		update.AddExpr(TableUserPermissionsColumnUpdatedAt, "=", p.UpdatedAt)
		update.AddExpr(TableUserPermissionsColumnUpdatedBy, "=", p.UpdatedBy)
		update.AddExpr(TableUserPermissionsColumnUserID, "=", p.UserID)
	}

	b := bytes.NewBufferString("INSERT INTO " + r.table)

	if insert.Len() > 0 {
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
	}
	b.WriteString(" ON CONFLICT ")
	if len(inf) > 0 && update.Len() > 0 {
		b.WriteString(" (")
		for j, i := range inf {
			if j != 0 {
				b.WriteString(", ")
			}
			b.WriteString(i)
		}
		b.WriteString(") ")
		b.WriteString(" DO UPDATE SET ")
		for update.Next() {
			if !update.First() {
				b.WriteString(", ")
			}

			b.WriteString(update.Key())
			b.WriteString(" ")
			b.WriteString(update.Oper())
			b.WriteString(" ")
			b.WriteString(update.PlaceHolder())
		}
	} else {
		b.WriteString(" DO NOTHING ")
	}
	if insert.Len() > 0 {
		if len(r.columns) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(r.columns, ", "))
		}
	}

	if r.dbg {
		if err := r.log.Log("msg", b.String(), "function", "Upsert"); err != nil {
			return nil, err
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
func (r *UserPermissionsRepositoryBase) UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, patch *UserPermissionsPatch) (*UserPermissionsEntity, error) {
	update := pqcomp.New(4, 8)
	update.AddArg(userPermissionsUserID)
	update.AddArg(userPermissionsPermissionSubsystem)
	update.AddArg(userPermissionsPermissionModule)
	update.AddArg(userPermissionsPermissionAction)
	if patch.CreatedAt != nil {
		update.AddExpr(TableUserPermissionsColumnCreatedAt, pqcomp.Equal, patch.CreatedAt)

	}
	update.AddExpr(TableUserPermissionsColumnCreatedBy, pqcomp.Equal, patch.CreatedBy)
	update.AddExpr(TableUserPermissionsColumnPermissionAction, pqcomp.Equal, patch.PermissionAction)
	update.AddExpr(TableUserPermissionsColumnPermissionModule, pqcomp.Equal, patch.PermissionModule)
	update.AddExpr(TableUserPermissionsColumnPermissionSubsystem, pqcomp.Equal, patch.PermissionSubsystem)
	if patch.UpdatedAt != nil {
		update.AddExpr(TableUserPermissionsColumnUpdatedAt, pqcomp.Equal, patch.UpdatedAt)
	} else {
		update.AddExpr(TableUserPermissionsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(TableUserPermissionsColumnUpdatedBy, pqcomp.Equal, patch.UpdatedBy)
	update.AddExpr(TableUserPermissionsColumnUserID, pqcomp.Equal, patch.UserID)

	if update.Len() == 0 {
		return nil, errors.New("UserPermissions update failure, nothing to update")
	}
	query := "UPDATE charon.user_permissions SET "
	for update.Next() {
		if !update.First() {
			query += ", "
		}

		query += update.Key() + " " + update.Oper() + " " + update.PlaceHolder()
	}
	query += " WHERE user_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction"); err != nil {
			return nil, err
		}
	}
	var e UserPermissionsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
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

	return &e, nil
}

const SQL = `
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
