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
	confirmationToken []byte
	createdAt         time.Time
	createdBy         *ntypes.Int64
	firstName         string
	id                int64
	isActive          bool
	isConfirmed       bool
	isStaff           bool
	isSuperuser       bool
	lastLoginAt       *time.Time
	lastName          string
	password          []byte
	updatedAt         *time.Time
	updatedBy         *ntypes.Int64
	username          string
	author            *userEntity
	modifier          *userEntity
	permission        []*permissionEntity
	group             []*groupEntity
}

func (e *userEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserColumnConfirmationToken:
		return &e.confirmationToken, true
	case tableUserColumnCreatedAt:
		return &e.createdAt, true
	case tableUserColumnCreatedBy:
		return &e.createdBy, true
	case tableUserColumnFirstName:
		return &e.firstName, true
	case tableUserColumnID:
		return &e.id, true
	case tableUserColumnIsActive:
		return &e.isActive, true
	case tableUserColumnIsConfirmed:
		return &e.isConfirmed, true
	case tableUserColumnIsStaff:
		return &e.isStaff, true
	case tableUserColumnIsSuperuser:
		return &e.isSuperuser, true
	case tableUserColumnLastLoginAt:
		return &e.lastLoginAt, true
	case tableUserColumnLastName:
		return &e.lastName, true
	case tableUserColumnPassword:
		return &e.password, true
	case tableUserColumnUpdatedAt:
		return &e.updatedAt, true
	case tableUserColumnUpdatedBy:
		return &e.updatedBy, true
	case tableUserColumnUsername:
		return &e.username, true
	default:
		return nil, false
	}
}
func (e *userEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

func scanUserRows(rows *sql.Rows) ([]*userEntity, error) {
	var (
		entities []*userEntity
		err      error
	)
	for rows.Next() {
		var ent userEntity
		err = rows.Scan(
			&ent.confirmationToken,
			&ent.createdAt,
			&ent.createdBy,
			&ent.firstName,
			&ent.id,
			&ent.isActive,
			&ent.isConfirmed,
			&ent.isStaff,
			&ent.isSuperuser,
			&ent.lastLoginAt,
			&ent.lastName,
			&ent.password,
			&ent.updatedAt,
			&ent.updatedBy,
			&ent.username,
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

func (r *userRepositoryBase) count(c *userCriteria) (int64, error) {

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

func (r *userRepositoryBase) find(c *userCriteria) ([]*userEntity, error) {

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

	return scanUserRows(rows)
}
func (r *userRepositoryBase) findIter(c *userCriteria) (*userIterator, error) {

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
func (r *userRepositoryBase) findOneByID(id int64) (*userEntity, error) {
	var (
		entity userEntity
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
		&entity.confirmationToken,
		&entity.createdAt,
		&entity.createdBy,
		&entity.firstName,
		&entity.id,
		&entity.isActive,
		&entity.isConfirmed,
		&entity.isStaff,
		&entity.isSuperuser,
		&entity.lastLoginAt,
		&entity.lastName,
		&entity.password,
		&entity.updatedAt,
		&entity.updatedBy,
		&entity.username,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *userRepositoryBase) insert(e *userEntity) (*userEntity, error) {
	insert := pqcomp.New(0, 15)
	insert.AddExpr(tableUserColumnConfirmationToken, "", e.confirmationToken)
	insert.AddExpr(tableUserColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserColumnFirstName, "", e.firstName)
	insert.AddExpr(tableUserColumnIsActive, "", e.isActive)
	insert.AddExpr(tableUserColumnIsConfirmed, "", e.isConfirmed)
	insert.AddExpr(tableUserColumnIsStaff, "", e.isStaff)
	insert.AddExpr(tableUserColumnIsSuperuser, "", e.isSuperuser)
	insert.AddExpr(tableUserColumnLastLoginAt, "", e.lastLoginAt)
	insert.AddExpr(tableUserColumnLastName, "", e.lastName)
	insert.AddExpr(tableUserColumnPassword, "", e.password)
	insert.AddExpr(tableUserColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserColumnUsername, "", e.username)

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
		&e.confirmationToken,
		&e.createdAt,
		&e.createdBy,
		&e.firstName,
		&e.id,
		&e.isActive,
		&e.isConfirmed,
		&e.isStaff,
		&e.isSuperuser,
		&e.lastLoginAt,
		&e.lastName,
		&e.password,
		&e.updatedAt,
		&e.updatedBy,
		&e.username,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userRepositoryBase) upsert(e *userEntity, p *userPatch, inf ...string) (*userEntity, error) {
	insert := pqcomp.New(0, 15)
	update := insert.Compose(15)
	insert.AddExpr(tableUserColumnConfirmationToken, "", e.confirmationToken)
	insert.AddExpr(tableUserColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserColumnFirstName, "", e.firstName)
	insert.AddExpr(tableUserColumnIsActive, "", e.isActive)
	insert.AddExpr(tableUserColumnIsConfirmed, "", e.isConfirmed)
	insert.AddExpr(tableUserColumnIsStaff, "", e.isStaff)
	insert.AddExpr(tableUserColumnIsSuperuser, "", e.isSuperuser)
	insert.AddExpr(tableUserColumnLastLoginAt, "", e.lastLoginAt)
	insert.AddExpr(tableUserColumnLastName, "", e.lastName)
	insert.AddExpr(tableUserColumnPassword, "", e.password)
	insert.AddExpr(tableUserColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserColumnUsername, "", e.username)
	if len(inf) > 0 {
		update.AddExpr(tableUserColumnConfirmationToken, "=", p.confirmationToken)
		update.AddExpr(tableUserColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tableUserColumnCreatedBy, "=", p.createdBy)
		update.AddExpr(tableUserColumnFirstName, "=", p.firstName)
		update.AddExpr(tableUserColumnIsActive, "=", p.isActive)
		update.AddExpr(tableUserColumnIsConfirmed, "=", p.isConfirmed)
		update.AddExpr(tableUserColumnIsStaff, "=", p.isStaff)
		update.AddExpr(tableUserColumnIsSuperuser, "=", p.isSuperuser)
		update.AddExpr(tableUserColumnLastLoginAt, "=", p.lastLoginAt)
		update.AddExpr(tableUserColumnLastName, "=", p.lastName)
		update.AddExpr(tableUserColumnPassword, "=", p.password)
		update.AddExpr(tableUserColumnUpdatedAt, "=", p.updatedAt)
		update.AddExpr(tableUserColumnUpdatedBy, "=", p.updatedBy)
		update.AddExpr(tableUserColumnUsername, "=", p.username)
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
		&e.confirmationToken,
		&e.createdAt,
		&e.createdBy,
		&e.firstName,
		&e.id,
		&e.isActive,
		&e.isConfirmed,
		&e.isStaff,
		&e.isSuperuser,
		&e.lastLoginAt,
		&e.lastName,
		&e.password,
		&e.updatedAt,
		&e.updatedBy,
		&e.username,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userRepositoryBase) updateOneByID(id int64, patch *userPatch) (*userEntity, error) {
	update := pqcomp.New(1, 15)
	update.AddArg(id)

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
		&e.confirmationToken,
		&e.createdAt,
		&e.createdBy,
		&e.firstName,
		&e.id,
		&e.isActive,
		&e.isConfirmed,
		&e.isStaff,
		&e.isSuperuser,
		&e.lastLoginAt,
		&e.lastName,
		&e.password,
		&e.updatedAt,
		&e.updatedBy,
		&e.username,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *userRepositoryBase) deleteOneByID(id int64) (int64, error) {
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
	createdAt   time.Time
	createdBy   *ntypes.Int64
	description *ntypes.String
	id          int64
	name        string
	updatedAt   *time.Time
	updatedBy   *ntypes.Int64
	author      *userEntity
	modifier    *userEntity
	permission  []*permissionEntity
	user        []*userEntity
}

func (e *groupEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tableGroupColumnCreatedAt:
		return &e.createdAt, true
	case tableGroupColumnCreatedBy:
		return &e.createdBy, true
	case tableGroupColumnDescription:
		return &e.description, true
	case tableGroupColumnID:
		return &e.id, true
	case tableGroupColumnName:
		return &e.name, true
	case tableGroupColumnUpdatedAt:
		return &e.updatedAt, true
	case tableGroupColumnUpdatedBy:
		return &e.updatedBy, true
	default:
		return nil, false
	}
}
func (e *groupEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

func scanGroupRows(rows *sql.Rows) ([]*groupEntity, error) {
	var (
		entities []*groupEntity
		err      error
	)
	for rows.Next() {
		var ent groupEntity
		err = rows.Scan(
			&ent.createdAt,
			&ent.createdBy,
			&ent.description,
			&ent.id,
			&ent.name,
			&ent.updatedAt,
			&ent.updatedBy,
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

func (r *groupRepositoryBase) count(c *groupCriteria) (int64, error) {

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

func (r *groupRepositoryBase) find(c *groupCriteria) ([]*groupEntity, error) {

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

	return scanGroupRows(rows)
}
func (r *groupRepositoryBase) findIter(c *groupCriteria) (*groupIterator, error) {

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
func (r *groupRepositoryBase) findOneByID(id int64) (*groupEntity, error) {
	var (
		entity groupEntity
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
		&entity.createdAt,
		&entity.createdBy,
		&entity.description,
		&entity.id,
		&entity.name,
		&entity.updatedAt,
		&entity.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *groupRepositoryBase) insert(e *groupEntity) (*groupEntity, error) {
	insert := pqcomp.New(0, 7)
	insert.AddExpr(tableGroupColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableGroupColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableGroupColumnDescription, "", e.description)
	insert.AddExpr(tableGroupColumnName, "", e.name)
	insert.AddExpr(tableGroupColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableGroupColumnUpdatedBy, "", e.updatedBy)

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
		&e.createdAt,
		&e.createdBy,
		&e.description,
		&e.id,
		&e.name,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *groupRepositoryBase) upsert(e *groupEntity, p *groupPatch, inf ...string) (*groupEntity, error) {
	insert := pqcomp.New(0, 7)
	update := insert.Compose(7)
	insert.AddExpr(tableGroupColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableGroupColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableGroupColumnDescription, "", e.description)
	insert.AddExpr(tableGroupColumnName, "", e.name)
	insert.AddExpr(tableGroupColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableGroupColumnUpdatedBy, "", e.updatedBy)
	if len(inf) > 0 {
		update.AddExpr(tableGroupColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tableGroupColumnCreatedBy, "=", p.createdBy)
		update.AddExpr(tableGroupColumnDescription, "=", p.description)
		update.AddExpr(tableGroupColumnName, "=", p.name)
		update.AddExpr(tableGroupColumnUpdatedAt, "=", p.updatedAt)
		update.AddExpr(tableGroupColumnUpdatedBy, "=", p.updatedBy)
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
		&e.createdAt,
		&e.createdBy,
		&e.description,
		&e.id,
		&e.name,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *groupRepositoryBase) updateOneByID(id int64, patch *groupPatch) (*groupEntity, error) {
	update := pqcomp.New(1, 7)
	update.AddArg(id)

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
		&e.createdAt,
		&e.createdBy,
		&e.description,
		&e.id,
		&e.name,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *groupRepositoryBase) deleteOneByID(id int64) (int64, error) {
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
	action    string
	createdAt time.Time
	id        int64
	module    string
	subsystem string
	updatedAt *time.Time
	group     []*groupEntity
	user      []*userEntity
}

func (e *permissionEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tablePermissionColumnAction:
		return &e.action, true
	case tablePermissionColumnCreatedAt:
		return &e.createdAt, true
	case tablePermissionColumnID:
		return &e.id, true
	case tablePermissionColumnModule:
		return &e.module, true
	case tablePermissionColumnSubsystem:
		return &e.subsystem, true
	case tablePermissionColumnUpdatedAt:
		return &e.updatedAt, true
	default:
		return nil, false
	}
}
func (e *permissionEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

func scanPermissionRows(rows *sql.Rows) ([]*permissionEntity, error) {
	var (
		entities []*permissionEntity
		err      error
	)
	for rows.Next() {
		var ent permissionEntity
		err = rows.Scan(
			&ent.action,
			&ent.createdAt,
			&ent.id,
			&ent.module,
			&ent.subsystem,
			&ent.updatedAt,
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

func (r *permissionRepositoryBase) count(c *permissionCriteria) (int64, error) {

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

func (r *permissionRepositoryBase) find(c *permissionCriteria) ([]*permissionEntity, error) {

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

	return scanPermissionRows(rows)
}
func (r *permissionRepositoryBase) findIter(c *permissionCriteria) (*permissionIterator, error) {

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
func (r *permissionRepositoryBase) findOneByID(id int64) (*permissionEntity, error) {
	var (
		entity permissionEntity
	)
	query := `SELECT action,
created_at,
id,
module,
subsystem,
updated_at
 FROM charon.permission WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&entity.action,
		&entity.createdAt,
		&entity.id,
		&entity.module,
		&entity.subsystem,
		&entity.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *permissionRepositoryBase) findOneBySubsystemAndModuleAndAction(subsystem string, module string, action string) (*permissionEntity, error) {
	var (
		entity permissionEntity
	)
	query := `SELECT action, created_at, id, module, subsystem, updated_at FROM charon.permission WHERE subsystem = $1 AND module = $2 AND action = $3`
	err := r.db.QueryRow(query, subsystem, module, action).Scan(
		&entity.action,
		&entity.createdAt,
		&entity.id,
		&entity.module,
		&entity.subsystem,
		&entity.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *permissionRepositoryBase) insert(e *permissionEntity) (*permissionEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tablePermissionColumnAction, "", e.action)
	insert.AddExpr(tablePermissionColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tablePermissionColumnModule, "", e.module)
	insert.AddExpr(tablePermissionColumnSubsystem, "", e.subsystem)
	insert.AddExpr(tablePermissionColumnUpdatedAt, "", e.updatedAt)

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
		&e.action,
		&e.createdAt,
		&e.id,
		&e.module,
		&e.subsystem,
		&e.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *permissionRepositoryBase) upsert(e *permissionEntity, p *permissionPatch, inf ...string) (*permissionEntity, error) {
	insert := pqcomp.New(0, 6)
	update := insert.Compose(6)
	insert.AddExpr(tablePermissionColumnAction, "", e.action)
	insert.AddExpr(tablePermissionColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tablePermissionColumnModule, "", e.module)
	insert.AddExpr(tablePermissionColumnSubsystem, "", e.subsystem)
	insert.AddExpr(tablePermissionColumnUpdatedAt, "", e.updatedAt)
	if len(inf) > 0 {
		update.AddExpr(tablePermissionColumnAction, "=", p.action)
		update.AddExpr(tablePermissionColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tablePermissionColumnModule, "=", p.module)
		update.AddExpr(tablePermissionColumnSubsystem, "=", p.subsystem)
		update.AddExpr(tablePermissionColumnUpdatedAt, "=", p.updatedAt)
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
		&e.action,
		&e.createdAt,
		&e.id,
		&e.module,
		&e.subsystem,
		&e.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *permissionRepositoryBase) updateOneByID(id int64, patch *permissionPatch) (*permissionEntity, error) {
	update := pqcomp.New(1, 6)
	update.AddArg(id)

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
		&e.action,
		&e.createdAt,
		&e.id,
		&e.module,
		&e.subsystem,
		&e.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
func (r *permissionRepositoryBase) updateOneBySubsystemAndModuleAndAction(subsystem string, module string, action string, patch *permissionPatch) (*permissionEntity, error) {
	update := pqcomp.New(2, 6)
	update.AddArg(subsystem)
	update.AddArg(module)
	update.AddArg(action)
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
	query += " WHERE subsystem = $1 AND module = $2 AND action = $3 RETURNING " + strings.Join(r.columns, ", ")
	if r.dbg {
		if err := r.log.Log("msg", query, "function", "UpdateOneBySubsystemAndModuleAndAction"); err != nil {
			return nil, err
		}
	}
	var e permissionEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.action,
		&e.createdAt,
		&e.id,
		&e.module,
		&e.subsystem,
		&e.updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *permissionRepositoryBase) deleteOneByID(id int64) (int64, error) {
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
	createdAt time.Time
	createdBy *ntypes.Int64
	groupID   int64
	updatedAt *time.Time
	updatedBy *ntypes.Int64
	userID    int64
	user      *userEntity
	group     *groupEntity
	author    *userEntity
	modifier  *userEntity
}

func (e *userGroupsEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserGroupsColumnCreatedAt:
		return &e.createdAt, true
	case tableUserGroupsColumnCreatedBy:
		return &e.createdBy, true
	case tableUserGroupsColumnGroupID:
		return &e.groupID, true
	case tableUserGroupsColumnUpdatedAt:
		return &e.updatedAt, true
	case tableUserGroupsColumnUpdatedBy:
		return &e.updatedBy, true
	case tableUserGroupsColumnUserID:
		return &e.userID, true
	default:
		return nil, false
	}
}
func (e *userGroupsEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

type userGroupsPatch struct {
	createdAt *time.Time
	createdBy *ntypes.Int64
	groupID   *ntypes.Int64
	updatedAt *time.Time
	updatedBy *ntypes.Int64
	userID    *ntypes.Int64
}

type userGroupsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func scanUserGroupsRows(rows *sql.Rows) ([]*userGroupsEntity, error) {
	var (
		entities []*userGroupsEntity
		err      error
	)
	for rows.Next() {
		var ent userGroupsEntity
		err = rows.Scan(
			&ent.createdAt,
			&ent.createdBy,
			&ent.groupID,
			&ent.updatedAt,
			&ent.updatedBy,
			&ent.userID,
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

func (r *userGroupsRepositoryBase) count(c *userGroupsCriteria) (int64, error) {

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

func (r *userGroupsRepositoryBase) find(c *userGroupsCriteria) ([]*userGroupsEntity, error) {

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

	return scanUserGroupsRows(rows)
}
func (r *userGroupsRepositoryBase) findIter(c *userGroupsCriteria) (*userGroupsIterator, error) {

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
func (r *userGroupsRepositoryBase) findOneByUserIDAndGroupID(userID int64, groupID int64) (*userGroupsEntity, error) {
	var (
		entity userGroupsEntity
	)
	query := `SELECT created_at, created_by, group_id, updated_at, updated_by, user_id FROM charon.user_groups WHERE user_id = $1 AND group_id = $2`
	err := r.db.QueryRow(query, userID, groupID).Scan(
		&entity.createdAt,
		&entity.createdBy,
		&entity.groupID,
		&entity.updatedAt,
		&entity.updatedBy,
		&entity.userID,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *userGroupsRepositoryBase) insert(e *userGroupsEntity) (*userGroupsEntity, error) {
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tableUserGroupsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserGroupsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserGroupsColumnGroupID, "", e.groupID)
	insert.AddExpr(tableUserGroupsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserGroupsColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserGroupsColumnUserID, "", e.userID)

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
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userGroupsRepositoryBase) upsert(e *userGroupsEntity, p *userGroupsPatch, inf ...string) (*userGroupsEntity, error) {
	insert := pqcomp.New(0, 6)
	update := insert.Compose(6)
	insert.AddExpr(tableUserGroupsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserGroupsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserGroupsColumnGroupID, "", e.groupID)
	insert.AddExpr(tableUserGroupsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserGroupsColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserGroupsColumnUserID, "", e.userID)
	if len(inf) > 0 {
		update.AddExpr(tableUserGroupsColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tableUserGroupsColumnCreatedBy, "=", p.createdBy)
		update.AddExpr(tableUserGroupsColumnGroupID, "=", p.groupID)
		update.AddExpr(tableUserGroupsColumnUpdatedAt, "=", p.updatedAt)
		update.AddExpr(tableUserGroupsColumnUpdatedBy, "=", p.updatedBy)
		update.AddExpr(tableUserGroupsColumnUserID, "=", p.userID)
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
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userGroupsRepositoryBase) updateOneByUserIDAndGroupID(userID int64, groupID int64, patch *userGroupsPatch) (*userGroupsEntity, error) {
	update := pqcomp.New(2, 6)
	update.AddArg(userID)
	update.AddArg(groupID)
	if patch.createdAt != nil {
		update.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.Equal, patch.createdBy)
	update.AddExpr(tableUserGroupsColumnGroupID, pqcomp.Equal, patch.groupID)
	if patch.updatedAt != nil {
		update.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.Equal, patch.updatedBy)
	update.AddExpr(tableUserGroupsColumnUserID, pqcomp.Equal, patch.userID)

	if update.Len() == 0 {
		return nil, errors.New("userGroups update failure, nothing to update")
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
	var e userGroupsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
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
	createdAt           time.Time
	createdBy           *ntypes.Int64
	groupID             int64
	permissionAction    string
	permissionModule    string
	permissionSubsystem string
	updatedAt           *time.Time
	updatedBy           *ntypes.Int64
	group               *groupEntity
	permission          *permissionEntity
	author              *userEntity
	modifier            *userEntity
}

func (e *groupPermissionsEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tableGroupPermissionsColumnCreatedAt:
		return &e.createdAt, true
	case tableGroupPermissionsColumnCreatedBy:
		return &e.createdBy, true
	case tableGroupPermissionsColumnGroupID:
		return &e.groupID, true
	case tableGroupPermissionsColumnPermissionAction:
		return &e.permissionAction, true
	case tableGroupPermissionsColumnPermissionModule:
		return &e.permissionModule, true
	case tableGroupPermissionsColumnPermissionSubsystem:
		return &e.permissionSubsystem, true
	case tableGroupPermissionsColumnUpdatedAt:
		return &e.updatedAt, true
	case tableGroupPermissionsColumnUpdatedBy:
		return &e.updatedBy, true
	default:
		return nil, false
	}
}
func (e *groupPermissionsEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

type groupPermissionsPatch struct {
	createdAt           *time.Time
	createdBy           *ntypes.Int64
	groupID             *ntypes.Int64
	permissionAction    *ntypes.String
	permissionModule    *ntypes.String
	permissionSubsystem *ntypes.String
	updatedAt           *time.Time
	updatedBy           *ntypes.Int64
}

type groupPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func scanGroupPermissionsRows(rows *sql.Rows) ([]*groupPermissionsEntity, error) {
	var (
		entities []*groupPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent groupPermissionsEntity
		err = rows.Scan(
			&ent.createdAt,
			&ent.createdBy,
			&ent.groupID,
			&ent.permissionAction,
			&ent.permissionModule,
			&ent.permissionSubsystem,
			&ent.updatedAt,
			&ent.updatedBy,
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

func (r *groupPermissionsRepositoryBase) count(c *groupPermissionsCriteria) (int64, error) {

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

func (r *groupPermissionsRepositoryBase) find(c *groupPermissionsCriteria) ([]*groupPermissionsEntity, error) {

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

	return scanGroupPermissionsRows(rows)
}
func (r *groupPermissionsRepositoryBase) findIter(c *groupPermissionsCriteria) (*groupPermissionsIterator, error) {

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
func (r *groupPermissionsRepositoryBase) findOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(groupID int64, permissionSubsystem string, permissionModule string, permissionAction string) (*groupPermissionsEntity, error) {
	var (
		entity groupPermissionsEntity
	)
	query := `SELECT created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by FROM charon.group_permissions WHERE group_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4`
	err := r.db.QueryRow(query, groupID, permissionSubsystem, permissionModule, permissionAction).Scan(
		&entity.createdAt,
		&entity.createdBy,
		&entity.groupID,
		&entity.permissionAction,
		&entity.permissionModule,
		&entity.permissionSubsystem,
		&entity.updatedAt,
		&entity.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *groupPermissionsRepositoryBase) insert(e *groupPermissionsEntity) (*groupPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(tableGroupPermissionsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableGroupPermissionsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableGroupPermissionsColumnGroupID, "", e.groupID)
	insert.AddExpr(tableGroupPermissionsColumnPermissionAction, "", e.permissionAction)
	insert.AddExpr(tableGroupPermissionsColumnPermissionModule, "", e.permissionModule)
	insert.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, "", e.permissionSubsystem)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedBy, "", e.updatedBy)

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
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *groupPermissionsRepositoryBase) upsert(e *groupPermissionsEntity, p *groupPermissionsPatch, inf ...string) (*groupPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	update := insert.Compose(8)
	insert.AddExpr(tableGroupPermissionsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableGroupPermissionsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableGroupPermissionsColumnGroupID, "", e.groupID)
	insert.AddExpr(tableGroupPermissionsColumnPermissionAction, "", e.permissionAction)
	insert.AddExpr(tableGroupPermissionsColumnPermissionModule, "", e.permissionModule)
	insert.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, "", e.permissionSubsystem)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedBy, "", e.updatedBy)
	if len(inf) > 0 {
		update.AddExpr(tableGroupPermissionsColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tableGroupPermissionsColumnCreatedBy, "=", p.createdBy)
		update.AddExpr(tableGroupPermissionsColumnGroupID, "=", p.groupID)
		update.AddExpr(tableGroupPermissionsColumnPermissionAction, "=", p.permissionAction)
		update.AddExpr(tableGroupPermissionsColumnPermissionModule, "=", p.permissionModule)
		update.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, "=", p.permissionSubsystem)
		update.AddExpr(tableGroupPermissionsColumnUpdatedAt, "=", p.updatedAt)
		update.AddExpr(tableGroupPermissionsColumnUpdatedBy, "=", p.updatedBy)
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
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *groupPermissionsRepositoryBase) updateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(groupID int64, permissionSubsystem string, permissionModule string, permissionAction string, patch *groupPermissionsPatch) (*groupPermissionsEntity, error) {
	update := pqcomp.New(2, 8)
	update.AddArg(groupID)
	update.AddArg(permissionSubsystem)
	update.AddArg(permissionModule)
	update.AddArg(permissionAction)
	if patch.createdAt != nil {
		update.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.Equal, patch.createdBy)
	update.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.Equal, patch.groupID)
	update.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.Equal, patch.permissionAction)
	update.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.Equal, patch.permissionModule)
	update.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.Equal, patch.permissionSubsystem)
	if patch.updatedAt != nil {
		update.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.Equal, patch.updatedBy)

	if update.Len() == 0 {
		return nil, errors.New("groupPermissions update failure, nothing to update")
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
	var e groupPermissionsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.createdAt,
		&e.createdBy,
		&e.groupID,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
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
	createdAt           time.Time
	createdBy           *ntypes.Int64
	permissionAction    string
	permissionModule    string
	permissionSubsystem string
	updatedAt           *time.Time
	updatedBy           *ntypes.Int64
	userID              int64
	user                *userEntity
	permission          *permissionEntity
	author              *userEntity
	modifier            *userEntity
}

func (e *userPermissionsEntity) prop(cn string) (interface{}, bool) {
	switch cn {
	case tableUserPermissionsColumnCreatedAt:
		return &e.createdAt, true
	case tableUserPermissionsColumnCreatedBy:
		return &e.createdBy, true
	case tableUserPermissionsColumnPermissionAction:
		return &e.permissionAction, true
	case tableUserPermissionsColumnPermissionModule:
		return &e.permissionModule, true
	case tableUserPermissionsColumnPermissionSubsystem:
		return &e.permissionSubsystem, true
	case tableUserPermissionsColumnUpdatedAt:
		return &e.updatedAt, true
	case tableUserPermissionsColumnUpdatedBy:
		return &e.updatedBy, true
	case tableUserPermissionsColumnUserID:
		return &e.userID, true
	default:
		return nil, false
	}
}
func (e *userPermissionsEntity) props(cns ...string) ([]interface{}, error) {

	res := make([]interface{}, 0, len(cns))
	for _, cn := range cns {
		if prop, ok := e.prop(cn); ok {
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

	props, err := ent.props(cols...)
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

	if len(c.sort) > 0 {
		i := 0
		com.WriteString(" ORDER BY ")
		for cn, asc := range c.sort {
			if i > 0 {
				com.WriteString(", ")
			}
			com.WriteString(cn)
			if !asc {
				com.WriteString(" DESC ")
			}
			i++
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

type userPermissionsPatch struct {
	createdAt           *time.Time
	createdBy           *ntypes.Int64
	permissionAction    *ntypes.String
	permissionModule    *ntypes.String
	permissionSubsystem *ntypes.String
	updatedAt           *time.Time
	updatedBy           *ntypes.Int64
	userID              *ntypes.Int64
}

type userPermissionsRepositoryBase struct {
	table   string
	columns []string
	db      *sql.DB
	dbg     bool
	log     log.Logger
}

func scanUserPermissionsRows(rows *sql.Rows) ([]*userPermissionsEntity, error) {
	var (
		entities []*userPermissionsEntity
		err      error
	)
	for rows.Next() {
		var ent userPermissionsEntity
		err = rows.Scan(
			&ent.createdAt,
			&ent.createdBy,
			&ent.permissionAction,
			&ent.permissionModule,
			&ent.permissionSubsystem,
			&ent.updatedAt,
			&ent.updatedBy,
			&ent.userID,
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

func (r *userPermissionsRepositoryBase) count(c *userPermissionsCriteria) (int64, error) {

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

func (r *userPermissionsRepositoryBase) find(c *userPermissionsCriteria) ([]*userPermissionsEntity, error) {

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

	return scanUserPermissionsRows(rows)
}
func (r *userPermissionsRepositoryBase) findIter(c *userPermissionsCriteria) (*userPermissionsIterator, error) {

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
func (r *userPermissionsRepositoryBase) findOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(userID int64, permissionSubsystem string, permissionModule string, permissionAction string) (*userPermissionsEntity, error) {
	var (
		entity userPermissionsEntity
	)
	query := `SELECT created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id FROM charon.user_permissions WHERE user_id = $1 AND permission_subsystem = $2 AND permission_module = $3 AND permission_action = $4`
	err := r.db.QueryRow(query, userID, permissionSubsystem, permissionModule, permissionAction).Scan(
		&entity.createdAt,
		&entity.createdBy,
		&entity.permissionAction,
		&entity.permissionModule,
		&entity.permissionSubsystem,
		&entity.updatedAt,
		&entity.updatedBy,
		&entity.userID,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *userPermissionsRepositoryBase) insert(e *userPermissionsEntity) (*userPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	insert.AddExpr(tableUserPermissionsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserPermissionsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserPermissionsColumnPermissionAction, "", e.permissionAction)
	insert.AddExpr(tableUserPermissionsColumnPermissionModule, "", e.permissionModule)
	insert.AddExpr(tableUserPermissionsColumnPermissionSubsystem, "", e.permissionSubsystem)
	insert.AddExpr(tableUserPermissionsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserPermissionsColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserPermissionsColumnUserID, "", e.userID)

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
		&e.createdAt,
		&e.createdBy,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userPermissionsRepositoryBase) upsert(e *userPermissionsEntity, p *userPermissionsPatch, inf ...string) (*userPermissionsEntity, error) {
	insert := pqcomp.New(0, 8)
	update := insert.Compose(8)
	insert.AddExpr(tableUserPermissionsColumnCreatedAt, "", e.createdAt)
	insert.AddExpr(tableUserPermissionsColumnCreatedBy, "", e.createdBy)
	insert.AddExpr(tableUserPermissionsColumnPermissionAction, "", e.permissionAction)
	insert.AddExpr(tableUserPermissionsColumnPermissionModule, "", e.permissionModule)
	insert.AddExpr(tableUserPermissionsColumnPermissionSubsystem, "", e.permissionSubsystem)
	insert.AddExpr(tableUserPermissionsColumnUpdatedAt, "", e.updatedAt)
	insert.AddExpr(tableUserPermissionsColumnUpdatedBy, "", e.updatedBy)
	insert.AddExpr(tableUserPermissionsColumnUserID, "", e.userID)
	if len(inf) > 0 {
		update.AddExpr(tableUserPermissionsColumnCreatedAt, "=", p.createdAt)
		update.AddExpr(tableUserPermissionsColumnCreatedBy, "=", p.createdBy)
		update.AddExpr(tableUserPermissionsColumnPermissionAction, "=", p.permissionAction)
		update.AddExpr(tableUserPermissionsColumnPermissionModule, "=", p.permissionModule)
		update.AddExpr(tableUserPermissionsColumnPermissionSubsystem, "=", p.permissionSubsystem)
		update.AddExpr(tableUserPermissionsColumnUpdatedAt, "=", p.updatedAt)
		update.AddExpr(tableUserPermissionsColumnUpdatedBy, "=", p.updatedBy)
		update.AddExpr(tableUserPermissionsColumnUserID, "=", p.userID)
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
		&e.createdAt,
		&e.createdBy,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (r *userPermissionsRepositoryBase) updateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(userID int64, permissionSubsystem string, permissionModule string, permissionAction string, patch *userPermissionsPatch) (*userPermissionsEntity, error) {
	update := pqcomp.New(2, 8)
	update.AddArg(userID)
	update.AddArg(permissionSubsystem)
	update.AddArg(permissionModule)
	update.AddArg(permissionAction)
	if patch.createdAt != nil {
		update.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.Equal, patch.createdAt)

	}
	update.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.Equal, patch.createdBy)
	update.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.Equal, patch.permissionAction)
	update.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.Equal, patch.permissionModule)
	update.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.Equal, patch.permissionSubsystem)
	if patch.updatedAt != nil {
		update.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.Equal, patch.updatedAt)
	} else {
		update.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.Equal, "NOW()")
	}
	update.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.Equal, patch.updatedBy)
	update.AddExpr(tableUserPermissionsColumnUserID, pqcomp.Equal, patch.userID)

	if update.Len() == 0 {
		return nil, errors.New("userPermissions update failure, nothing to update")
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
	var e userPermissionsEntity
	err := r.db.QueryRow(query, update.Args()...).Scan(
		&e.createdAt,
		&e.createdBy,
		&e.permissionAction,
		&e.permissionModule,
		&e.permissionSubsystem,
		&e.updatedAt,
		&e.updatedBy,
		&e.userID,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
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
