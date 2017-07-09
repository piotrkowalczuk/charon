package model

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

func joinClause(comp *Composer, jt JoinType, on string) (ok bool, err error) {
	if jt != JoinDoNot {
		switch jt {
		case JoinInner:
			if _, err = comp.WriteString(" INNER JOIN "); err != nil {
				return
			}
		case JoinLeft:
			if _, err = comp.WriteString(" LEFT JOIN "); err != nil {
				return
			}
		case JoinRight:
			if _, err = comp.WriteString(" RIGHT JOIN "); err != nil {
				return
			}
		case JoinCross:
			if _, err = comp.WriteString(" CROSS JOIN "); err != nil {
				return
			}
		default:
			return
		}
		if _, err = comp.WriteString(on); err != nil {
			return
		}
		comp.Dirty = true
		ok = true
		return
	}
	return
}

// LogFunc represents function that can be passed into repository to log query result.
type LogFunc func(err error, ent, fnc, sql string, args ...interface{})

// Rows ...
type Rows interface {
	io.Closer
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...interface{}) error
}

const (
	TableUser                         = "charon.user"
	TableUserColumnConfirmationToken  = "confirmation_token"
	TableUserColumnCreatedAt          = "created_at"
	TableUserColumnCreatedBy          = "created_by"
	TableUserColumnFirstName          = "first_name"
	TableUserColumnID                 = "id"
	TableUserColumnIsActive           = "is_active"
	TableUserColumnIsConfirmed        = "is_confirmed"
	TableUserColumnIsStaff            = "is_staff"
	TableUserColumnIsSuperuser        = "is_superuser"
	TableUserColumnLastLoginAt        = "last_login_at"
	TableUserColumnLastName           = "last_name"
	TableUserColumnPassword           = "password"
	TableUserColumnUpdatedAt          = "updated_at"
	TableUserColumnUpdatedBy          = "updated_by"
	TableUserColumnUsername           = "username"
	TableUserConstraintUsernameUnique = "charon.user_username_key"

	TableUserConstraintPrimaryKey = "charon.user_id_pkey"

	TableUserConstraintCreatedByForeignKey = "charon.user_created_by_fkey"

	TableUserConstraintUpdatedByForeignKey = "charon.user_updated_by_fkey"
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

// UserEntity ...
type UserEntity struct {
	// ConfirmationToken ...
	ConfirmationToken []byte
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
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
	LastLoginAt pq.NullTime
	// LastName ...
	LastName string
	// Password ...
	Password []byte
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// UpdatedBy ...
	UpdatedBy ntypes.Int64
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
	if len(cns) == 0 {
		cns = TableUserColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	ConfirmationToken []byte
	CreatedAt         *qtypes.Timestamp
	CreatedBy         *qtypes.Int64
	FirstName         *qtypes.String
	ID                *qtypes.Int64
	IsActive          ntypes.Bool
	IsConfirmed       ntypes.Bool
	IsStaff           ntypes.Bool
	IsSuperuser       ntypes.Bool
	LastLoginAt       *qtypes.Timestamp
	LastName          *qtypes.String
	Password          []byte
	UpdatedAt         *qtypes.Timestamp
	UpdatedBy         *qtypes.Int64
	Username          *qtypes.String
}

type UserFindExpr struct {
	Where         *UserCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type UserCountExpr struct {
	Where        *UserCriteria
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserJoin struct {
	On, Where    *UserCriteria
	Fetch        bool
	Kind         JoinType
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserPatch struct {
	ConfirmationToken []byte
	CreatedAt         pq.NullTime
	CreatedBy         ntypes.Int64
	FirstName         ntypes.String
	IsActive          ntypes.Bool
	IsConfirmed       ntypes.Bool
	IsStaff           ntypes.Bool
	IsSuperuser       ntypes.Bool
	LastLoginAt       pq.NullTime
	LastName          ntypes.String
	Password          []byte
	UpdatedAt         pq.NullTime
	UpdatedBy         ntypes.Int64
	Username          ntypes.String
}

func ScanUserRows(rows Rows) (entities []*UserEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type UserRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserRepositoryBase) InsertQuery(e *UserEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(15)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if e.ConfirmationToken != nil {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnConfirmationToken); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.ConfirmationToken)
		insert.Dirty = true
	}

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.CreatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnFirstName); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.FirstName)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsActive); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.IsActive)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsConfirmed); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.IsConfirmed)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsStaff); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.IsStaff)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsSuperuser); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.IsSuperuser)
	insert.Dirty = true

	if e.LastLoginAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnLastLoginAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.LastLoginAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnLastName); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.LastName)
	insert.Dirty = true

	if e.Password != nil {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnPassword); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.Password)
		insert.Dirty = true
	}

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UpdatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnUsername); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Username)
	insert.Dirty = true

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *UserRepositoryBase) Insert(ctx context.Context, e *UserEntity) (*UserEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.ConfirmationToken,
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
	if r.Log != nil {
		r.Log(err, "User", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func UserCriteriaWhereClause(comp *Composer, c *UserCriteria, id int) error {
	if c.ConfirmationToken != nil {
		if comp.Dirty {
			comp.WriteString(" AND ")
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnConfirmationToken); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.ConfirmationToken)
		comp.Dirty = true
	}

	QueryTimestampWhereClause(c.CreatedAt, id, TableUserColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableUserColumnCreatedBy, comp, And)

	QueryStringWhereClause(c.FirstName, id, TableUserColumnFirstName, comp, And)

	QueryInt64WhereClause(c.ID, id, TableUserColumnID, comp, And)

	if c.IsActive.Valid {
		if comp.Dirty {
			if _, err := comp.WriteString(" AND "); err != nil {
				return err
			}
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnIsActive); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.IsActive)
		comp.Dirty = true
	}

	if c.IsConfirmed.Valid {
		if comp.Dirty {
			if _, err := comp.WriteString(" AND "); err != nil {
				return err
			}
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnIsConfirmed); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.IsConfirmed)
		comp.Dirty = true
	}

	if c.IsStaff.Valid {
		if comp.Dirty {
			if _, err := comp.WriteString(" AND "); err != nil {
				return err
			}
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnIsStaff); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.IsStaff)
		comp.Dirty = true
	}

	if c.IsSuperuser.Valid {
		if comp.Dirty {
			if _, err := comp.WriteString(" AND "); err != nil {
				return err
			}
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnIsSuperuser); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.IsSuperuser)
		comp.Dirty = true
	}

	QueryTimestampWhereClause(c.LastLoginAt, id, TableUserColumnLastLoginAt, comp, And)

	QueryStringWhereClause(c.LastName, id, TableUserColumnLastName, comp, And)

	if c.Password != nil {
		if comp.Dirty {
			comp.WriteString(" AND ")
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableUserColumnPassword); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.Password)
		comp.Dirty = true
	}

	QueryTimestampWhereClause(c.UpdatedAt, id, TableUserColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableUserColumnUpdatedBy, comp, And)

	QueryStringWhereClause(c.Username, id, TableUserColumnUsername, comp, And)

	return nil
}

func (r *UserRepositoryBase) FindQuery(fe *UserFindExpr) (string, []interface{}, error) {
	comp := NewComposer(15)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.confirmation_token, t0.created_at, t0.created_by, t0.first_name, t0.id, t0.is_active, t0.is_confirmed, t0.is_staff, t0.is_superuser, t0.last_login_at, t0.last_name, t0.password, t0.updated_at, t0.updated_by, t0.username")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}

	if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinAuthor != nil {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t1 ON t0.created_by=t1.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 1); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinModifier != nil {
		joinClause(comp, fe.JoinModifier.Kind, "charon.user AS t2 ON t0.updated_by=t2.id")
		if fe.JoinModifier.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinModifier.On, 2); err != nil {
				return "", nil, err
			}
		}
	}

	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinModifier.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TableUserColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *UserRepositoryBase) Find(ctx context.Context, fe *UserFindExpr) ([]*UserEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "User", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*UserEntity
	var props []interface{}
	for rows.Next() {
		var ent UserEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
			ent.Modifier = &UserEntity{}
			if prop, err = ent.Modifier.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "User", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *UserRepositoryBase) FindIter(ctx context.Context, fe *UserFindExpr) (*UserIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "User", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &UserIterator{
		rows: rows,
		cols: []string{"confirmation_token", "created_at", "created_by", "first_name", "id", "is_active", "is_confirmed", "is_staff", "is_superuser", "last_login_at", "last_name", "password", "updated_at", "updated_by", "username"},
	}, nil
}
func (r *UserRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*UserEntity, error) {
	find := NewComposer(15)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableUser)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		ent UserEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "User", "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *UserRepositoryBase) FindOneByUsername(ctx context.Context, userUsername string) (*UserEntity, error) {
	find := NewComposer(15)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableUser)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserColumnUsername)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userUsername)

	var (
		ent UserEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *UserRepositoryBase) UpdateOneByIDQuery(pk int64, p *UserPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(15)
	if p.ConfirmationToken != nil {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnConfirmationToken); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.ConfirmationToken)
		update.Dirty = true

	}

	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.FirstName.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnFirstName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.FirstName)
		update.Dirty = true
	}

	if p.IsActive.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsActive); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsActive)
		update.Dirty = true
	}

	if p.IsConfirmed.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsConfirmed); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsConfirmed)
		update.Dirty = true
	}

	if p.IsStaff.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsStaff); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsStaff)
		update.Dirty = true
	}

	if p.IsSuperuser.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsSuperuser); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsSuperuser)
		update.Dirty = true
	}

	if p.LastLoginAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnLastLoginAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.LastLoginAt)
		update.Dirty = true

	}

	if p.LastName.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnLastName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.LastName)
		update.Dirty = true
	}

	if p.Password != nil {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnPassword); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Password)
		update.Dirty = true

	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if p.Username.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUsername); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Username)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("User update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")

	update.WriteString(TableUserColumnID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(pk)

	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
	}
	return buf.String(), update.Args(), nil
}
func (r *UserRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *UserPatch) (*UserEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent UserEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "User", "update by primary key", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *UserRepositoryBase) UpdateOneByUsernameQuery(userUsername string, p *UserPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(1)
	if p.ConfirmationToken != nil {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnConfirmationToken); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.ConfirmationToken)
		update.Dirty = true

	}

	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.FirstName.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnFirstName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.FirstName)
		update.Dirty = true
	}

	if p.IsActive.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsActive); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsActive)
		update.Dirty = true
	}

	if p.IsConfirmed.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsConfirmed); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsConfirmed)
		update.Dirty = true
	}

	if p.IsStaff.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsStaff); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsStaff)
		update.Dirty = true
	}

	if p.IsSuperuser.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnIsSuperuser); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.IsSuperuser)
		update.Dirty = true
	}

	if p.LastLoginAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnLastLoginAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.LastLoginAt)
		update.Dirty = true

	}

	if p.LastName.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnLastName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.LastName)
		update.Dirty = true
	}

	if p.Password != nil {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnPassword); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Password)
		update.Dirty = true

	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if p.Username.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserColumnUsername); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Username)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("User update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableUserColumnUsername)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userUsername)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
	}
	return buf.String(), update.Args(), nil
}
func (r *UserRepositoryBase) UpdateOneByUsername(ctx context.Context, userUsername string, p *UserPatch) (*UserEntity, error) {
	query, args, err := r.UpdateOneByUsernameQuery(userUsername, p)
	if err != nil {
		return nil, err
	}
	var ent UserEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "User", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *UserRepositoryBase) UpsertQuery(e *UserEntity, p *UserPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(30)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if e.ConfirmationToken != nil {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnConfirmationToken); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.ConfirmationToken)
		upsert.Dirty = true
	}

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.CreatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnFirstName); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.FirstName)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsActive); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.IsActive)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsConfirmed); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.IsConfirmed)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsStaff); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.IsStaff)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnIsSuperuser); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.IsSuperuser)
	upsert.Dirty = true

	if e.LastLoginAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnLastLoginAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.LastLoginAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnLastName); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.LastName)
	upsert.Dirty = true

	if e.Password != nil {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnPassword); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.Password)
		upsert.Dirty = true
	}

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UpdatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserColumnUsername); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Username)
	upsert.Dirty = true

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.ConfirmationToken != nil {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnConfirmationToken); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.ConfirmationToken)
			upsert.Dirty = true

		}

		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.CreatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnCreatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedBy)
			upsert.Dirty = true
		}

		if p.FirstName.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnFirstName); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.FirstName)
			upsert.Dirty = true
		}

		if p.IsActive.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnIsActive); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.IsActive)
			upsert.Dirty = true
		}

		if p.IsConfirmed.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnIsConfirmed); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.IsConfirmed)
			upsert.Dirty = true
		}

		if p.IsStaff.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnIsStaff); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.IsStaff)
			upsert.Dirty = true
		}

		if p.IsSuperuser.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnIsSuperuser); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.IsSuperuser)
			upsert.Dirty = true
		}

		if p.LastLoginAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnLastLoginAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.LastLoginAt)
			upsert.Dirty = true

		}

		if p.LastName.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnLastName); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.LastName)
			upsert.Dirty = true
		}

		if p.Password != nil {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnPassword); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Password)
			upsert.Dirty = true

		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

		if p.UpdatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnUpdatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedBy)
			upsert.Dirty = true
		}

		if p.Username.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserColumnUsername); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Username)
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("confirmation_token, created_at, created_by, first_name, id, is_active, is_confirmed, is_staff, is_superuser, last_login_at, last_name, password, updated_at, updated_by, username")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *UserRepositoryBase) Upsert(ctx context.Context, e *UserEntity, p *UserPatch, inf ...string) (*UserEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.ConfirmationToken,
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
	if r.Log != nil {
		r.Log(err, "User", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *UserRepositoryBase) Count(ctx context.Context, c *UserCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},

		JoinAuthor:   c.JoinAuthor,
		JoinModifier: c.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "User", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *UserRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	find := NewComposer(15)
	find.WriteString("DELETE FROM ")
	find.WriteString(TableUser)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	res, err := r.DB.ExecContext(ctx, find.String(), find.Args()...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	TableGroup                     = "charon.group"
	TableGroupColumnCreatedAt      = "created_at"
	TableGroupColumnCreatedBy      = "created_by"
	TableGroupColumnDescription    = "description"
	TableGroupColumnID             = "id"
	TableGroupColumnName           = "name"
	TableGroupColumnUpdatedAt      = "updated_at"
	TableGroupColumnUpdatedBy      = "updated_by"
	TableGroupConstraintNameUnique = "charon.group_name_key"

	TableGroupConstraintPrimaryKey = "charon.group_id_pkey"

	TableGroupConstraintCreatedByForeignKey = "charon.group_created_by_fkey"

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

// GroupEntity ...
type GroupEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
	// Description ...
	Description ntypes.String
	// ID ...
	ID int64
	// Name ...
	Name string
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// UpdatedBy ...
	UpdatedBy ntypes.Int64
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
	if len(cns) == 0 {
		cns = TableGroupColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	CreatedAt   *qtypes.Timestamp
	CreatedBy   *qtypes.Int64
	Description *qtypes.String
	ID          *qtypes.Int64
	Name        *qtypes.String
	UpdatedAt   *qtypes.Timestamp
	UpdatedBy   *qtypes.Int64
}

type GroupFindExpr struct {
	Where         *GroupCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type GroupCountExpr struct {
	Where        *GroupCriteria
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupJoin struct {
	On, Where    *GroupCriteria
	Fetch        bool
	Kind         JoinType
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupPatch struct {
	CreatedAt   pq.NullTime
	CreatedBy   ntypes.Int64
	Description ntypes.String
	Name        ntypes.String
	UpdatedAt   pq.NullTime
	UpdatedBy   ntypes.Int64
}

func ScanGroupRows(rows Rows) (entities []*GroupEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type GroupRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *GroupRepositoryBase) InsertQuery(e *GroupEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(7)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.CreatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnDescription); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Description)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnName); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Name)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UpdatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *GroupRepositoryBase) Insert(ctx context.Context, e *GroupEntity) (*GroupEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		r.Log(err, "Group", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func GroupCriteriaWhereClause(comp *Composer, c *GroupCriteria, id int) error {
	QueryTimestampWhereClause(c.CreatedAt, id, TableGroupColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableGroupColumnCreatedBy, comp, And)

	QueryStringWhereClause(c.Description, id, TableGroupColumnDescription, comp, And)

	QueryInt64WhereClause(c.ID, id, TableGroupColumnID, comp, And)

	QueryStringWhereClause(c.Name, id, TableGroupColumnName, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TableGroupColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableGroupColumnUpdatedBy, comp, And)

	return nil
}

func (r *GroupRepositoryBase) FindQuery(fe *GroupFindExpr) (string, []interface{}, error) {
	comp := NewComposer(7)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.created_at, t0.created_by, t0.description, t0.id, t0.name, t0.updated_at, t0.updated_by")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}

	if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinAuthor != nil {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t1 ON t0.created_by=t1.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 1); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinModifier != nil {
		joinClause(comp, fe.JoinModifier.Kind, "charon.user AS t2 ON t0.updated_by=t2.id")
		if fe.JoinModifier.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinModifier.On, 2); err != nil {
				return "", nil, err
			}
		}
	}

	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := GroupCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinModifier.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TableGroupColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *GroupRepositoryBase) Find(ctx context.Context, fe *GroupFindExpr) ([]*GroupEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "Group", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*GroupEntity
	var props []interface{}
	for rows.Next() {
		var ent GroupEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
			ent.Modifier = &UserEntity{}
			if prop, err = ent.Modifier.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "Group", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *GroupRepositoryBase) FindIter(ctx context.Context, fe *GroupFindExpr) (*GroupIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "Group", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &GroupIterator{
		rows: rows,
		cols: []string{"created_at", "created_by", "description", "id", "name", "updated_at", "updated_by"},
	}, nil
}
func (r *GroupRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*GroupEntity, error) {
	find := NewComposer(7)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableGroup)
	find.WriteString(" WHERE ")
	find.WriteString(TableGroupColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		ent GroupEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Group", "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *GroupRepositoryBase) FindOneByName(ctx context.Context, groupName string) (*GroupEntity, error) {
	find := NewComposer(7)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableGroup)
	find.WriteString(" WHERE ")
	find.WriteString(TableGroupColumnName)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(groupName)

	var (
		ent GroupEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *GroupRepositoryBase) UpdateOneByIDQuery(pk int64, p *GroupPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(7)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.Description.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnDescription); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Description)
		update.Dirty = true
	}

	if p.Name.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Name)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("Group update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")

	update.WriteString(TableGroupColumnID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(pk)

	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
	}
	return buf.String(), update.Args(), nil
}
func (r *GroupRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *GroupPatch) (*GroupEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent GroupEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Group", "update by primary key", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *GroupRepositoryBase) UpdateOneByNameQuery(groupName string, p *GroupPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(1)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.Description.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnDescription); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Description)
		update.Dirty = true
	}

	if p.Name.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnName); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Name)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("Group update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableGroupColumnName)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(groupName)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
	}
	return buf.String(), update.Args(), nil
}
func (r *GroupRepositoryBase) UpdateOneByName(ctx context.Context, groupName string, p *GroupPatch) (*GroupEntity, error) {
	query, args, err := r.UpdateOneByNameQuery(groupName, p)
	if err != nil {
		return nil, err
	}
	var ent GroupEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Group", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *GroupRepositoryBase) UpsertQuery(e *GroupEntity, p *GroupPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(14)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.CreatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnDescription); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Description)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnName); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Name)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UpdatedBy)
	upsert.Dirty = true

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.CreatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnCreatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedBy)
			upsert.Dirty = true
		}

		if p.Description.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnDescription); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Description)
			upsert.Dirty = true
		}

		if p.Name.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnName); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Name)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

		if p.UpdatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupColumnUpdatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedBy)
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("created_at, created_by, description, id, name, updated_at, updated_by")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *GroupRepositoryBase) Upsert(ctx context.Context, e *GroupEntity, p *GroupPatch, inf ...string) (*GroupEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		r.Log(err, "Group", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *GroupRepositoryBase) Count(ctx context.Context, c *GroupCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&GroupFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},

		JoinAuthor:   c.JoinAuthor,
		JoinModifier: c.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "Group", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *GroupRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	find := NewComposer(7)
	find.WriteString("DELETE FROM ")
	find.WriteString(TableGroup)
	find.WriteString(" WHERE ")
	find.WriteString(TableGroupColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	res, err := r.DB.ExecContext(ctx, find.String(), find.Args()...)
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
	TablePermissionConstraintSubsystemModuleActionUnique = "charon.permission_subsystem_module_action_key"

	TablePermissionConstraintPrimaryKey = "charon.permission_id_pkey"
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

// PermissionEntity ...
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
	UpdatedAt pq.NullTime
	// Users ...
	Users []*UserEntity
	// Groups ...
	Groups []*GroupEntity
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
	if len(cns) == 0 {
		cns = TablePermissionColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	Action    *qtypes.String
	CreatedAt *qtypes.Timestamp
	ID        *qtypes.Int64
	Module    *qtypes.String
	Subsystem *qtypes.String
	UpdatedAt *qtypes.Timestamp
}

type PermissionFindExpr struct {
	Where         *PermissionCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
}

type PermissionCountExpr struct {
	Where *PermissionCriteria
}

type PermissionJoin struct {
	On, Where *PermissionCriteria
	Fetch     bool
	Kind      JoinType
}

type PermissionPatch struct {
	Action    ntypes.String
	CreatedAt pq.NullTime
	Module    ntypes.String
	Subsystem ntypes.String
	UpdatedAt pq.NullTime
}

func ScanPermissionRows(rows Rows) (entities []*PermissionEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type PermissionRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *PermissionRepositoryBase) InsertQuery(e *PermissionEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(6)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnAction); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Action)
	insert.Dirty = true

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TablePermissionColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnModule); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Module)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnSubsystem); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.Subsystem)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("action, created_at, id, module, subsystem, updated_at")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *PermissionRepositoryBase) Insert(ctx context.Context, e *PermissionEntity) (*PermissionEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if r.Log != nil {
		r.Log(err, "Permission", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func PermissionCriteriaWhereClause(comp *Composer, c *PermissionCriteria, id int) error {
	QueryStringWhereClause(c.Action, id, TablePermissionColumnAction, comp, And)

	QueryTimestampWhereClause(c.CreatedAt, id, TablePermissionColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.ID, id, TablePermissionColumnID, comp, And)

	QueryStringWhereClause(c.Module, id, TablePermissionColumnModule, comp, And)

	QueryStringWhereClause(c.Subsystem, id, TablePermissionColumnSubsystem, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TablePermissionColumnUpdatedAt, comp, And)

	return nil
}

func (r *PermissionRepositoryBase) FindQuery(fe *PermissionFindExpr) (string, []interface{}, error) {
	comp := NewComposer(6)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.action, t0.created_at, t0.id, t0.module, t0.subsystem, t0.updated_at")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := PermissionCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TablePermissionColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *PermissionRepositoryBase) Find(ctx context.Context, fe *PermissionFindExpr) ([]*PermissionEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "Permission", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*PermissionEntity
	var props []interface{}
	for rows.Next() {
		var ent PermissionEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "Permission", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *PermissionRepositoryBase) FindIter(ctx context.Context, fe *PermissionFindExpr) (*PermissionIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "Permission", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &PermissionIterator{
		rows: rows,
		cols: []string{"action", "created_at", "id", "module", "subsystem", "updated_at"},
	}, nil
}
func (r *PermissionRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*PermissionEntity, error) {
	find := NewComposer(6)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("action, created_at, id, module, subsystem, updated_at")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TablePermission)
	find.WriteString(" WHERE ")
	find.WriteString(TablePermissionColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		ent PermissionEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Permission", "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *PermissionRepositoryBase) FindOneBySubsystemAndModuleAndAction(ctx context.Context, permissionSubsystem string, permissionModule string, permissionAction string) (*PermissionEntity, error) {
	find := NewComposer(6)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("action, created_at, id, module, subsystem, updated_at")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TablePermission)
	find.WriteString(" WHERE ")
	find.WriteString(TablePermissionColumnSubsystem)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(permissionSubsystem)
	find.WriteString(" AND ")
	find.WriteString(TablePermissionColumnModule)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(permissionModule)
	find.WriteString(" AND ")
	find.WriteString(TablePermissionColumnAction)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(permissionAction)

	var (
		ent PermissionEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *PermissionRepositoryBase) UpdateOneByIDQuery(pk int64, p *PermissionPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(6)
	if p.Action.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnAction); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Action)
		update.Dirty = true
	}

	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.Module.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnModule); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Module)
		update.Dirty = true
	}

	if p.Subsystem.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnSubsystem); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Subsystem)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("Permission update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")

	update.WriteString(TablePermissionColumnID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(pk)

	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("action, created_at, id, module, subsystem, updated_at")
	}
	return buf.String(), update.Args(), nil
}
func (r *PermissionRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *PermissionPatch) (*PermissionEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent PermissionEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Permission", "update by primary key", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *PermissionRepositoryBase) UpdateOneBySubsystemAndModuleAndActionQuery(permissionSubsystem string, permissionModule string, permissionAction string, p *PermissionPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(3)
	if p.Action.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnAction); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Action)
		update.Dirty = true
	}

	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.Module.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnModule); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Module)
		update.Dirty = true
	}

	if p.Subsystem.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnSubsystem); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Subsystem)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("Permission update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TablePermissionColumnSubsystem)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(permissionSubsystem)
	update.WriteString(" AND ")
	update.WriteString(TablePermissionColumnModule)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(permissionModule)
	update.WriteString(" AND ")
	update.WriteString(TablePermissionColumnAction)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(permissionAction)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("action, created_at, id, module, subsystem, updated_at")
	}
	return buf.String(), update.Args(), nil
}
func (r *PermissionRepositoryBase) UpdateOneBySubsystemAndModuleAndAction(ctx context.Context, permissionSubsystem string, permissionModule string, permissionAction string, p *PermissionPatch) (*PermissionEntity, error) {
	query, args, err := r.UpdateOneBySubsystemAndModuleAndActionQuery(permissionSubsystem, permissionModule, permissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent PermissionEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "Permission", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *PermissionRepositoryBase) UpsertQuery(e *PermissionEntity, p *PermissionPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(12)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnAction); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Action)
	upsert.Dirty = true

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TablePermissionColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnModule); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Module)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TablePermissionColumnSubsystem); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.Subsystem)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TablePermissionColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.Action.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnAction); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Action)
			upsert.Dirty = true
		}

		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.Module.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnModule); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Module)
			upsert.Dirty = true
		}

		if p.Subsystem.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnSubsystem); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Subsystem)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TablePermissionColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("action, created_at, id, module, subsystem, updated_at")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *PermissionRepositoryBase) Upsert(ctx context.Context, e *PermissionEntity, p *PermissionPatch, inf ...string) (*PermissionEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if r.Log != nil {
		r.Log(err, "Permission", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *PermissionRepositoryBase) Count(ctx context.Context, c *PermissionCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&PermissionFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "Permission", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *PermissionRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	find := NewComposer(6)
	find.WriteString("DELETE FROM ")
	find.WriteString(TablePermission)
	find.WriteString(" WHERE ")
	find.WriteString(TablePermissionColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	res, err := r.DB.ExecContext(ctx, find.String(), find.Args()...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

const (
	TableUserGroups                           = "charon.user_groups"
	TableUserGroupsColumnCreatedAt            = "created_at"
	TableUserGroupsColumnCreatedBy            = "created_by"
	TableUserGroupsColumnGroupID              = "group_id"
	TableUserGroupsColumnUpdatedAt            = "updated_at"
	TableUserGroupsColumnUpdatedBy            = "updated_by"
	TableUserGroupsColumnUserID               = "user_id"
	TableUserGroupsConstraintUserIDForeignKey = "charon.user_groups_user_id_fkey"

	TableUserGroupsConstraintGroupIDForeignKey = "charon.user_groups_group_id_fkey"

	TableUserGroupsConstraintUserIDGroupIDUnique = "charon.user_groups_user_id_group_id_key"

	TableUserGroupsConstraintCreatedByForeignKey = "charon.user_groups_created_by_fkey"

	TableUserGroupsConstraintUpdatedByForeignKey = "charon.user_groups_updated_by_fkey"
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

// UserGroupsEntity ...
type UserGroupsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
	// GroupID ...
	GroupID int64
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// UpdatedBy ...
	UpdatedBy ntypes.Int64
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
	if len(cns) == 0 {
		cns = TableUserGroupsColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	CreatedAt *qtypes.Timestamp
	CreatedBy *qtypes.Int64
	GroupID   *qtypes.Int64
	UpdatedAt *qtypes.Timestamp
	UpdatedBy *qtypes.Int64
	UserID    *qtypes.Int64
}

type UserGroupsFindExpr struct {
	Where         *UserGroupsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
	JoinUser      *UserJoin
	JoinGroup     *GroupJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type UserGroupsCountExpr struct {
	Where        *UserGroupsCriteria
	JoinUser     *UserJoin
	JoinGroup    *GroupJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserGroupsJoin struct {
	On, Where    *UserGroupsCriteria
	Fetch        bool
	Kind         JoinType
	JoinUser     *UserJoin
	JoinGroup    *GroupJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserGroupsPatch struct {
	CreatedAt pq.NullTime
	CreatedBy ntypes.Int64
	GroupID   ntypes.Int64
	UpdatedAt pq.NullTime
	UpdatedBy ntypes.Int64
	UserID    ntypes.Int64
}

func ScanUserGroupsRows(rows Rows) (entities []*UserGroupsEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type UserGroupsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserGroupsRepositoryBase) InsertQuery(e *UserGroupsEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(6)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserGroupsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.CreatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnGroupID); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.GroupID)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UpdatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnUserID); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UserID)
	insert.Dirty = true

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("created_at, created_by, group_id, updated_at, updated_by, user_id")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *UserGroupsRepositoryBase) Insert(ctx context.Context, e *UserGroupsEntity) (*UserGroupsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		r.Log(err, "UserGroups", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func UserGroupsCriteriaWhereClause(comp *Composer, c *UserGroupsCriteria, id int) error {
	QueryTimestampWhereClause(c.CreatedAt, id, TableUserGroupsColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableUserGroupsColumnCreatedBy, comp, And)

	QueryInt64WhereClause(c.GroupID, id, TableUserGroupsColumnGroupID, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TableUserGroupsColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableUserGroupsColumnUpdatedBy, comp, And)

	QueryInt64WhereClause(c.UserID, id, TableUserGroupsColumnUserID, comp, And)

	return nil
}

func (r *UserGroupsRepositoryBase) FindQuery(fe *UserGroupsFindExpr) (string, []interface{}, error) {
	comp := NewComposer(6)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.created_at, t0.created_by, t0.group_id, t0.updated_at, t0.updated_by, t0.user_id")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinUser != nil && fe.JoinUser.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}

	if fe.JoinGroup != nil && fe.JoinGroup.Fetch {
		buf.WriteString(", t2.created_at, t2.created_by, t2.description, t2.id, t2.name, t2.updated_at, t2.updated_by")
	}

	if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}

	if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
		buf.WriteString(", t4.confirmation_token, t4.created_at, t4.created_by, t4.first_name, t4.id, t4.is_active, t4.is_confirmed, t4.is_staff, t4.is_superuser, t4.last_login_at, t4.last_name, t4.password, t4.updated_at, t4.updated_by, t4.username")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinUser != nil {
		joinClause(comp, fe.JoinUser.Kind, "charon.user AS t1 ON t0.user_id=t1.id")
		if fe.JoinUser.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinUser.On, 1); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinGroup != nil {
		joinClause(comp, fe.JoinGroup.Kind, "charon.group AS t2 ON t0.group_id=t2.id")
		if fe.JoinGroup.On != nil {
			comp.Dirty = true
			if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.On, 2); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinAuthor != nil {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t3 ON t0.created_by=t3.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 3); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinModifier != nil {
		joinClause(comp, fe.JoinModifier.Kind, "charon.user AS t4 ON t0.updated_by=t4.id")
		if fe.JoinModifier.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinModifier.On, 4); err != nil {
				return "", nil, err
			}
		}
	}

	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := UserGroupsCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinUser != nil && fe.JoinUser.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinUser.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Where != nil {
		if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 3); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinModifier.Where, 4); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TableUserGroupsColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *UserGroupsRepositoryBase) Find(ctx context.Context, fe *UserGroupsFindExpr) ([]*UserGroupsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "UserGroups", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*UserGroupsEntity
	var props []interface{}
	for rows.Next() {
		var ent UserGroupsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinUser != nil && fe.JoinUser.Fetch {
			ent.User = &UserEntity{}
			if prop, err = ent.User.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinGroup != nil && fe.JoinGroup.Fetch {
			ent.Group = &GroupEntity{}
			if prop, err = ent.Group.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
			ent.Modifier = &UserEntity{}
			if prop, err = ent.Modifier.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "UserGroups", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *UserGroupsRepositoryBase) FindIter(ctx context.Context, fe *UserGroupsFindExpr) (*UserGroupsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "UserGroups", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &UserGroupsIterator{
		rows: rows,
		cols: []string{"created_at", "created_by", "group_id", "updated_at", "updated_by", "user_id"},
	}, nil
}
func (r *UserGroupsRepositoryBase) FindOneByUserIDAndGroupID(ctx context.Context, userGroupsUserID int64, userGroupsGroupID int64) (*UserGroupsEntity, error) {
	find := NewComposer(6)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, group_id, updated_at, updated_by, user_id")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableUserGroups)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserGroupsColumnUserID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userGroupsUserID)
	find.WriteString(" AND ")
	find.WriteString(TableUserGroupsColumnGroupID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userGroupsGroupID)

	var (
		ent UserGroupsEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *UserGroupsRepositoryBase) UpdateOneByUserIDAndGroupIDQuery(userGroupsUserID int64, userGroupsGroupID int64, p *UserGroupsPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(2)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.GroupID.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnGroupID); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.GroupID)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if p.UserID.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserGroupsColumnUserID); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UserID)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("UserGroups update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableUserGroupsColumnUserID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userGroupsUserID)
	update.WriteString(" AND ")
	update.WriteString(TableUserGroupsColumnGroupID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userGroupsGroupID)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, group_id, updated_at, updated_by, user_id")
	}
	return buf.String(), update.Args(), nil
}
func (r *UserGroupsRepositoryBase) UpdateOneByUserIDAndGroupID(ctx context.Context, userGroupsUserID int64, userGroupsGroupID int64, p *UserGroupsPatch) (*UserGroupsEntity, error) {
	query, args, err := r.UpdateOneByUserIDAndGroupIDQuery(userGroupsUserID, userGroupsGroupID, p)
	if err != nil {
		return nil, err
	}
	var ent UserGroupsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "UserGroups", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *UserGroupsRepositoryBase) UpsertQuery(e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(12)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserGroupsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.CreatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnGroupID); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.GroupID)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UpdatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserGroupsColumnUserID); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UserID)
	upsert.Dirty = true

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.CreatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnCreatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedBy)
			upsert.Dirty = true
		}

		if p.GroupID.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnGroupID); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.GroupID)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

		if p.UpdatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnUpdatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedBy)
			upsert.Dirty = true
		}

		if p.UserID.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserGroupsColumnUserID); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UserID)
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("created_at, created_by, group_id, updated_at, updated_by, user_id")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *UserGroupsRepositoryBase) Upsert(ctx context.Context, e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (*UserGroupsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		r.Log(err, "UserGroups", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *UserGroupsRepositoryBase) Count(ctx context.Context, c *UserGroupsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserGroupsFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},

		JoinUser:     c.JoinUser,
		JoinGroup:    c.JoinGroup,
		JoinAuthor:   c.JoinAuthor,
		JoinModifier: c.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "UserGroups", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

const (
	TableGroupPermissions                            = "charon.group_permissions"
	TableGroupPermissionsColumnCreatedAt             = "created_at"
	TableGroupPermissionsColumnCreatedBy             = "created_by"
	TableGroupPermissionsColumnGroupID               = "group_id"
	TableGroupPermissionsColumnPermissionAction      = "permission_action"
	TableGroupPermissionsColumnPermissionModule      = "permission_module"
	TableGroupPermissionsColumnPermissionSubsystem   = "permission_subsystem"
	TableGroupPermissionsColumnUpdatedAt             = "updated_at"
	TableGroupPermissionsColumnUpdatedBy             = "updated_by"
	TableGroupPermissionsConstraintGroupIDForeignKey = "charon.group_permissions_group_id_fkey"

	TableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey = "charon.group_permissions_subsystem_module_action_fkey"

	TableGroupPermissionsConstraintGroupIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.group_permissions_group_id_subsystem_module_action_key"

	TableGroupPermissionsConstraintCreatedByForeignKey = "charon.group_permissions_created_by_fkey"

	TableGroupPermissionsConstraintUpdatedByForeignKey = "charon.group_permissions_updated_by_fkey"
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

// GroupPermissionsEntity ...
type GroupPermissionsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
	// GroupID ...
	GroupID int64
	// PermissionAction ...
	PermissionAction string
	// PermissionModule ...
	PermissionModule string
	// PermissionSubsystem ...
	PermissionSubsystem string
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// UpdatedBy ...
	UpdatedBy ntypes.Int64
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
	if len(cns) == 0 {
		cns = TableGroupPermissionsColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	CreatedAt           *qtypes.Timestamp
	CreatedBy           *qtypes.Int64
	GroupID             *qtypes.Int64
	PermissionAction    *qtypes.String
	PermissionModule    *qtypes.String
	PermissionSubsystem *qtypes.String
	UpdatedAt           *qtypes.Timestamp
	UpdatedBy           *qtypes.Int64
}

type GroupPermissionsFindExpr struct {
	Where         *GroupPermissionsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
	JoinGroup     *GroupJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type GroupPermissionsCountExpr struct {
	Where        *GroupPermissionsCriteria
	JoinGroup    *GroupJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupPermissionsJoin struct {
	On, Where    *GroupPermissionsCriteria
	Fetch        bool
	Kind         JoinType
	JoinGroup    *GroupJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupPermissionsPatch struct {
	CreatedAt           pq.NullTime
	CreatedBy           ntypes.Int64
	GroupID             ntypes.Int64
	PermissionAction    ntypes.String
	PermissionModule    ntypes.String
	PermissionSubsystem ntypes.String
	UpdatedAt           pq.NullTime
	UpdatedBy           ntypes.Int64
}

func ScanGroupPermissionsRows(rows Rows) (entities []*GroupPermissionsEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type GroupPermissionsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *GroupPermissionsRepositoryBase) InsertQuery(e *GroupPermissionsEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(8)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.CreatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnGroupID); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.GroupID)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionAction); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionAction)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionModule); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionModule)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionSubsystem); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionSubsystem)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UpdatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *GroupPermissionsRepositoryBase) Insert(ctx context.Context, e *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func GroupPermissionsCriteriaWhereClause(comp *Composer, c *GroupPermissionsCriteria, id int) error {
	QueryTimestampWhereClause(c.CreatedAt, id, TableGroupPermissionsColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableGroupPermissionsColumnCreatedBy, comp, And)

	QueryInt64WhereClause(c.GroupID, id, TableGroupPermissionsColumnGroupID, comp, And)

	QueryStringWhereClause(c.PermissionAction, id, TableGroupPermissionsColumnPermissionAction, comp, And)

	QueryStringWhereClause(c.PermissionModule, id, TableGroupPermissionsColumnPermissionModule, comp, And)

	QueryStringWhereClause(c.PermissionSubsystem, id, TableGroupPermissionsColumnPermissionSubsystem, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TableGroupPermissionsColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableGroupPermissionsColumnUpdatedBy, comp, And)

	return nil
}

func (r *GroupPermissionsRepositoryBase) FindQuery(fe *GroupPermissionsFindExpr) (string, []interface{}, error) {
	comp := NewComposer(8)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.created_at, t0.created_by, t0.group_id, t0.permission_action, t0.permission_module, t0.permission_subsystem, t0.updated_at, t0.updated_by")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Fetch {
		buf.WriteString(", t1.created_at, t1.created_by, t1.description, t1.id, t1.name, t1.updated_at, t1.updated_by")
	}

	if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}

	if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinGroup != nil {
		joinClause(comp, fe.JoinGroup.Kind, "charon.group AS t1 ON t0.group_id=t1.id")
		if fe.JoinGroup.On != nil {
			comp.Dirty = true
			if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.On, 1); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinAuthor != nil {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t2 ON t0.created_by=t2.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 2); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinModifier != nil {
		joinClause(comp, fe.JoinModifier.Kind, "charon.user AS t3 ON t0.updated_by=t3.id")
		if fe.JoinModifier.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinModifier.On, 3); err != nil {
				return "", nil, err
			}
		}
	}

	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := GroupPermissionsCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Where != nil {
		if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinModifier.Where, 3); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TableGroupPermissionsColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *GroupPermissionsRepositoryBase) Find(ctx context.Context, fe *GroupPermissionsFindExpr) ([]*GroupPermissionsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*GroupPermissionsEntity
	var props []interface{}
	for rows.Next() {
		var ent GroupPermissionsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinGroup != nil && fe.JoinGroup.Fetch {
			ent.Group = &GroupEntity{}
			if prop, err = ent.Group.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
			ent.Modifier = &UserEntity{}
			if prop, err = ent.Modifier.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *GroupPermissionsRepositoryBase) FindIter(ctx context.Context, fe *GroupPermissionsFindExpr) (*GroupPermissionsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &GroupPermissionsIterator{
		rows: rows,
		cols: []string{"created_at", "created_by", "group_id", "permission_action", "permission_module", "permission_subsystem", "updated_at", "updated_by"},
	}, nil
}
func (r *GroupPermissionsRepositoryBase) FindOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string) (*GroupPermissionsEntity, error) {
	find := NewComposer(8)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableGroupPermissions)
	find.WriteString(" WHERE ")
	find.WriteString(TableGroupPermissionsColumnGroupID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(groupPermissionsGroupID)
	find.WriteString(" AND ")
	find.WriteString(TableGroupPermissionsColumnPermissionSubsystem)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(groupPermissionsPermissionSubsystem)
	find.WriteString(" AND ")
	find.WriteString(TableGroupPermissionsColumnPermissionModule)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(groupPermissionsPermissionModule)
	find.WriteString(" AND ")
	find.WriteString(TableGroupPermissionsColumnPermissionAction)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(groupPermissionsPermissionAction)

	var (
		ent GroupPermissionsEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *GroupPermissionsRepositoryBase) UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, p *GroupPermissionsPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(4)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.GroupID.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnGroupID); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.GroupID)
		update.Dirty = true
	}

	if p.PermissionAction.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnPermissionAction); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionAction)
		update.Dirty = true
	}

	if p.PermissionModule.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnPermissionModule); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionModule)
		update.Dirty = true
	}

	if p.PermissionSubsystem.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnPermissionSubsystem); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionSubsystem)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableGroupPermissionsColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("GroupPermissions update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableGroupPermissionsColumnGroupID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(groupPermissionsGroupID)
	update.WriteString(" AND ")
	update.WriteString(TableGroupPermissionsColumnPermissionSubsystem)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(groupPermissionsPermissionSubsystem)
	update.WriteString(" AND ")
	update.WriteString(TableGroupPermissionsColumnPermissionModule)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(groupPermissionsPermissionModule)
	update.WriteString(" AND ")
	update.WriteString(TableGroupPermissionsColumnPermissionAction)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(groupPermissionsPermissionAction)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by")
	}
	return buf.String(), update.Args(), nil
}
func (r *GroupPermissionsRepositoryBase) UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, p *GroupPermissionsPatch) (*GroupPermissionsEntity, error) {
	query, args, err := r.UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent GroupPermissionsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *GroupPermissionsRepositoryBase) UpsertQuery(e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(16)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.CreatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnGroupID); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.GroupID)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionAction); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionAction)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionModule); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionModule)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnPermissionSubsystem); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionSubsystem)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableGroupPermissionsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UpdatedBy)
	upsert.Dirty = true

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.CreatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnCreatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedBy)
			upsert.Dirty = true
		}

		if p.GroupID.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnGroupID); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.GroupID)
			upsert.Dirty = true
		}

		if p.PermissionAction.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnPermissionAction); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionAction)
			upsert.Dirty = true
		}

		if p.PermissionModule.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnPermissionModule); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionModule)
			upsert.Dirty = true
		}

		if p.PermissionSubsystem.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnPermissionSubsystem); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionSubsystem)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

		if p.UpdatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableGroupPermissionsColumnUpdatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedBy)
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("created_at, created_by, group_id, permission_action, permission_module, permission_subsystem, updated_at, updated_by")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *GroupPermissionsRepositoryBase) Upsert(ctx context.Context, e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (*GroupPermissionsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *GroupPermissionsRepositoryBase) Count(ctx context.Context, c *GroupPermissionsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&GroupPermissionsFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},

		JoinGroup:    c.JoinGroup,
		JoinAuthor:   c.JoinAuthor,
		JoinModifier: c.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "GroupPermissions", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

const (
	TableUserPermissions                           = "charon.user_permissions"
	TableUserPermissionsColumnCreatedAt            = "created_at"
	TableUserPermissionsColumnCreatedBy            = "created_by"
	TableUserPermissionsColumnPermissionAction     = "permission_action"
	TableUserPermissionsColumnPermissionModule     = "permission_module"
	TableUserPermissionsColumnPermissionSubsystem  = "permission_subsystem"
	TableUserPermissionsColumnUpdatedAt            = "updated_at"
	TableUserPermissionsColumnUpdatedBy            = "updated_by"
	TableUserPermissionsColumnUserID               = "user_id"
	TableUserPermissionsConstraintUserIDForeignKey = "charon.user_permissions_user_id_fkey"

	TableUserPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey = "charon.user_permissions_subsystem_module_action_fkey"

	TableUserPermissionsConstraintUserIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.user_permissions_user_id_subsystem_module_action_key"

	TableUserPermissionsConstraintCreatedByForeignKey = "charon.user_permissions_created_by_fkey"

	TableUserPermissionsConstraintUpdatedByForeignKey = "charon.user_permissions_updated_by_fkey"
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

// UserPermissionsEntity ...
type UserPermissionsEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
	// PermissionAction ...
	PermissionAction string
	// PermissionModule ...
	PermissionModule string
	// PermissionSubsystem ...
	PermissionSubsystem string
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// UpdatedBy ...
	UpdatedBy ntypes.Int64
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
	if len(cns) == 0 {
		cns = TableUserPermissionsColumns
	}
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
	rows Rows
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

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
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
	cols, err := i.Columns()
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
	CreatedAt           *qtypes.Timestamp
	CreatedBy           *qtypes.Int64
	PermissionAction    *qtypes.String
	PermissionModule    *qtypes.String
	PermissionSubsystem *qtypes.String
	UpdatedAt           *qtypes.Timestamp
	UpdatedBy           *qtypes.Int64
	UserID              *qtypes.Int64
}

type UserPermissionsFindExpr struct {
	Where         *UserPermissionsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       map[string]bool
	JoinUser      *UserJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type UserPermissionsCountExpr struct {
	Where        *UserPermissionsCriteria
	JoinUser     *UserJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserPermissionsJoin struct {
	On, Where    *UserPermissionsCriteria
	Fetch        bool
	Kind         JoinType
	JoinUser     *UserJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserPermissionsPatch struct {
	CreatedAt           pq.NullTime
	CreatedBy           ntypes.Int64
	PermissionAction    ntypes.String
	PermissionModule    ntypes.String
	PermissionSubsystem ntypes.String
	UpdatedAt           pq.NullTime
	UpdatedBy           ntypes.Int64
	UserID              ntypes.Int64
}

func ScanUserPermissionsRows(rows Rows) (entities []*UserPermissionsEntity, err error) {
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
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

type UserPermissionsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserPermissionsRepositoryBase) InsertQuery(e *UserPermissionsEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(8)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.CreatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.CreatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionAction); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionAction)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionModule); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionModule)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionSubsystem); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.PermissionSubsystem)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if insert.Dirty {
			if _, err := insert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := insert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		insert.Add(e.UpdatedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UpdatedBy)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnUserID); err != nil {
		return "", nil, err
	}
	if insert.Dirty {
		if _, err := insert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := insert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	insert.Add(e.UserID)
	insert.Dirty = true

	if columns.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(insert)
		buf.WriteString(") ")
		if read {
			buf.WriteString("RETURNING ")
			if len(r.Columns) > 0 {
				buf.WriteString(strings.Join(r.Columns, ", "))
			} else {
				buf.WriteString("created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}
func (r *UserPermissionsRepositoryBase) Insert(ctx context.Context, e *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "insert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func UserPermissionsCriteriaWhereClause(comp *Composer, c *UserPermissionsCriteria, id int) error {
	QueryTimestampWhereClause(c.CreatedAt, id, TableUserPermissionsColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableUserPermissionsColumnCreatedBy, comp, And)

	QueryStringWhereClause(c.PermissionAction, id, TableUserPermissionsColumnPermissionAction, comp, And)

	QueryStringWhereClause(c.PermissionModule, id, TableUserPermissionsColumnPermissionModule, comp, And)

	QueryStringWhereClause(c.PermissionSubsystem, id, TableUserPermissionsColumnPermissionSubsystem, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TableUserPermissionsColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableUserPermissionsColumnUpdatedBy, comp, And)

	QueryInt64WhereClause(c.UserID, id, TableUserPermissionsColumnUserID, comp, And)

	return nil
}

func (r *UserPermissionsRepositoryBase) FindQuery(fe *UserPermissionsFindExpr) (string, []interface{}, error) {
	comp := NewComposer(8)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.created_at, t0.created_by, t0.permission_action, t0.permission_module, t0.permission_subsystem, t0.updated_at, t0.updated_by, t0.user_id")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinUser != nil && fe.JoinUser.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}

	if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}

	if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinUser != nil {
		joinClause(comp, fe.JoinUser.Kind, "charon.user AS t1 ON t0.user_id=t1.id")
		if fe.JoinUser.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinUser.On, 1); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinAuthor != nil {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t2 ON t0.created_by=t2.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 2); err != nil {
				return "", nil, err
			}
		}
	}

	if fe.JoinModifier != nil {
		joinClause(comp, fe.JoinModifier.Kind, "charon.user AS t3 ON t0.updated_by=t3.id")
		if fe.JoinModifier.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinModifier.On, 3); err != nil {
				return "", nil, err
			}
		}
	}

	if comp.Dirty {
		buf.ReadFrom(comp)
		comp.Dirty = false
	}
	if fe.Where != nil {
		if err := UserPermissionsCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinUser != nil && fe.JoinUser.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinUser.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinModifier.Where, 3); err != nil {
			return "", nil, err
		}
	}
	if comp.Dirty {
		if _, err := buf.WriteString(" WHERE "); err != nil {
			return "", nil, err
		}
		buf.ReadFrom(comp)
	}

	if len(fe.OrderBy) > 0 {
		i := 0
		for cn, asc := range fe.OrderBy {
			for _, tcn := range TableUserPermissionsColumns {
				if cn == tcn {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(cn); err != nil {
						return "", nil, err
					}
					if !asc {
						if _, err := comp.WriteString(" DESC "); err != nil {
							return "", nil, err
						}
					}
					i++
					break
				}
			}
		}
	}
	if fe.Offset > 0 {
		if _, err := comp.WriteString(" OFFSET "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Offset)
	}
	if fe.Limit > 0 {
		if _, err := comp.WriteString(" LIMIT "); err != nil {
			return "", nil, err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		if _, err := comp.WriteString(" "); err != nil {
			return "", nil, err
		}
		comp.Add(fe.Limit)
	}

	buf.ReadFrom(comp)

	return buf.String(), comp.Args(), nil
}

func (r *UserPermissionsRepositoryBase) Find(ctx context.Context, fe *UserPermissionsFindExpr) ([]*UserPermissionsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*UserPermissionsEntity
	var props []interface{}
	for rows.Next() {
		var ent UserPermissionsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinUser != nil && fe.JoinUser.Fetch {
			ent.User = &UserEntity{}
			if prop, err = ent.User.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Fetch {
			ent.Modifier = &UserEntity{}
			if prop, err = ent.Modifier.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		err = rows.Scan(props...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &ent)
	}
	err = rows.Err()
	if r.Log != nil {
		r.Log(err, "UserPermissions", "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (r *UserPermissionsRepositoryBase) FindIter(ctx context.Context, fe *UserPermissionsFindExpr) (*UserPermissionsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "find iter", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &UserPermissionsIterator{
		rows: rows,
		cols: []string{"created_at", "created_by", "permission_action", "permission_module", "permission_subsystem", "updated_at", "updated_by", "user_id"},
	}, nil
}
func (r *UserPermissionsRepositoryBase) FindOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string) (*UserPermissionsEntity, error) {
	find := NewComposer(8)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableUserPermissions)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserPermissionsColumnUserID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userPermissionsUserID)
	find.WriteString(" AND ")
	find.WriteString(TableUserPermissionsColumnPermissionSubsystem)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userPermissionsPermissionSubsystem)
	find.WriteString(" AND ")
	find.WriteString(TableUserPermissionsColumnPermissionModule)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userPermissionsPermissionModule)
	find.WriteString(" AND ")
	find.WriteString(TableUserPermissionsColumnPermissionAction)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(userPermissionsPermissionAction)

	var (
		ent UserPermissionsEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
func (r *UserPermissionsRepositoryBase) UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, p *UserPermissionsPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(4)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedAt)
		update.Dirty = true

	}

	if p.CreatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnCreatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.CreatedBy)
		update.Dirty = true
	}

	if p.PermissionAction.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnPermissionAction); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionAction)
		update.Dirty = true
	}

	if p.PermissionModule.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnPermissionModule); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionModule)
		update.Dirty = true
	}

	if p.PermissionSubsystem.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnPermissionSubsystem); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.PermissionSubsystem)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedAt)
		update.Dirty = true

	} else {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("=NOW()"); err != nil {
			return "", nil, err
		}
		update.Dirty = true
	}

	if p.UpdatedBy.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnUpdatedBy); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UpdatedBy)
		update.Dirty = true
	}

	if p.UserID.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableUserPermissionsColumnUserID); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.UserID)
		update.Dirty = true
	}

	if !update.Dirty {
		return "", nil, errors.New("UserPermissions update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableUserPermissionsColumnUserID)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userPermissionsUserID)
	update.WriteString(" AND ")
	update.WriteString(TableUserPermissionsColumnPermissionSubsystem)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userPermissionsPermissionSubsystem)
	update.WriteString(" AND ")
	update.WriteString(TableUserPermissionsColumnPermissionModule)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userPermissionsPermissionModule)
	update.WriteString(" AND ")
	update.WriteString(TableUserPermissionsColumnPermissionAction)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(userPermissionsPermissionAction)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id")
	}
	return buf.String(), update.Args(), nil
}
func (r *UserPermissionsRepositoryBase) UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, p *UserPermissionsPatch) (*UserPermissionsEntity, error) {
	query, args, err := r.UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent UserPermissionsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "update one by unique", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
func (r *UserPermissionsRepositoryBase) UpsertQuery(e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(16)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserPermissionsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.CreatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnCreatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.CreatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionAction); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionAction)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionModule); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionModule)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnPermissionSubsystem); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.PermissionSubsystem)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if upsert.Dirty {
			if _, err := upsert.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if err := upsert.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		upsert.Add(e.UpdatedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnUpdatedBy); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UpdatedBy)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableUserPermissionsColumnUserID); err != nil {
		return "", nil, err
	}
	if upsert.Dirty {
		if _, err := upsert.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if err := upsert.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	upsert.Add(e.UserID)
	upsert.Dirty = true

	if upsert.Dirty {
		buf.WriteString(" (")
		buf.ReadFrom(columns)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(upsert)
		buf.WriteString(")")
	}
	buf.WriteString(" ON CONFLICT ")
	if len(inf) > 0 {
		upsert.Dirty = false
		if p.CreatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnCreatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedAt)
			upsert.Dirty = true

		}

		if p.CreatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnCreatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.CreatedBy)
			upsert.Dirty = true
		}

		if p.PermissionAction.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnPermissionAction); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionAction)
			upsert.Dirty = true
		}

		if p.PermissionModule.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnPermissionModule); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionModule)
			upsert.Dirty = true
		}

		if p.PermissionSubsystem.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnPermissionSubsystem); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.PermissionSubsystem)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedAt)
			upsert.Dirty = true

		} else {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnUpdatedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("=NOW()"); err != nil {
				return "", nil, err
			}
			upsert.Dirty = true
		}

		if p.UpdatedBy.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnUpdatedBy); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UpdatedBy)
			upsert.Dirty = true
		}

		if p.UserID.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableUserPermissionsColumnUserID); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.UserID)
			upsert.Dirty = true
		}

	}

	if len(inf) > 0 && upsert.Dirty {
		buf.WriteString("(")
		for j, i := range inf {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(i)
		}
		buf.WriteString(")")
		buf.WriteString(" DO UPDATE SET ")
		buf.ReadFrom(upsert)
	} else {
		buf.WriteString(" DO NOTHING ")
	}
	if upsert.Dirty {
		buf.WriteString(" RETURNING ")
		if len(r.Columns) > 0 {
			buf.WriteString(strings.Join(r.Columns, ", "))
		} else {
			buf.WriteString("created_at, created_by, permission_action, permission_module, permission_subsystem, updated_at, updated_by, user_id")
		}
	}
	return buf.String(), upsert.Args(), nil
}
func (r *UserPermissionsRepositoryBase) Upsert(ctx context.Context, e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (*UserPermissionsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "upsert", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *UserPermissionsRepositoryBase) Count(ctx context.Context, c *UserPermissionsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserPermissionsFindExpr{
		Where:   c.Where,
		Columns: []string{"COUNT(*)"},

		JoinUser:     c.JoinUser,
		JoinAuthor:   c.JoinAuthor,
		JoinModifier: c.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.Log != nil {
		r.Log(err, "UserPermissions", "count", query, args...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

const (
	JoinDoNot = iota
	JoinInner
	JoinLeft
	JoinRight
	JoinCross
)

type JoinType int

func (jt JoinType) String() string {
	switch jt {

	case JoinInner:
		return "INNER JOIN"
	case JoinLeft:
		return "LEFT JOIN"
	case JoinRight:
		return "RIGHT JOIN"
	case JoinCross:
		return "CROSS JOIN"
	default:
		return ""
	}
}

// ErrorConstraint returns the error constraint of err if it was produced by the pq library.
// Otherwise, it returns empty string.
func ErrorConstraint(err error) string {
	if err == nil {
		return "nil"
	}
	if pqerr, ok := err.(*pq.Error); ok {
		return pqerr.Constraint
	}

	return fmt.Sprintf("%T", err)
}

type NullInt64Array struct {
	pq.Int64Array
	Valid bool
}

func (n *NullInt64Array) Scan(value interface{}) error {
	if value == nil {
		n.Int64Array, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return n.Int64Array.Scan(value)
}

type NullFloat64Array struct {
	pq.Float64Array
	Valid bool
}

func (n *NullFloat64Array) Scan(value interface{}) error {
	if value == nil {
		n.Float64Array, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return n.Float64Array.Scan(value)
}

type NullBoolArray struct {
	pq.BoolArray
	Valid bool
}

func (n *NullBoolArray) Scan(value interface{}) error {
	if value == nil {
		n.BoolArray, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return n.BoolArray.Scan(value)
}

type NullStringArray struct {
	pq.StringArray
	Valid bool
}

func (n *NullStringArray) Scan(value interface{}) error {
	if value == nil {
		n.StringArray, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return n.StringArray.Scan(value)
}

type NullByteaArray struct {
	pq.ByteaArray
	Valid bool
}

func (n *NullByteaArray) Scan(value interface{}) error {
	if value == nil {
		n.ByteaArray, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return n.ByteaArray.Scan(value)
}

const (
	jsonArraySeparator     = ","
	jsonArrayBeginningChar = "["
	jsonArrayEndChar       = "]"
)

// JSONArrayInt64 is a slice of int64s that implements necessary interfaces.
type JSONArrayInt64 []int64

// Scan satisfy sql.Scanner interface.
func (a *JSONArrayInt64) Scan(src interface{}) error {
	if src == nil {
		if a == nil {
			*a = make(JSONArrayInt64, 0)
		}
		return nil
	}

	var tmp []string
	var srcs string

	switch t := src.(type) {
	case []byte:
		srcs = string(t)
	case string:
		srcs = t
	default:
		return fmt.Errorf("expected slice of bytes or string as a source argument in Scan, not %T", src)
	}

	l := len(srcs)

	if l < 2 {
		return fmt.Errorf("expected to get source argument in format '[1,2,...,N]', but got %s", srcs)
	}

	if l == 2 {
		*a = make(JSONArrayInt64, 0)
		return nil
	}

	if string(srcs[0]) != jsonArrayBeginningChar || string(srcs[l-1]) != jsonArrayEndChar {
		return fmt.Errorf("expected to get source argument in format '[1,2,...,N]', but got %s", srcs)
	}

	tmp = strings.Split(string(srcs[1:l-1]), jsonArraySeparator)
	*a = make(JSONArrayInt64, 0, len(tmp))
	for i, v := range tmp {
		j, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("expected to get source argument in format '[1,2,...,N]', but got %s at index %d", v, i)
		}

		*a = append(*a, j)
	}

	return nil
}

// Value satisfy driver.Valuer interface.
func (a JSONArrayInt64) Value() (driver.Value, error) {
	var (
		buffer bytes.Buffer
		err    error
	)

	if _, err = buffer.WriteString(jsonArrayBeginningChar); err != nil {
		return nil, err
	}

	for i, v := range a {
		if i > 0 {
			if _, err := buffer.WriteString(jsonArraySeparator); err != nil {
				return nil, err
			}
		}
		if _, err := buffer.WriteString(strconv.FormatInt(v, 10)); err != nil {
			return nil, err
		}
	}

	if _, err = buffer.WriteString(jsonArrayEndChar); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// JSONArrayString is a slice of strings that implements necessary interfaces.
type JSONArrayString []string

// Scan satisfy sql.Scanner interface.
func (a *JSONArrayString) Scan(src interface{}) error {
	if src == nil {
		if a == nil {
			*a = make(JSONArrayString, 0)
		}
		return nil
	}

	switch t := src.(type) {
	case []byte:
		return json.Unmarshal(t, a)
	default:
		return fmt.Errorf("expected slice of bytes or string as a source argument in Scan, not %T", src)
	}
}

// Value satisfy driver.Valuer interface.
func (a JSONArrayString) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// JSONArrayFloat64 is a slice of int64s that implements necessary interfaces.
type JSONArrayFloat64 []float64

// Scan satisfy sql.Scanner interface.
func (a *JSONArrayFloat64) Scan(src interface{}) error {
	if src == nil {
		if a == nil {
			*a = make(JSONArrayFloat64, 0)
		}
		return nil
	}

	var tmp []string
	var srcs string

	switch t := src.(type) {
	case []byte:
		srcs = string(t)
	case string:
		srcs = t
	default:
		return fmt.Errorf("expected slice of bytes or string as a source argument in Scan, not %T", src)
	}

	l := len(srcs)

	if l < 2 {
		return fmt.Errorf("expected to get source argument in format '[1.3,2.4,...,N.M]', but got %s", srcs)
	}

	if l == 2 {
		*a = make(JSONArrayFloat64, 0)
		return nil
	}

	if string(srcs[0]) != jsonArrayBeginningChar || string(srcs[l-1]) != jsonArrayEndChar {
		return fmt.Errorf("expected to get source argument in format '[1.3,2.4,...,N.M]', but got %s", srcs)
	}

	tmp = strings.Split(string(srcs[1:l-1]), jsonArraySeparator)
	*a = make(JSONArrayFloat64, 0, len(tmp))
	for i, v := range tmp {
		j, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("expected to get source argument in format '[1.3,2.4,...,N.M]', but got %s at index %d", v, i)
		}

		*a = append(*a, j)
	}

	return nil
}

// Value satisfy driver.Valuer interface.
func (a JSONArrayFloat64) Value() (driver.Value, error) {
	var (
		buffer bytes.Buffer
		err    error
	)

	if _, err = buffer.WriteString(jsonArrayBeginningChar); err != nil {
		return nil, err
	}

	for i, v := range a {
		if i > 0 {
			if _, err := buffer.WriteString(jsonArraySeparator); err != nil {
				return nil, err
			}
		}
		if _, err := buffer.WriteString(strconv.FormatFloat(v, 'f', -1, 64)); err != nil {
			return nil, err
		}
	}

	if _, err = buffer.WriteString(jsonArrayEndChar); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

var (
	// Space is a shorthand composition option that holds space.
	Space = &CompositionOpts{
		Joint: " ",
	}
	// And is a shorthand composition option that holds AND operator.
	And = &CompositionOpts{
		Joint: " AND ",
	}
	// Or is a shorthand composition option that holds OR operator.
	Or = &CompositionOpts{
		Joint: " OR ",
	}
	// Comma is a shorthand composition option that holds comma.
	Comma = &CompositionOpts{
		Joint: ", ",
	}
)

// CompositionOpts is a container for modification that can be applied.
type CompositionOpts struct {
	Joint                           string
	PlaceholderFuncs, SelectorFuncs []string
	PlaceholderCast, SelectorCast   string
	IsJSON                          bool
	IsDynamic                       bool
}

// CompositionWriter is a simple wrapper for WriteComposition function.
type CompositionWriter interface {
	// WriteComposition is a function that allow custom struct type to be used as a part of criteria.
	// It gives possibility to write custom query based on object that implements this interface.
	WriteComposition(string, *Composer, *CompositionOpts) error
}

// Composer holds buffer, arguments and placeholders count.
// In combination with external buffet can be also used to also generate sub-queries.
// To do that simply write buffer to the parent buffer, composer will hold all arguments and remember number of last placeholder.
type Composer struct {
	buf     bytes.Buffer
	args    []interface{}
	counter int
	Dirty   bool
}

// NewComposer allocates new Composer with inner slice of arguments of given size.
func NewComposer(size int64) *Composer {
	return &Composer{
		counter: 1,
		args:    make([]interface{}, 0, size),
	}
}

// WriteString appends the contents of s to the query buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with bytes ErrTooLarge.
func (c *Composer) WriteString(s string) (int, error) {
	return c.buf.WriteString(s)
}

// Write implements io Writer interface.
func (c *Composer) Write(b []byte) (int, error) {
	return c.buf.Write(b)
}

// Read implements io Reader interface.
func (c *Composer) Read(b []byte) (int, error) {
	return c.buf.Read(b)
}

// ResetBuf resets internal buffer.
func (c *Composer) ResetBuf() {
	c.buf.Reset()
}

// String implements fmt Stringer interface.
func (c *Composer) String() string {
	return c.buf.String()
}

// WritePlaceholder writes appropriate placeholder to the query buffer based on current state of the composer.
func (c *Composer) WritePlaceholder() error {
	if _, err := c.buf.WriteString("$"); err != nil {
		return err
	}
	if _, err := c.buf.WriteString(strconv.Itoa(c.counter)); err != nil {
		return err
	}

	c.counter++
	return nil
}

func (c *Composer) WriteAlias(i int) error {
	if i < 0 {
		return nil
	}
	if _, err := c.buf.WriteString("t"); err != nil {
		return err
	}
	if _, err := c.buf.WriteString(strconv.Itoa(i)); err != nil {
		return err
	}
	if _, err := c.buf.WriteString("."); err != nil {
		return err
	}
	return nil
}

// Len returns number of arguments.
func (c *Composer) Len() int {
	return c.counter
}

// Add appends list with new element.
func (c *Composer) Add(arg interface{}) {
	c.args = append(c.args, arg)
}

// Args returns all arguments stored as a slice.
func (c *Composer) Args() []interface{} {
	return c.args
}
func QueryInt64WhereClause(i *qtypes.Int64, id int, sel string, com *Composer, opt *CompositionOpts) (err error) {
	if i == nil || !i.Valid {
		return nil
	}
	if i.Type == qtypes.QueryType_IN {
		if len(i.Values) == 0 {
			return nil
		}
	}
	if i.Type != qtypes.QueryType_BETWEEN {
		if com.Dirty {
			if _, err = com.WriteString(opt.Joint); err != nil {
				return
			}
		}

		if i.Negation {
			switch i.Type {
			case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_OVERLAP, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if _, err = com.WriteString(" NOT "); err != nil {
					return
				}
			}
		}

		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			case qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if !opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			}
		} else {
			for _, sf := range opt.SelectorFuncs {
				if _, err = com.WriteString(sf); err != nil {
					return err
				}
				if _, err = com.WriteString("("); err != nil {
					return err
				}
			}
		}
		if !opt.IsDynamic {
			if err = com.WriteAlias(id); err != nil {
				return err
			}
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString("("); err != nil {
				return
			}
		}
		if _, err := com.WriteString(sel); err != nil {
			return err
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString(")::"); err != nil {
				return
			}
			if _, err = com.WriteString(opt.SelectorCast); err != nil {
				return
			}
		}
		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("))"); err != nil {
						return err
					}
				}
			}
		} else {
			for range opt.SelectorFuncs {
				if _, err = com.WriteString(")"); err != nil {
					return err
				}
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_NULL:
		if i.Negation {
			if _, err = com.WriteString(" IS NOT NULL"); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IS NULL"); err != nil {
				return
			}
		}
		com.Dirty = true
		return nil

	case qtypes.QueryType_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" <> "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" = "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER:
		if i.Negation {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS:
		if i.Negation {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_CONTAINS:
		if _, err = com.WriteString(" @> "); err != nil {
			return
		}
	case qtypes.QueryType_IS_CONTAINED_BY:
		if _, err = com.WriteString(" <@ "); err != nil {
			return
		}
	case qtypes.QueryType_OVERLAP:
		if _, err = com.WriteString(" && "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ANY_ELEMENT:
		if _, err = com.WriteString(" ?| "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ALL_ELEMENTS:
		if _, err = com.WriteString(" ?& "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ELEMENT:
		if _, err = com.WriteString(" ? "); err != nil {
			return
		}
	case qtypes.QueryType_IN:
		if i.Negation {
			if _, err = com.WriteString(" NOT IN ("); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IN ("); err != nil {
				return
			}
		}
		for i, v := range i.Values {
			if i != 0 {
				if _, err = com.WriteString(","); err != nil {
					return
				}
			}
			if err = com.WritePlaceholder(); err != nil {
				return
			}

			com.Add(v)
			com.Dirty = true
		}
		if _, err = com.WriteString(")"); err != nil {
			return
		}
	case qtypes.QueryType_BETWEEN:
		cpy := *i
		cpy.Values = i.Values[:1]
		cpy.Type = qtypes.QueryType_GREATER
		if err := QueryInt64WhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
		cpy.Values = i.Values[1:]
		cpy.Type = qtypes.QueryType_LESS
		if err := QueryInt64WhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
	default:
		return
	}
	if i.Type != qtypes.QueryType_BETWEEN && i.Type != qtypes.QueryType_IN {
		for _, pf := range opt.PlaceholderFuncs {
			if _, err := com.WriteString(pf); err != nil {
				return err
			}
			if _, err := com.WriteString("("); err != nil {
				return err
			}
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		for range opt.PlaceholderFuncs {
			if _, err := com.WriteString(")"); err != nil {
				return err
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS:
		switch opt.IsJSON {
		case true:
			com.Add(JSONArrayInt64(i.Values))
		case false:
			com.Add(pq.Int64Array(i.Values))
		}
	case qtypes.QueryType_OVERLAP:
		com.Add(pq.Int64Array(i.Values))
	case qtypes.QueryType_SUBSTRING:
		com.Add(fmt.Sprintf("%%%d%%", i.Value()))
	case qtypes.QueryType_HAS_PREFIX:
		com.Add(fmt.Sprintf("%d%%", i.Value()))
	case qtypes.QueryType_HAS_SUFFIX:
		com.Add(fmt.Sprintf("%%%d", i.Value()))
	case qtypes.QueryType_IN:
		// already handled
	case qtypes.QueryType_BETWEEN:
		// already handled by recursive call
	default:

		com.Add(i.Value())
	}

	com.Dirty = true
	return nil
}
func QueryFloat64WhereClause(i *qtypes.Float64, id int, sel string, com *Composer, opt *CompositionOpts) (err error) {
	if i == nil || !i.Valid {
		return nil
	}
	if i.Type == qtypes.QueryType_IN {
		if len(i.Values) == 0 {
			return nil
		}
	}
	if i.Type != qtypes.QueryType_BETWEEN {
		if com.Dirty {
			if _, err = com.WriteString(opt.Joint); err != nil {
				return
			}
		}

		if i.Negation {
			switch i.Type {
			case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_OVERLAP, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if _, err = com.WriteString(" NOT "); err != nil {
					return
				}
			}
		}

		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			case qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if !opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			}
		} else {
			for _, sf := range opt.SelectorFuncs {
				if _, err = com.WriteString(sf); err != nil {
					return err
				}
				if _, err = com.WriteString("("); err != nil {
					return err
				}
			}
		}
		if !opt.IsDynamic {
			if err = com.WriteAlias(id); err != nil {
				return err
			}
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString("("); err != nil {
				return
			}
		}
		if _, err := com.WriteString(sel); err != nil {
			return err
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString(")::"); err != nil {
				return
			}
			if _, err = com.WriteString(opt.SelectorCast); err != nil {
				return
			}
		}
		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("))"); err != nil {
						return err
					}
				}
			}
		} else {
			for range opt.SelectorFuncs {
				if _, err = com.WriteString(")"); err != nil {
					return err
				}
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_NULL:
		if i.Negation {
			if _, err = com.WriteString(" IS NOT NULL"); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IS NULL"); err != nil {
				return
			}
		}
		com.Dirty = true
		return nil

	case qtypes.QueryType_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" <> "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" = "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER:
		if i.Negation {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS:
		if i.Negation {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_CONTAINS:
		if _, err = com.WriteString(" @> "); err != nil {
			return
		}
	case qtypes.QueryType_IS_CONTAINED_BY:
		if _, err = com.WriteString(" <@ "); err != nil {
			return
		}
	case qtypes.QueryType_OVERLAP:
		if _, err = com.WriteString(" && "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ANY_ELEMENT:
		if _, err = com.WriteString(" ?| "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ALL_ELEMENTS:
		if _, err = com.WriteString(" ?& "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ELEMENT:
		if _, err = com.WriteString(" ? "); err != nil {
			return
		}
	case qtypes.QueryType_IN:
		if i.Negation {
			if _, err = com.WriteString(" NOT IN ("); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IN ("); err != nil {
				return
			}
		}
		for i, v := range i.Values {
			if i != 0 {
				if _, err = com.WriteString(","); err != nil {
					return
				}
			}
			if err = com.WritePlaceholder(); err != nil {
				return
			}

			com.Add(v)
			com.Dirty = true
		}
		if _, err = com.WriteString(")"); err != nil {
			return
		}
	case qtypes.QueryType_BETWEEN:
		cpy := *i
		cpy.Values = i.Values[:1]
		cpy.Type = qtypes.QueryType_GREATER
		if err := QueryFloat64WhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
		cpy.Values = i.Values[1:]
		cpy.Type = qtypes.QueryType_LESS
		if err := QueryFloat64WhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
	default:
		return
	}
	if i.Type != qtypes.QueryType_BETWEEN && i.Type != qtypes.QueryType_IN {
		for _, pf := range opt.PlaceholderFuncs {
			if _, err := com.WriteString(pf); err != nil {
				return err
			}
			if _, err := com.WriteString("("); err != nil {
				return err
			}
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		for range opt.PlaceholderFuncs {
			if _, err := com.WriteString(")"); err != nil {
				return err
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS:
		switch opt.IsJSON {
		case true:
			com.Add(JSONArrayFloat64(i.Values))
		case false:
			com.Add(pq.Float64Array(i.Values))
		}
	case qtypes.QueryType_OVERLAP:
		com.Add(pq.Float64Array(i.Values))
	case qtypes.QueryType_SUBSTRING:
		com.Add(fmt.Sprintf("%%%g%%", i.Value()))
	case qtypes.QueryType_HAS_PREFIX:
		com.Add(fmt.Sprintf("%g%%", i.Value()))
	case qtypes.QueryType_HAS_SUFFIX:
		com.Add(fmt.Sprintf("%%%g", i.Value()))
	case qtypes.QueryType_IN:
		// already handled
	case qtypes.QueryType_BETWEEN:
		// already handled by recursive call
	default:

		com.Add(i.Value())
	}

	com.Dirty = true
	return nil
}
func QueryStringWhereClause(i *qtypes.String, id int, sel string, com *Composer, opt *CompositionOpts) (err error) {
	if i == nil || !i.Valid {
		return nil
	}
	if i.Type == qtypes.QueryType_IN {
		if len(i.Values) == 0 {
			return nil
		}
	}
	if i.Type != qtypes.QueryType_BETWEEN {
		if com.Dirty {
			if _, err = com.WriteString(opt.Joint); err != nil {
				return
			}
		}

		if i.Negation {
			switch i.Type {
			case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_OVERLAP, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if _, err = com.WriteString(" NOT "); err != nil {
					return
				}
			}
		}

		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			case qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if !opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			}
		} else {
			for _, sf := range opt.SelectorFuncs {
				if _, err = com.WriteString(sf); err != nil {
					return err
				}
				if _, err = com.WriteString("("); err != nil {
					return err
				}
			}
		}
		if !opt.IsDynamic {
			if err = com.WriteAlias(id); err != nil {
				return err
			}
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString("("); err != nil {
				return
			}
		}
		if _, err := com.WriteString(sel); err != nil {
			return err
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString(")::"); err != nil {
				return
			}
			if _, err = com.WriteString(opt.SelectorCast); err != nil {
				return
			}
		}
		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("))"); err != nil {
						return err
					}
				}
			}
		} else {
			for range opt.SelectorFuncs {
				if _, err = com.WriteString(")"); err != nil {
					return err
				}
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_NULL:
		if i.Negation {
			if _, err = com.WriteString(" IS NOT NULL"); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IS NULL"); err != nil {
				return
			}
		}
		com.Dirty = true
		return nil

	case qtypes.QueryType_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" <> "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" = "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER:
		if i.Negation {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS:
		if i.Negation {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_CONTAINS:
		if _, err = com.WriteString(" @> "); err != nil {
			return
		}
	case qtypes.QueryType_IS_CONTAINED_BY:
		if _, err = com.WriteString(" <@ "); err != nil {
			return
		}
	case qtypes.QueryType_OVERLAP:
		if _, err = com.WriteString(" && "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ANY_ELEMENT:
		if _, err = com.WriteString(" ?| "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ALL_ELEMENTS:
		if _, err = com.WriteString(" ?& "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ELEMENT:
		if _, err = com.WriteString(" ? "); err != nil {
			return
		}
	case qtypes.QueryType_SUBSTRING, qtypes.QueryType_HAS_PREFIX, qtypes.QueryType_HAS_SUFFIX:
		if i.Negation {
			if _, err = com.WriteString(" NOT"); err != nil {
				return
			}
		}
		if i.Insensitive {
			if _, err = com.WriteString(" ILIKE "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" LIKE "); err != nil {
				return
			}
		}
	case qtypes.QueryType_IN:
		if i.Negation {
			if _, err = com.WriteString(" NOT IN ("); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IN ("); err != nil {
				return
			}
		}
		for i, v := range i.Values {
			if i != 0 {
				if _, err = com.WriteString(","); err != nil {
					return
				}
			}
			if err = com.WritePlaceholder(); err != nil {
				return
			}

			com.Add(v)
			com.Dirty = true
		}
		if _, err = com.WriteString(")"); err != nil {
			return
		}
	case qtypes.QueryType_BETWEEN:
		cpy := *i
		cpy.Values = i.Values[:1]
		cpy.Type = qtypes.QueryType_GREATER
		if err := QueryStringWhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
		cpy.Values = i.Values[1:]
		cpy.Type = qtypes.QueryType_LESS
		if err := QueryStringWhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
	default:
		return
	}
	if i.Type != qtypes.QueryType_BETWEEN && i.Type != qtypes.QueryType_IN {
		for _, pf := range opt.PlaceholderFuncs {
			if _, err := com.WriteString(pf); err != nil {
				return err
			}
			if _, err := com.WriteString("("); err != nil {
				return err
			}
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		for range opt.PlaceholderFuncs {
			if _, err := com.WriteString(")"); err != nil {
				return err
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS:
		switch opt.IsJSON {
		case true:
			com.Add(JSONArrayString(i.Values))
		case false:
			com.Add(pq.StringArray(i.Values))
		}
	case qtypes.QueryType_OVERLAP:
		com.Add(pq.StringArray(i.Values))
	case qtypes.QueryType_SUBSTRING:
		com.Add(fmt.Sprintf("%%%s%%", i.Value()))
	case qtypes.QueryType_HAS_PREFIX:
		com.Add(fmt.Sprintf("%s%%", i.Value()))
	case qtypes.QueryType_HAS_SUFFIX:
		com.Add(fmt.Sprintf("%%%s", i.Value()))
	case qtypes.QueryType_IN:
		// already handled
	case qtypes.QueryType_BETWEEN:
		// already handled by recursive call
	default:

		com.Add(i.Value())
	}

	com.Dirty = true
	return nil
}
func QueryTimestampWhereClause(i *qtypes.Timestamp, id int, sel string, com *Composer, opt *CompositionOpts) (err error) {
	if i == nil || !i.Valid {
		return nil
	}
	if i.Type == qtypes.QueryType_IN {
		if len(i.Values) == 0 {
			return nil
		}
	}
	if i.Type != qtypes.QueryType_BETWEEN {
		if com.Dirty {
			if _, err = com.WriteString(opt.Joint); err != nil {
				return
			}
		}

		if i.Negation {
			switch i.Type {
			case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_OVERLAP, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if _, err = com.WriteString(" NOT "); err != nil {
					return
				}
			}
		}

		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			case qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS, qtypes.QueryType_HAS_ELEMENT:
				if !opt.IsJSON {
					if _, err = com.WriteString("ARRAY(SELECT jsonb_array_elements_text("); err != nil {
						return err
					}
				}
			}
		} else {
			for _, sf := range opt.SelectorFuncs {
				if _, err = com.WriteString(sf); err != nil {
					return err
				}
				if _, err = com.WriteString("("); err != nil {
					return err
				}
			}
		}
		if !opt.IsDynamic {
			if err = com.WriteAlias(id); err != nil {
				return err
			}
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString("("); err != nil {
				return
			}
		}
		if _, err := com.WriteString(sel); err != nil {
			return err
		}
		if opt.SelectorCast != "" {
			if _, err = com.WriteString(")::"); err != nil {
				return
			}
			if _, err = com.WriteString(opt.SelectorCast); err != nil {
				return
			}
		}
		if len(opt.SelectorFuncs) == 0 {
			switch i.Type {
			case qtypes.QueryType_OVERLAP:
				if opt.IsJSON {
					if _, err = com.WriteString("))"); err != nil {
						return err
					}
				}
			}
		} else {
			for range opt.SelectorFuncs {
				if _, err = com.WriteString(")"); err != nil {
					return err
				}
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_NULL:
		if i.Negation {
			if _, err = com.WriteString(" IS NOT NULL"); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IS NULL"); err != nil {
				return
			}
		}
		com.Dirty = true
		return nil

	case qtypes.QueryType_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" <> "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" = "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER:
		if i.Negation {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		}
	case qtypes.QueryType_GREATER_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS:
		if i.Negation {
			if _, err = com.WriteString(" >= "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" < "); err != nil {
				return
			}
		}
	case qtypes.QueryType_LESS_EQUAL:
		if i.Negation {
			if _, err = com.WriteString(" > "); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" <= "); err != nil {
				return
			}
		}
	case qtypes.QueryType_CONTAINS:
		if _, err = com.WriteString(" @> "); err != nil {
			return
		}
	case qtypes.QueryType_IS_CONTAINED_BY:
		if _, err = com.WriteString(" <@ "); err != nil {
			return
		}
	case qtypes.QueryType_OVERLAP:
		if _, err = com.WriteString(" && "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ANY_ELEMENT:
		if _, err = com.WriteString(" ?| "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ALL_ELEMENTS:
		if _, err = com.WriteString(" ?& "); err != nil {
			return
		}
	case qtypes.QueryType_HAS_ELEMENT:
		if _, err = com.WriteString(" ? "); err != nil {
			return
		}
	case qtypes.QueryType_IN:
		if i.Negation {
			if _, err = com.WriteString(" NOT IN ("); err != nil {
				return
			}
		} else {
			if _, err = com.WriteString(" IN ("); err != nil {
				return
			}
		}
		for i, v := range i.Values {
			if i != 0 {
				if _, err = com.WriteString(","); err != nil {
					return
				}
			}
			if err = com.WritePlaceholder(); err != nil {
				return
			}

			ts, err := ptypes.Timestamp(v)
			if err != nil {
				return err
			}
			com.Add(ts)
			com.Dirty = true
		}
		if _, err = com.WriteString(")"); err != nil {
			return
		}
	case qtypes.QueryType_BETWEEN:
		cpy := *i
		cpy.Values = i.Values[:1]
		cpy.Type = qtypes.QueryType_GREATER
		if err := QueryTimestampWhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
		cpy.Values = i.Values[1:]
		cpy.Type = qtypes.QueryType_LESS
		if err := QueryTimestampWhereClause(&cpy, id, sel, com, opt); err != nil {
			return err
		}
	default:
		return
	}
	if i.Type != qtypes.QueryType_BETWEEN && i.Type != qtypes.QueryType_IN {
		for _, pf := range opt.PlaceholderFuncs {
			if _, err := com.WriteString(pf); err != nil {
				return err
			}
			if _, err := com.WriteString("("); err != nil {
				return err
			}
		}
		if err = com.WritePlaceholder(); err != nil {
			return
		}
		for range opt.PlaceholderFuncs {
			if _, err := com.WriteString(")"); err != nil {
				return err
			}
		}
	}
	switch i.Type {
	case qtypes.QueryType_CONTAINS, qtypes.QueryType_IS_CONTAINED_BY, qtypes.QueryType_HAS_ANY_ELEMENT, qtypes.QueryType_HAS_ALL_ELEMENTS:
		return errors.New("query type not supported for timestamp")
	case qtypes.QueryType_SUBSTRING:
		com.Add(fmt.Sprintf("%%%s%%", i.Value()))
	case qtypes.QueryType_HAS_PREFIX:
		com.Add(fmt.Sprintf("%s%%", i.Value()))
	case qtypes.QueryType_HAS_SUFFIX:
		com.Add(fmt.Sprintf("%%%s", i.Value()))
	case qtypes.QueryType_IN:
		// already handled
	case qtypes.QueryType_BETWEEN:
		// already handled by recursive call
	default:

		ts, err := ptypes.Timestamp(i.Value())
		if err != nil {
			return err
		}
		com.Add(ts)
	}

	com.Dirty = true
	return nil
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

	CONSTRAINT "charon.user_username_key" UNIQUE (username),
	CONSTRAINT "charon.user_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.user_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.group (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	description TEXT,
	id BIGSERIAL,
	name TEXT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,

	CONSTRAINT "charon.group_name_key" UNIQUE (name),
	CONSTRAINT "charon.group_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.group_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

CREATE TABLE IF NOT EXISTS charon.permission (
	action TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	module TEXT NOT NULL,
	subsystem TEXT NOT NULL,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "charon.permission_subsystem_module_action_key" UNIQUE (subsystem, module, action),
	CONSTRAINT "charon.permission_id_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS charon.user_groups (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	group_id BIGINT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	user_id BIGINT NOT NULL,

	CONSTRAINT "charon.user_groups_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_groups_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
	CONSTRAINT "charon.user_groups_user_id_group_id_key" UNIQUE (user_id, group_id),
	CONSTRAINT "charon.user_groups_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_groups_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
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

	CONSTRAINT "charon.group_permissions_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
	CONSTRAINT "charon.group_permissions_subsystem_module_action_fkey" FOREIGN KEY (permission_subsystem, permission_module, permission_action) REFERENCES charon.permission (subsystem, module, action),
	CONSTRAINT "charon.group_permissions_group_id_subsystem_module_action_key" UNIQUE (group_id, permission_subsystem, permission_module, permission_action),
	CONSTRAINT "charon.group_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
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

	CONSTRAINT "charon.user_permissions_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_subsystem_module_action_fkey" FOREIGN KEY (permission_subsystem, permission_module, permission_action) REFERENCES charon.permission (subsystem, module, action),
	CONSTRAINT "charon.user_permissions_user_id_subsystem_module_action_key" UNIQUE (user_id, permission_subsystem, permission_module, permission_action),
	CONSTRAINT "charon.user_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

`
