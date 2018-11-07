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

// LogFunc represents function that can be passed into repository to log query result.
type LogFunc func(err error, ent, fnc, sql string, args ...interface{})

// RetryTransaction can be returned by user defined function when a transaction is rolled back and logic repeated.
var RetryTransaction = errors.New("retry transaction")

func RunInTransaction(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error, attempts int) (err error) {
	for n := 0; n < attempts; n++ {
		if err = func() error {
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				return err
			}

			defer func() {
				if p := recover(); p != nil {
					_ = tx.Rollback()
					panic(p)
				} else if err != nil {
					_ = tx.Rollback()
				}
			}()

			if err = f(tx); err != nil {
				_ = tx.Rollback()
				return err
			}

			return tx.Commit()
		}(); err == RetryTransaction {
			continue
		}
		return err
	}
	return err
}

// Rows ...
type Rows interface {
	io.Closer
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dst ...interface{}) error
}

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

const (
	TableUserConstraintPrimaryKey          = "charon.user_id_pkey"
	TableUserConstraintUsernameUnique      = "charon.user_username_key"
	TableUserConstraintCreatedByForeignKey = "charon.user_created_by_fkey"
	TableUserConstraintUpdatedByForeignKey = "charon.user_updated_by_fkey"
)

const (
	TableUser                        = "charon.user"
	TableUserColumnConfirmationToken = "confirmation_token"
	TableUserColumnCreatedAt         = "created_at"
	TableUserColumnCreatedBy         = "created_by"
	TableUserColumnFirstName         = "first_name"
	TableUserColumnID                = "id"
	TableUserColumnIsActive          = "is_active"
	TableUserColumnIsConfirmed       = "is_confirmed"
	TableUserColumnIsStaff           = "is_staff"
	TableUserColumnIsSuperuser       = "is_superuser"
	TableUserColumnLastLoginAt       = "last_login_at"
	TableUserColumnLastName          = "last_name"
	TableUserColumnPassword          = "password"
	TableUserColumnUpdatedAt         = "updated_at"
	TableUserColumnUpdatedBy         = "updated_by"
	TableUserColumnUsername          = "username"
)

var TableUserColumns = []string{
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

// ScanUserRows helps to scan rows straight to the slice of entities.
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

// UserIterator is not thread safe.
type UserIterator struct {
	rows Rows
	cols []string
	expr *UserFindExpr
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
	var prop []interface{}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type UserCriteria struct {
	ConfirmationToken      []byte
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	FirstName              *qtypes.String
	ID                     *qtypes.Int64
	IsActive               ntypes.Bool
	IsConfirmed            ntypes.Bool
	IsStaff                ntypes.Bool
	IsSuperuser            ntypes.Bool
	LastLoginAt            *qtypes.Timestamp
	LastName               *qtypes.String
	Password               []byte
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	Username               *qtypes.String
	operator               string
	child, sibling, parent *UserCriteria
}

func UserOperand(operator string, operands ...*UserCriteria) *UserCriteria {
	if len(operands) == 0 {
		return &UserCriteria{operator: operator}
	}

	parent := &UserCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func UserOr(operands ...*UserCriteria) *UserCriteria {
	return UserOperand("OR", operands...)
}

func UserAnd(operands ...*UserCriteria) *UserCriteria {
	return UserOperand("AND", operands...)
}

type UserFindExpr struct {
	Where         *UserCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type UserJoin struct {
	On, Where    *UserCriteria
	Fetch        bool
	Kind         JoinType
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserCountExpr struct {
	Where        *UserCriteria
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

type UserRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserRepositoryBase) Tx(tx *sql.Tx) (*UserRepositoryBaseTx, error) {
	return &UserRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *UserRepositoryBase) BeginTx(ctx context.Context) (*UserRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r UserRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *UserRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *UserRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *UserEntity) (*UserEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
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
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "insert", query, args...)
		} else {
			r.Log(err, TableUser, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserRepositoryBase) Insert(ctx context.Context, e *UserEntity) (*UserEntity, error) {
	return r.insert(ctx, nil, e)
}

func UserCriteriaWhereClause(comp *Composer, c *UserCriteria, id int) error {
	if c.child == nil {
		return _UserCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _UserCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _UserCriteriaWhereClause(comp *Composer, c *UserCriteria, id int) error {
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
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t1 ON t0.created_by=t1.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableUserColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *UserRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *UserFindExpr) ([]*UserEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "find", query, args...)
		} else {
			r.Log(err, TableUser, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*UserEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent UserEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableUser, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *UserRepositoryBase) Find(ctx context.Context, fe *UserFindExpr) ([]*UserEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *UserRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *UserFindExpr) (*UserIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "find iter", query, args...)
		} else {
			r.Log(err, TableUser, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &UserIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *UserRepositoryBase) FindIter(ctx context.Context, fe *UserFindExpr) (*UserIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *UserRepositoryBase) findOneByID(ctx context.Context, tx *sql.Tx, pk int64) (*UserEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "find by primary key", find.String(), find.Args()...)
		} else {
			r.Log(err, TableUser, "find by primary key tx", find.String(), find.Args()...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *UserRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*UserEntity, error) {
	return r.findOneByID(ctx, nil, pk)
}

func (r *UserRepositoryBase) findOneByUsername(ctx context.Context, tx *sql.Tx, userUsername string) (*UserEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *UserRepositoryBase) FindOneByUsername(ctx context.Context, userUsername string) (*UserEntity, error) {
	return r.findOneByUsername(ctx, nil, userUsername)
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

func (r *UserRepositoryBase) updateOneByID(ctx context.Context, tx *sql.Tx, pk int64, p *UserPatch) (*UserEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent UserEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "update by primary key", query, args...)
		} else {
			r.Log(err, TableUser, "update by primary key tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *UserRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *UserPatch) (*UserEntity, error) {
	return r.updateOneByID(ctx, nil, pk, p)
}

func (r *UserRepositoryBase) FindOneByIDAndUpdate(ctx context.Context, pk int64, p *UserPatch) (before, after *UserEntity, err error) {
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
	find.WriteString(" FOR UPDATE")
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return
	}
	var (
		oldEnt, newEnt UserEntity
	)
	oldProps, err := oldEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	newProps, err := newEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	tx, err := r.DB.Begin()
	if err != nil {
		return
	}
	err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(oldProps...)
	if r.Log != nil {
		r.Log(err, TableUser, "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.QueryRowContext(ctx, query, args...).Scan(newProps...)
	if r.Log != nil {
		r.Log(err, TableUser, "update by primary key", query, args...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return &oldEnt, &newEnt, nil
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
		return "", nil, errors.New("user update failure, nothing to update")
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

func (r *UserRepositoryBase) updateOneByUsername(ctx context.Context, tx *sql.Tx, userUsername string, p *UserPatch) (*UserEntity, error) {
	query, args, err := r.UpdateOneByUsernameQuery(userUsername, p)
	if err != nil {
		return nil, err
	}
	var ent UserEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "update one by unique", query, args...)
		} else {
			r.Log(err, TableUser, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *UserRepositoryBase) UpdateOneByUsername(ctx context.Context, userUsername string, p *UserPatch) (*UserEntity, error) {
	return r.updateOneByUsername(ctx, nil, userUsername, p)
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

func (r *UserRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *UserEntity, p *UserPatch, inf ...string) (*UserEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
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
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "upsert", query, args...)
		} else {
			r.Log(err, TableUser, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserRepositoryBase) Upsert(ctx context.Context, e *UserEntity, p *UserPatch, inf ...string) (*UserEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *UserRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *UserCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUser, "count", query, args...)
		} else {
			r.Log(err, TableUser, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepositoryBase) Count(ctx context.Context, exp *UserCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

func (r *UserRepositoryBase) deleteOneByID(ctx context.Context, tx *sql.Tx, pk int64) (int64, error) {
	find := NewComposer(15)
	find.WriteString("DELETE FROM ")
	find.WriteString(TableUser)
	find.WriteString(" WHERE ")
	find.WriteString(TableUserColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		err error
		res sql.Result
	)
	if tx == nil {
		res, err = r.DB.ExecContext(ctx, find.String(), find.Args()...)
	} else {
		res, err = tx.ExecContext(ctx, find.String(), find.Args()...)
	}
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *UserRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.deleteOneByID(ctx, nil, pk)
}

type UserRepositoryBaseTx struct {
	base *UserRepositoryBase
	tx   *sql.Tx
}

func (r UserRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r UserRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *UserRepositoryBaseTx) Insert(ctx context.Context, e *UserEntity) (*UserEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *UserRepositoryBaseTx) Find(ctx context.Context, fe *UserFindExpr) ([]*UserEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *UserRepositoryBaseTx) FindIter(ctx context.Context, fe *UserFindExpr) (*UserIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *UserRepositoryBaseTx) FindOneByID(ctx context.Context, pk int64) (*UserEntity, error) {
	return r.base.findOneByID(ctx, r.tx, pk)
}

func (r *UserRepositoryBaseTx) UpdateOneByID(ctx context.Context, pk int64, p *UserPatch) (*UserEntity, error) {
	return r.base.updateOneByID(ctx, r.tx, pk, p)
}

func (r *UserRepositoryBaseTx) UpdateOneByUsername(ctx context.Context, userUsername string, p *UserPatch) (*UserEntity, error) {
	return r.base.updateOneByUsername(ctx, r.tx, userUsername, p)
}

func (r *UserRepositoryBaseTx) Upsert(ctx context.Context, e *UserEntity, p *UserPatch, inf ...string) (*UserEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *UserRepositoryBaseTx) Count(ctx context.Context, exp *UserCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

func (r *UserRepositoryBaseTx) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.base.deleteOneByID(ctx, r.tx, pk)
}

const (
	TableGroupConstraintNameUnique          = "charon.group_name_key"
	TableGroupConstraintPrimaryKey          = "charon.group_id_pkey"
	TableGroupConstraintCreatedByForeignKey = "charon.group_created_by_fkey"
	TableGroupConstraintUpdatedByForeignKey = "charon.group_updated_by_fkey"
)

const (
	TableGroup                  = "charon.group"
	TableGroupColumnCreatedAt   = "created_at"
	TableGroupColumnCreatedBy   = "created_by"
	TableGroupColumnDescription = "description"
	TableGroupColumnID          = "id"
	TableGroupColumnName        = "name"
	TableGroupColumnUpdatedAt   = "updated_at"
	TableGroupColumnUpdatedBy   = "updated_by"
)

var TableGroupColumns = []string{
	TableGroupColumnCreatedAt,
	TableGroupColumnCreatedBy,
	TableGroupColumnDescription,
	TableGroupColumnID,
	TableGroupColumnName,
	TableGroupColumnUpdatedAt,
	TableGroupColumnUpdatedBy,
}

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

// ScanGroupRows helps to scan rows straight to the slice of entities.
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

// GroupIterator is not thread safe.
type GroupIterator struct {
	rows Rows
	cols []string
	expr *GroupFindExpr
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
	var prop []interface{}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type GroupCriteria struct {
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	Description            *qtypes.String
	ID                     *qtypes.Int64
	Name                   *qtypes.String
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	operator               string
	child, sibling, parent *GroupCriteria
}

func GroupOperand(operator string, operands ...*GroupCriteria) *GroupCriteria {
	if len(operands) == 0 {
		return &GroupCriteria{operator: operator}
	}

	parent := &GroupCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func GroupOr(operands ...*GroupCriteria) *GroupCriteria {
	return GroupOperand("OR", operands...)
}

func GroupAnd(operands ...*GroupCriteria) *GroupCriteria {
	return GroupOperand("AND", operands...)
}

type GroupFindExpr struct {
	Where         *GroupCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type GroupJoin struct {
	On, Where    *GroupCriteria
	Fetch        bool
	Kind         JoinType
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupCountExpr struct {
	Where        *GroupCriteria
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

type GroupRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *GroupRepositoryBase) Tx(tx *sql.Tx) (*GroupRepositoryBaseTx, error) {
	return &GroupRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *GroupRepositoryBase) BeginTx(ctx context.Context) (*GroupRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r GroupRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *GroupRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *GroupRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *GroupEntity) (*GroupEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "insert", query, args...)
		} else {
			r.Log(err, TableGroup, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *GroupRepositoryBase) Insert(ctx context.Context, e *GroupEntity) (*GroupEntity, error) {
	return r.insert(ctx, nil, e)
}

func GroupCriteriaWhereClause(comp *Composer, c *GroupCriteria, id int) error {
	if c.child == nil {
		return _GroupCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _GroupCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _GroupCriteriaWhereClause(comp *Composer, c *GroupCriteria, id int) error {
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
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t1 ON t0.created_by=t1.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableGroupColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *GroupRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *GroupFindExpr) ([]*GroupEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "find", query, args...)
		} else {
			r.Log(err, TableGroup, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*GroupEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent GroupEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableGroup, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GroupRepositoryBase) Find(ctx context.Context, fe *GroupFindExpr) ([]*GroupEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *GroupRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *GroupFindExpr) (*GroupIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "find iter", query, args...)
		} else {
			r.Log(err, TableGroup, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &GroupIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *GroupRepositoryBase) FindIter(ctx context.Context, fe *GroupFindExpr) (*GroupIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *GroupRepositoryBase) findOneByID(ctx context.Context, tx *sql.Tx, pk int64) (*GroupEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "find by primary key", find.String(), find.Args()...)
		} else {
			r.Log(err, TableGroup, "find by primary key tx", find.String(), find.Args()...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *GroupRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*GroupEntity, error) {
	return r.findOneByID(ctx, nil, pk)
}

func (r *GroupRepositoryBase) findOneByName(ctx context.Context, tx *sql.Tx, groupName string) (*GroupEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *GroupRepositoryBase) FindOneByName(ctx context.Context, groupName string) (*GroupEntity, error) {
	return r.findOneByName(ctx, nil, groupName)
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

func (r *GroupRepositoryBase) updateOneByID(ctx context.Context, tx *sql.Tx, pk int64, p *GroupPatch) (*GroupEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent GroupEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "update by primary key", query, args...)
		} else {
			r.Log(err, TableGroup, "update by primary key tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *GroupRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *GroupPatch) (*GroupEntity, error) {
	return r.updateOneByID(ctx, nil, pk, p)
}

func (r *GroupRepositoryBase) FindOneByIDAndUpdate(ctx context.Context, pk int64, p *GroupPatch) (before, after *GroupEntity, err error) {
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
	find.WriteString(" FOR UPDATE")
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return
	}
	var (
		oldEnt, newEnt GroupEntity
	)
	oldProps, err := oldEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	newProps, err := newEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	tx, err := r.DB.Begin()
	if err != nil {
		return
	}
	err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(oldProps...)
	if r.Log != nil {
		r.Log(err, TableGroup, "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.QueryRowContext(ctx, query, args...).Scan(newProps...)
	if r.Log != nil {
		r.Log(err, TableGroup, "update by primary key", query, args...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return &oldEnt, &newEnt, nil
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
		return "", nil, errors.New("group update failure, nothing to update")
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

func (r *GroupRepositoryBase) updateOneByName(ctx context.Context, tx *sql.Tx, groupName string, p *GroupPatch) (*GroupEntity, error) {
	query, args, err := r.UpdateOneByNameQuery(groupName, p)
	if err != nil {
		return nil, err
	}
	var ent GroupEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "update one by unique", query, args...)
		} else {
			r.Log(err, TableGroup, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *GroupRepositoryBase) UpdateOneByName(ctx context.Context, groupName string, p *GroupPatch) (*GroupEntity, error) {
	return r.updateOneByName(ctx, nil, groupName, p)
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

func (r *GroupRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *GroupEntity, p *GroupPatch, inf ...string) (*GroupEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.Description,
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "upsert", query, args...)
		} else {
			r.Log(err, TableGroup, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *GroupRepositoryBase) Upsert(ctx context.Context, e *GroupEntity, p *GroupPatch, inf ...string) (*GroupEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *GroupRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *GroupCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&GroupFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroup, "count", query, args...)
		} else {
			r.Log(err, TableGroup, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupRepositoryBase) Count(ctx context.Context, exp *GroupCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

func (r *GroupRepositoryBase) deleteOneByID(ctx context.Context, tx *sql.Tx, pk int64) (int64, error) {
	find := NewComposer(7)
	find.WriteString("DELETE FROM ")
	find.WriteString(TableGroup)
	find.WriteString(" WHERE ")
	find.WriteString(TableGroupColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		err error
		res sql.Result
	)
	if tx == nil {
		res, err = r.DB.ExecContext(ctx, find.String(), find.Args()...)
	} else {
		res, err = tx.ExecContext(ctx, find.String(), find.Args()...)
	}
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *GroupRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.deleteOneByID(ctx, nil, pk)
}

type GroupRepositoryBaseTx struct {
	base *GroupRepositoryBase
	tx   *sql.Tx
}

func (r GroupRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r GroupRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *GroupRepositoryBaseTx) Insert(ctx context.Context, e *GroupEntity) (*GroupEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *GroupRepositoryBaseTx) Find(ctx context.Context, fe *GroupFindExpr) ([]*GroupEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *GroupRepositoryBaseTx) FindIter(ctx context.Context, fe *GroupFindExpr) (*GroupIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *GroupRepositoryBaseTx) FindOneByID(ctx context.Context, pk int64) (*GroupEntity, error) {
	return r.base.findOneByID(ctx, r.tx, pk)
}

func (r *GroupRepositoryBaseTx) UpdateOneByID(ctx context.Context, pk int64, p *GroupPatch) (*GroupEntity, error) {
	return r.base.updateOneByID(ctx, r.tx, pk, p)
}

func (r *GroupRepositoryBaseTx) UpdateOneByName(ctx context.Context, groupName string, p *GroupPatch) (*GroupEntity, error) {
	return r.base.updateOneByName(ctx, r.tx, groupName, p)
}

func (r *GroupRepositoryBaseTx) Upsert(ctx context.Context, e *GroupEntity, p *GroupPatch, inf ...string) (*GroupEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *GroupRepositoryBaseTx) Count(ctx context.Context, exp *GroupCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

func (r *GroupRepositoryBaseTx) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.base.deleteOneByID(ctx, r.tx, pk)
}

const (
	TablePermissionConstraintSubsystemModuleActionUnique = "charon.permission_subsystem_module_action_key"
	TablePermissionConstraintPrimaryKey                  = "charon.permission_id_pkey"
)

const (
	TablePermission                = "charon.permission"
	TablePermissionColumnAction    = "action"
	TablePermissionColumnCreatedAt = "created_at"
	TablePermissionColumnID        = "id"
	TablePermissionColumnModule    = "module"
	TablePermissionColumnSubsystem = "subsystem"
	TablePermissionColumnUpdatedAt = "updated_at"
)

var TablePermissionColumns = []string{
	TablePermissionColumnAction,
	TablePermissionColumnCreatedAt,
	TablePermissionColumnID,
	TablePermissionColumnModule,
	TablePermissionColumnSubsystem,
	TablePermissionColumnUpdatedAt,
}

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

// ScanPermissionRows helps to scan rows straight to the slice of entities.
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

// PermissionIterator is not thread safe.
type PermissionIterator struct {
	rows Rows
	cols []string
	expr *PermissionFindExpr
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
	Action                 *qtypes.String
	CreatedAt              *qtypes.Timestamp
	ID                     *qtypes.Int64
	Module                 *qtypes.String
	Subsystem              *qtypes.String
	UpdatedAt              *qtypes.Timestamp
	operator               string
	child, sibling, parent *PermissionCriteria
}

func PermissionOperand(operator string, operands ...*PermissionCriteria) *PermissionCriteria {
	if len(operands) == 0 {
		return &PermissionCriteria{operator: operator}
	}

	parent := &PermissionCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func PermissionOr(operands ...*PermissionCriteria) *PermissionCriteria {
	return PermissionOperand("OR", operands...)
}

func PermissionAnd(operands ...*PermissionCriteria) *PermissionCriteria {
	return PermissionOperand("AND", operands...)
}

type PermissionFindExpr struct {
	Where         *PermissionCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
}

type PermissionJoin struct {
	On, Where *PermissionCriteria
	Fetch     bool
	Kind      JoinType
}

type PermissionCountExpr struct {
	Where *PermissionCriteria
}

type PermissionPatch struct {
	Action    ntypes.String
	CreatedAt pq.NullTime
	Module    ntypes.String
	Subsystem ntypes.String
	UpdatedAt pq.NullTime
}

type PermissionRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *PermissionRepositoryBase) Tx(tx *sql.Tx) (*PermissionRepositoryBaseTx, error) {
	return &PermissionRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *PermissionRepositoryBase) BeginTx(ctx context.Context) (*PermissionRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r PermissionRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *PermissionRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *PermissionRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *PermissionEntity) (*PermissionEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "insert", query, args...)
		} else {
			r.Log(err, TablePermission, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PermissionRepositoryBase) Insert(ctx context.Context, e *PermissionEntity) (*PermissionEntity, error) {
	return r.insert(ctx, nil, e)
}

func PermissionCriteriaWhereClause(comp *Composer, c *PermissionCriteria, id int) error {
	if c.child == nil {
		return _PermissionCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _PermissionCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _PermissionCriteriaWhereClause(comp *Composer, c *PermissionCriteria, id int) error {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TablePermissionColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *PermissionRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *PermissionFindExpr) ([]*PermissionEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "find", query, args...)
		} else {
			r.Log(err, TablePermission, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*PermissionEntity
		props    []interface{}
	)
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
		r.Log(err, TablePermission, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *PermissionRepositoryBase) Find(ctx context.Context, fe *PermissionFindExpr) ([]*PermissionEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *PermissionRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *PermissionFindExpr) (*PermissionIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "find iter", query, args...)
		} else {
			r.Log(err, TablePermission, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &PermissionIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *PermissionRepositoryBase) FindIter(ctx context.Context, fe *PermissionFindExpr) (*PermissionIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *PermissionRepositoryBase) findOneByID(ctx context.Context, tx *sql.Tx, pk int64) (*PermissionEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "find by primary key", find.String(), find.Args()...)
		} else {
			r.Log(err, TablePermission, "find by primary key tx", find.String(), find.Args()...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *PermissionRepositoryBase) FindOneByID(ctx context.Context, pk int64) (*PermissionEntity, error) {
	return r.findOneByID(ctx, nil, pk)
}

func (r *PermissionRepositoryBase) findOneBySubsystemAndModuleAndAction(ctx context.Context, tx *sql.Tx, permissionSubsystem string, permissionModule string, permissionAction string) (*PermissionEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *PermissionRepositoryBase) FindOneBySubsystemAndModuleAndAction(ctx context.Context, permissionSubsystem string, permissionModule string, permissionAction string) (*PermissionEntity, error) {
	return r.findOneBySubsystemAndModuleAndAction(ctx, nil, permissionSubsystem, permissionModule, permissionAction)
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

func (r *PermissionRepositoryBase) updateOneByID(ctx context.Context, tx *sql.Tx, pk int64, p *PermissionPatch) (*PermissionEntity, error) {
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return nil, err
	}
	var ent PermissionEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(props...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "update by primary key", query, args...)
		} else {
			r.Log(err, TablePermission, "update by primary key tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *PermissionRepositoryBase) UpdateOneByID(ctx context.Context, pk int64, p *PermissionPatch) (*PermissionEntity, error) {
	return r.updateOneByID(ctx, nil, pk, p)
}

func (r *PermissionRepositoryBase) FindOneByIDAndUpdate(ctx context.Context, pk int64, p *PermissionPatch) (before, after *PermissionEntity, err error) {
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
	find.WriteString(" FOR UPDATE")
	query, args, err := r.UpdateOneByIDQuery(pk, p)
	if err != nil {
		return
	}
	var (
		oldEnt, newEnt PermissionEntity
	)
	oldProps, err := oldEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	newProps, err := newEnt.Props(r.Columns...)
	if err != nil {
		return
	}
	tx, err := r.DB.Begin()
	if err != nil {
		return
	}
	err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(oldProps...)
	if r.Log != nil {
		r.Log(err, TablePermission, "find by primary key", find.String(), find.Args()...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.QueryRowContext(ctx, query, args...).Scan(newProps...)
	if r.Log != nil {
		r.Log(err, TablePermission, "update by primary key", query, args...)
	}
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return &oldEnt, &newEnt, nil
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
		return "", nil, errors.New("permission update failure, nothing to update")
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

func (r *PermissionRepositoryBase) updateOneBySubsystemAndModuleAndAction(ctx context.Context, tx *sql.Tx, permissionSubsystem string, permissionModule string, permissionAction string, p *PermissionPatch) (*PermissionEntity, error) {
	query, args, err := r.UpdateOneBySubsystemAndModuleAndActionQuery(permissionSubsystem, permissionModule, permissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent PermissionEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "update one by unique", query, args...)
		} else {
			r.Log(err, TablePermission, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *PermissionRepositoryBase) UpdateOneBySubsystemAndModuleAndAction(ctx context.Context, permissionSubsystem string, permissionModule string, permissionAction string, p *PermissionPatch) (*PermissionEntity, error) {
	return r.updateOneBySubsystemAndModuleAndAction(ctx, nil, permissionSubsystem, permissionModule, permissionAction, p)
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

func (r *PermissionRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *PermissionEntity, p *PermissionPatch, inf ...string) (*PermissionEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.Action,
		&e.CreatedAt,
		&e.ID,
		&e.Module,
		&e.Subsystem,
		&e.UpdatedAt,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "upsert", query, args...)
		} else {
			r.Log(err, TablePermission, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PermissionRepositoryBase) Upsert(ctx context.Context, e *PermissionEntity, p *PermissionPatch, inf ...string) (*PermissionEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *PermissionRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *PermissionCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&PermissionFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TablePermission, "count", query, args...)
		} else {
			r.Log(err, TablePermission, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PermissionRepositoryBase) Count(ctx context.Context, exp *PermissionCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

func (r *PermissionRepositoryBase) deleteOneByID(ctx context.Context, tx *sql.Tx, pk int64) (int64, error) {
	find := NewComposer(6)
	find.WriteString("DELETE FROM ")
	find.WriteString(TablePermission)
	find.WriteString(" WHERE ")
	find.WriteString(TablePermissionColumnID)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(pk)
	var (
		err error
		res sql.Result
	)
	if tx == nil {
		res, err = r.DB.ExecContext(ctx, find.String(), find.Args()...)
	} else {
		res, err = tx.ExecContext(ctx, find.String(), find.Args()...)
	}
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *PermissionRepositoryBase) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.deleteOneByID(ctx, nil, pk)
}

type PermissionRepositoryBaseTx struct {
	base *PermissionRepositoryBase
	tx   *sql.Tx
}

func (r PermissionRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r PermissionRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *PermissionRepositoryBaseTx) Insert(ctx context.Context, e *PermissionEntity) (*PermissionEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *PermissionRepositoryBaseTx) Find(ctx context.Context, fe *PermissionFindExpr) ([]*PermissionEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *PermissionRepositoryBaseTx) FindIter(ctx context.Context, fe *PermissionFindExpr) (*PermissionIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *PermissionRepositoryBaseTx) FindOneByID(ctx context.Context, pk int64) (*PermissionEntity, error) {
	return r.base.findOneByID(ctx, r.tx, pk)
}

func (r *PermissionRepositoryBaseTx) UpdateOneByID(ctx context.Context, pk int64, p *PermissionPatch) (*PermissionEntity, error) {
	return r.base.updateOneByID(ctx, r.tx, pk, p)
}

func (r *PermissionRepositoryBaseTx) UpdateOneBySubsystemAndModuleAndAction(ctx context.Context, permissionSubsystem string, permissionModule string, permissionAction string, p *PermissionPatch) (*PermissionEntity, error) {
	return r.base.updateOneBySubsystemAndModuleAndAction(ctx, r.tx, permissionSubsystem, permissionModule, permissionAction, p)
}

func (r *PermissionRepositoryBaseTx) Upsert(ctx context.Context, e *PermissionEntity, p *PermissionPatch, inf ...string) (*PermissionEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *PermissionRepositoryBaseTx) Count(ctx context.Context, exp *PermissionCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

func (r *PermissionRepositoryBaseTx) DeleteOneByID(ctx context.Context, pk int64) (int64, error) {
	return r.base.deleteOneByID(ctx, r.tx, pk)
}

const (
	TableUserGroupsConstraintUserIDForeignKey    = "charon.user_groups_user_id_fkey"
	TableUserGroupsConstraintGroupIDForeignKey   = "charon.user_groups_group_id_fkey"
	TableUserGroupsConstraintUserIDGroupIDUnique = "charon.user_groups_user_id_group_id_key"
	TableUserGroupsConstraintCreatedByForeignKey = "charon.user_groups_created_by_fkey"
	TableUserGroupsConstraintUpdatedByForeignKey = "charon.user_groups_updated_by_fkey"
)

const (
	TableUserGroups                = "charon.user_groups"
	TableUserGroupsColumnCreatedAt = "created_at"
	TableUserGroupsColumnCreatedBy = "created_by"
	TableUserGroupsColumnGroupID   = "group_id"
	TableUserGroupsColumnUpdatedAt = "updated_at"
	TableUserGroupsColumnUpdatedBy = "updated_by"
	TableUserGroupsColumnUserID    = "user_id"
)

var TableUserGroupsColumns = []string{
	TableUserGroupsColumnCreatedAt,
	TableUserGroupsColumnCreatedBy,
	TableUserGroupsColumnGroupID,
	TableUserGroupsColumnUpdatedAt,
	TableUserGroupsColumnUpdatedBy,
	TableUserGroupsColumnUserID,
}

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

// ScanUserGroupsRows helps to scan rows straight to the slice of entities.
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

// UserGroupsIterator is not thread safe.
type UserGroupsIterator struct {
	rows Rows
	cols []string
	expr *UserGroupsFindExpr
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
	var prop []interface{}
	if i.expr.JoinUser != nil && i.expr.JoinUser.Kind.Actionable() && i.expr.JoinUser.Fetch {
		ent.User = &UserEntity{}
		if prop, err = ent.User.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinGroup != nil && i.expr.JoinGroup.Kind.Actionable() && i.expr.JoinGroup.Fetch {
		ent.Group = &GroupEntity{}
		if prop, err = ent.Group.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type UserGroupsCriteria struct {
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	GroupID                *qtypes.Int64
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	UserID                 *qtypes.Int64
	operator               string
	child, sibling, parent *UserGroupsCriteria
}

func UserGroupsOperand(operator string, operands ...*UserGroupsCriteria) *UserGroupsCriteria {
	if len(operands) == 0 {
		return &UserGroupsCriteria{operator: operator}
	}

	parent := &UserGroupsCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func UserGroupsOr(operands ...*UserGroupsCriteria) *UserGroupsCriteria {
	return UserGroupsOperand("OR", operands...)
}

func UserGroupsAnd(operands ...*UserGroupsCriteria) *UserGroupsCriteria {
	return UserGroupsOperand("AND", operands...)
}

type UserGroupsFindExpr struct {
	Where         *UserGroupsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinUser      *UserJoin
	JoinGroup     *GroupJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
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

type UserGroupsCountExpr struct {
	Where        *UserGroupsCriteria
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

type UserGroupsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserGroupsRepositoryBase) Tx(tx *sql.Tx) (*UserGroupsRepositoryBaseTx, error) {
	return &UserGroupsRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *UserGroupsRepositoryBase) BeginTx(ctx context.Context) (*UserGroupsRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r UserGroupsRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *UserGroupsRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *UserGroupsRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *UserGroupsEntity) (*UserGroupsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "insert", query, args...)
		} else {
			r.Log(err, TableUserGroups, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserGroupsRepositoryBase) Insert(ctx context.Context, e *UserGroupsEntity) (*UserGroupsEntity, error) {
	return r.insert(ctx, nil, e)
}

func UserGroupsCriteriaWhereClause(comp *Composer, c *UserGroupsCriteria, id int) error {
	if c.child == nil {
		return _UserGroupsCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _UserGroupsCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _UserGroupsCriteriaWhereClause(comp *Composer, c *UserGroupsCriteria, id int) error {
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
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Fetch {
		buf.WriteString(", t2.created_at, t2.created_by, t2.description, t2.id, t2.name, t2.updated_at, t2.updated_by")
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t4.confirmation_token, t4.created_at, t4.created_by, t4.first_name, t4.id, t4.is_active, t4.is_confirmed, t4.is_staff, t4.is_superuser, t4.last_login_at, t4.last_name, t4.password, t4.updated_at, t4.updated_by, t4.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() {
		joinClause(comp, fe.JoinUser.Kind, "charon.user AS t1 ON t0.user_id=t1.id")
		if fe.JoinUser.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinUser.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() {
		joinClause(comp, fe.JoinGroup.Kind, "charon.group AS t2 ON t0.group_id=t2.id")
		if fe.JoinGroup.On != nil {
			comp.Dirty = true
			if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.On, 2); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t3 ON t0.created_by=t3.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 3); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinUser.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Where != nil {
		if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 3); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableUserGroupsColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *UserGroupsRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *UserGroupsFindExpr) ([]*UserGroupsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "find", query, args...)
		} else {
			r.Log(err, TableUserGroups, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*UserGroupsEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent UserGroupsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
			ent.User = &UserEntity{}
			if prop, err = ent.User.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Fetch {
			ent.Group = &GroupEntity{}
			if prop, err = ent.Group.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableUserGroups, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *UserGroupsRepositoryBase) Find(ctx context.Context, fe *UserGroupsFindExpr) ([]*UserGroupsEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *UserGroupsRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *UserGroupsFindExpr) (*UserGroupsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "find iter", query, args...)
		} else {
			r.Log(err, TableUserGroups, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &UserGroupsIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *UserGroupsRepositoryBase) FindIter(ctx context.Context, fe *UserGroupsFindExpr) (*UserGroupsIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *UserGroupsRepositoryBase) findOneByUserIDAndGroupID(ctx context.Context, tx *sql.Tx, userGroupsUserID int64, userGroupsGroupID int64) (*UserGroupsEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *UserGroupsRepositoryBase) FindOneByUserIDAndGroupID(ctx context.Context, userGroupsUserID int64, userGroupsGroupID int64) (*UserGroupsEntity, error) {
	return r.findOneByUserIDAndGroupID(ctx, nil, userGroupsUserID, userGroupsGroupID)
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
		return "", nil, errors.New("user_groups update failure, nothing to update")
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

func (r *UserGroupsRepositoryBase) updateOneByUserIDAndGroupID(ctx context.Context, tx *sql.Tx, userGroupsUserID int64, userGroupsGroupID int64, p *UserGroupsPatch) (*UserGroupsEntity, error) {
	query, args, err := r.UpdateOneByUserIDAndGroupIDQuery(userGroupsUserID, userGroupsGroupID, p)
	if err != nil {
		return nil, err
	}
	var ent UserGroupsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "update one by unique", query, args...)
		} else {
			r.Log(err, TableUserGroups, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *UserGroupsRepositoryBase) UpdateOneByUserIDAndGroupID(ctx context.Context, userGroupsUserID int64, userGroupsGroupID int64, p *UserGroupsPatch) (*UserGroupsEntity, error) {
	return r.updateOneByUserIDAndGroupID(ctx, nil, userGroupsUserID, userGroupsGroupID, p)
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

func (r *UserGroupsRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (*UserGroupsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "upsert", query, args...)
		} else {
			r.Log(err, TableUserGroups, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserGroupsRepositoryBase) Upsert(ctx context.Context, e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (*UserGroupsEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *UserGroupsRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *UserGroupsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserGroupsFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinUser:     exp.JoinUser,
		JoinGroup:    exp.JoinGroup,
		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserGroups, "count", query, args...)
		} else {
			r.Log(err, TableUserGroups, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserGroupsRepositoryBase) Count(ctx context.Context, exp *UserGroupsCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

type UserGroupsRepositoryBaseTx struct {
	base *UserGroupsRepositoryBase
	tx   *sql.Tx
}

func (r UserGroupsRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r UserGroupsRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *UserGroupsRepositoryBaseTx) Insert(ctx context.Context, e *UserGroupsEntity) (*UserGroupsEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *UserGroupsRepositoryBaseTx) Find(ctx context.Context, fe *UserGroupsFindExpr) ([]*UserGroupsEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *UserGroupsRepositoryBaseTx) FindIter(ctx context.Context, fe *UserGroupsFindExpr) (*UserGroupsIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *UserGroupsRepositoryBaseTx) UpdateOneByUserIDAndGroupID(ctx context.Context, userGroupsUserID int64, userGroupsGroupID int64, p *UserGroupsPatch) (*UserGroupsEntity, error) {
	return r.base.updateOneByUserIDAndGroupID(ctx, r.tx, userGroupsUserID, userGroupsGroupID, p)
}

func (r *UserGroupsRepositoryBaseTx) Upsert(ctx context.Context, e *UserGroupsEntity, p *UserGroupsPatch, inf ...string) (*UserGroupsEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *UserGroupsRepositoryBaseTx) Count(ctx context.Context, exp *UserGroupsCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

const (
	TableGroupPermissionsConstraintGroupIDForeignKey                                                = "charon.group_permissions_group_id_fkey"
	TableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey    = "charon.group_permissions_subsystem_module_action_fkey"
	TableGroupPermissionsConstraintGroupIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.group_permissions_group_id_subsystem_module_action_key"
	TableGroupPermissionsConstraintCreatedByForeignKey                                              = "charon.group_permissions_created_by_fkey"
	TableGroupPermissionsConstraintUpdatedByForeignKey                                              = "charon.group_permissions_updated_by_fkey"
)

const (
	TableGroupPermissions                          = "charon.group_permissions"
	TableGroupPermissionsColumnCreatedAt           = "created_at"
	TableGroupPermissionsColumnCreatedBy           = "created_by"
	TableGroupPermissionsColumnGroupID             = "group_id"
	TableGroupPermissionsColumnPermissionAction    = "permission_action"
	TableGroupPermissionsColumnPermissionModule    = "permission_module"
	TableGroupPermissionsColumnPermissionSubsystem = "permission_subsystem"
	TableGroupPermissionsColumnUpdatedAt           = "updated_at"
	TableGroupPermissionsColumnUpdatedBy           = "updated_by"
)

var TableGroupPermissionsColumns = []string{
	TableGroupPermissionsColumnCreatedAt,
	TableGroupPermissionsColumnCreatedBy,
	TableGroupPermissionsColumnGroupID,
	TableGroupPermissionsColumnPermissionAction,
	TableGroupPermissionsColumnPermissionModule,
	TableGroupPermissionsColumnPermissionSubsystem,
	TableGroupPermissionsColumnUpdatedAt,
	TableGroupPermissionsColumnUpdatedBy,
}

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

// ScanGroupPermissionsRows helps to scan rows straight to the slice of entities.
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

// GroupPermissionsIterator is not thread safe.
type GroupPermissionsIterator struct {
	rows Rows
	cols []string
	expr *GroupPermissionsFindExpr
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
	var prop []interface{}
	if i.expr.JoinGroup != nil && i.expr.JoinGroup.Kind.Actionable() && i.expr.JoinGroup.Fetch {
		ent.Group = &GroupEntity{}
		if prop, err = ent.Group.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type GroupPermissionsCriteria struct {
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	GroupID                *qtypes.Int64
	PermissionAction       *qtypes.String
	PermissionModule       *qtypes.String
	PermissionSubsystem    *qtypes.String
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	operator               string
	child, sibling, parent *GroupPermissionsCriteria
}

func GroupPermissionsOperand(operator string, operands ...*GroupPermissionsCriteria) *GroupPermissionsCriteria {
	if len(operands) == 0 {
		return &GroupPermissionsCriteria{operator: operator}
	}

	parent := &GroupPermissionsCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func GroupPermissionsOr(operands ...*GroupPermissionsCriteria) *GroupPermissionsCriteria {
	return GroupPermissionsOperand("OR", operands...)
}

func GroupPermissionsAnd(operands ...*GroupPermissionsCriteria) *GroupPermissionsCriteria {
	return GroupPermissionsOperand("AND", operands...)
}

type GroupPermissionsFindExpr struct {
	Where         *GroupPermissionsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinGroup     *GroupJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type GroupPermissionsJoin struct {
	On, Where    *GroupPermissionsCriteria
	Fetch        bool
	Kind         JoinType
	JoinGroup    *GroupJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type GroupPermissionsCountExpr struct {
	Where        *GroupPermissionsCriteria
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

type GroupPermissionsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *GroupPermissionsRepositoryBase) Tx(tx *sql.Tx) (*GroupPermissionsRepositoryBaseTx, error) {
	return &GroupPermissionsRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *GroupPermissionsRepositoryBase) BeginTx(ctx context.Context) (*GroupPermissionsRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r GroupPermissionsRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *GroupPermissionsRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *GroupPermissionsRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "insert", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *GroupPermissionsRepositoryBase) Insert(ctx context.Context, e *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	return r.insert(ctx, nil, e)
}

func GroupPermissionsCriteriaWhereClause(comp *Composer, c *GroupPermissionsCriteria, id int) error {
	if c.child == nil {
		return _GroupPermissionsCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _GroupPermissionsCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _GroupPermissionsCriteriaWhereClause(comp *Composer, c *GroupPermissionsCriteria, id int) error {
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
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Fetch {
		buf.WriteString(", t1.created_at, t1.created_by, t1.description, t1.id, t1.name, t1.updated_at, t1.updated_by")
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() {
		joinClause(comp, fe.JoinGroup.Kind, "charon.group AS t1 ON t0.group_id=t1.id")
		if fe.JoinGroup.On != nil {
			comp.Dirty = true
			if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t2 ON t0.created_by=t2.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 2); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
	if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Where != nil {
		if err := GroupCriteriaWhereClause(comp, fe.JoinGroup.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableGroupPermissionsColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *GroupPermissionsRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *GroupPermissionsFindExpr) ([]*GroupPermissionsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "find", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*GroupPermissionsEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent GroupPermissionsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinGroup != nil && fe.JoinGroup.Kind.Actionable() && fe.JoinGroup.Fetch {
			ent.Group = &GroupEntity{}
			if prop, err = ent.Group.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableGroupPermissions, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GroupPermissionsRepositoryBase) Find(ctx context.Context, fe *GroupPermissionsFindExpr) ([]*GroupPermissionsEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *GroupPermissionsRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *GroupPermissionsFindExpr) (*GroupPermissionsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "find iter", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &GroupPermissionsIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *GroupPermissionsRepositoryBase) FindIter(ctx context.Context, fe *GroupPermissionsFindExpr) (*GroupPermissionsIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *GroupPermissionsRepositoryBase) findOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, tx *sql.Tx, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string) (*GroupPermissionsEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *GroupPermissionsRepositoryBase) FindOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string) (*GroupPermissionsEntity, error) {
	return r.findOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, nil, groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction)
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
		return "", nil, errors.New("group_permissions update failure, nothing to update")
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

func (r *GroupPermissionsRepositoryBase) updateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, tx *sql.Tx, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, p *GroupPermissionsPatch) (*GroupPermissionsEntity, error) {
	query, args, err := r.UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent GroupPermissionsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "update one by unique", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *GroupPermissionsRepositoryBase) UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, p *GroupPermissionsPatch) (*GroupPermissionsEntity, error) {
	return r.updateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, nil, groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction, p)
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

func (r *GroupPermissionsRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (*GroupPermissionsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "upsert", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *GroupPermissionsRepositoryBase) Upsert(ctx context.Context, e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (*GroupPermissionsEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *GroupPermissionsRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *GroupPermissionsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&GroupPermissionsFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinGroup:    exp.JoinGroup,
		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableGroupPermissions, "count", query, args...)
		} else {
			r.Log(err, TableGroupPermissions, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupPermissionsRepositoryBase) Count(ctx context.Context, exp *GroupPermissionsCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

type GroupPermissionsRepositoryBaseTx struct {
	base *GroupPermissionsRepositoryBase
	tx   *sql.Tx
}

func (r GroupPermissionsRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r GroupPermissionsRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *GroupPermissionsRepositoryBaseTx) Insert(ctx context.Context, e *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *GroupPermissionsRepositoryBaseTx) Find(ctx context.Context, fe *GroupPermissionsFindExpr) ([]*GroupPermissionsEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *GroupPermissionsRepositoryBaseTx) FindIter(ctx context.Context, fe *GroupPermissionsFindExpr) (*GroupPermissionsIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *GroupPermissionsRepositoryBaseTx) UpdateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, groupPermissionsGroupID int64, groupPermissionsPermissionSubsystem string, groupPermissionsPermissionModule string, groupPermissionsPermissionAction string, p *GroupPermissionsPatch) (*GroupPermissionsEntity, error) {
	return r.base.updateOneByGroupIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, r.tx, groupPermissionsGroupID, groupPermissionsPermissionSubsystem, groupPermissionsPermissionModule, groupPermissionsPermissionAction, p)
}

func (r *GroupPermissionsRepositoryBaseTx) Upsert(ctx context.Context, e *GroupPermissionsEntity, p *GroupPermissionsPatch, inf ...string) (*GroupPermissionsEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *GroupPermissionsRepositoryBaseTx) Count(ctx context.Context, exp *GroupPermissionsCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

const (
	TableUserPermissionsConstraintUserIDForeignKey                                                = "charon.user_permissions_user_id_fkey"
	TableUserPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey   = "charon.user_permissions_subsystem_module_action_fkey"
	TableUserPermissionsConstraintUserIDPermissionSubsystemPermissionModulePermissionActionUnique = "charon.user_permissions_user_id_subsystem_module_action_key"
	TableUserPermissionsConstraintCreatedByForeignKey                                             = "charon.user_permissions_created_by_fkey"
	TableUserPermissionsConstraintUpdatedByForeignKey                                             = "charon.user_permissions_updated_by_fkey"
)

const (
	TableUserPermissions                          = "charon.user_permissions"
	TableUserPermissionsColumnCreatedAt           = "created_at"
	TableUserPermissionsColumnCreatedBy           = "created_by"
	TableUserPermissionsColumnPermissionAction    = "permission_action"
	TableUserPermissionsColumnPermissionModule    = "permission_module"
	TableUserPermissionsColumnPermissionSubsystem = "permission_subsystem"
	TableUserPermissionsColumnUpdatedAt           = "updated_at"
	TableUserPermissionsColumnUpdatedBy           = "updated_by"
	TableUserPermissionsColumnUserID              = "user_id"
)

var TableUserPermissionsColumns = []string{
	TableUserPermissionsColumnCreatedAt,
	TableUserPermissionsColumnCreatedBy,
	TableUserPermissionsColumnPermissionAction,
	TableUserPermissionsColumnPermissionModule,
	TableUserPermissionsColumnPermissionSubsystem,
	TableUserPermissionsColumnUpdatedAt,
	TableUserPermissionsColumnUpdatedBy,
	TableUserPermissionsColumnUserID,
}

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

// ScanUserPermissionsRows helps to scan rows straight to the slice of entities.
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

// UserPermissionsIterator is not thread safe.
type UserPermissionsIterator struct {
	rows Rows
	cols []string
	expr *UserPermissionsFindExpr
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
	var prop []interface{}
	if i.expr.JoinUser != nil && i.expr.JoinUser.Kind.Actionable() && i.expr.JoinUser.Fetch {
		ent.User = &UserEntity{}
		if prop, err = ent.User.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type UserPermissionsCriteria struct {
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	PermissionAction       *qtypes.String
	PermissionModule       *qtypes.String
	PermissionSubsystem    *qtypes.String
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	UserID                 *qtypes.Int64
	operator               string
	child, sibling, parent *UserPermissionsCriteria
}

func UserPermissionsOperand(operator string, operands ...*UserPermissionsCriteria) *UserPermissionsCriteria {
	if len(operands) == 0 {
		return &UserPermissionsCriteria{operator: operator}
	}

	parent := &UserPermissionsCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func UserPermissionsOr(operands ...*UserPermissionsCriteria) *UserPermissionsCriteria {
	return UserPermissionsOperand("OR", operands...)
}

func UserPermissionsAnd(operands ...*UserPermissionsCriteria) *UserPermissionsCriteria {
	return UserPermissionsOperand("AND", operands...)
}

type UserPermissionsFindExpr struct {
	Where         *UserPermissionsCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinUser      *UserJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type UserPermissionsJoin struct {
	On, Where    *UserPermissionsCriteria
	Fetch        bool
	Kind         JoinType
	JoinUser     *UserJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type UserPermissionsCountExpr struct {
	Where        *UserPermissionsCriteria
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

type UserPermissionsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *UserPermissionsRepositoryBase) Tx(tx *sql.Tx) (*UserPermissionsRepositoryBaseTx, error) {
	return &UserPermissionsRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *UserPermissionsRepositoryBase) BeginTx(ctx context.Context) (*UserPermissionsRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r UserPermissionsRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *UserPermissionsRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
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

func (r *UserPermissionsRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "insert", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserPermissionsRepositoryBase) Insert(ctx context.Context, e *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	return r.insert(ctx, nil, e)
}

func UserPermissionsCriteriaWhereClause(comp *Composer, c *UserPermissionsCriteria, id int) error {
	if c.child == nil {
		return _UserPermissionsCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _UserPermissionsCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _UserPermissionsCriteriaWhereClause(comp *Composer, c *UserPermissionsCriteria, id int) error {
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
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() {
		joinClause(comp, fe.JoinUser.Kind, "charon.user AS t1 ON t0.user_id=t1.id")
		if fe.JoinUser.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinUser.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t2 ON t0.created_by=t2.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 2); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinUser.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableUserPermissionsColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *UserPermissionsRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *UserPermissionsFindExpr) ([]*UserPermissionsEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "find", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*UserPermissionsEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent UserPermissionsEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
			ent.User = &UserEntity{}
			if prop, err = ent.User.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableUserPermissions, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *UserPermissionsRepositoryBase) Find(ctx context.Context, fe *UserPermissionsFindExpr) ([]*UserPermissionsEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *UserPermissionsRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *UserPermissionsFindExpr) (*UserPermissionsIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "find iter", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &UserPermissionsIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *UserPermissionsRepositoryBase) FindIter(ctx context.Context, fe *UserPermissionsFindExpr) (*UserPermissionsIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *UserPermissionsRepositoryBase) findOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, tx *sql.Tx, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string) (*UserPermissionsEntity, error) {
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
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *UserPermissionsRepositoryBase) FindOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string) (*UserPermissionsEntity, error) {
	return r.findOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, nil, userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction)
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
		return "", nil, errors.New("user_permissions update failure, nothing to update")
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

func (r *UserPermissionsRepositoryBase) updateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, tx *sql.Tx, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, p *UserPermissionsPatch) (*UserPermissionsEntity, error) {
	query, args, err := r.UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionActionQuery(userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction, p)
	if err != nil {
		return nil, err
	}
	var ent UserPermissionsEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "update one by unique", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *UserPermissionsRepositoryBase) UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, p *UserPermissionsPatch) (*UserPermissionsEntity, error) {
	return r.updateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, nil, userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction, p)
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

func (r *UserPermissionsRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (*UserPermissionsEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionAction,
		&e.PermissionModule,
		&e.PermissionSubsystem,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "upsert", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserPermissionsRepositoryBase) Upsert(ctx context.Context, e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (*UserPermissionsEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *UserPermissionsRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *UserPermissionsCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&UserPermissionsFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinUser:     exp.JoinUser,
		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableUserPermissions, "count", query, args...)
		} else {
			r.Log(err, TableUserPermissions, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserPermissionsRepositoryBase) Count(ctx context.Context, exp *UserPermissionsCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

type UserPermissionsRepositoryBaseTx struct {
	base *UserPermissionsRepositoryBase
	tx   *sql.Tx
}

func (r UserPermissionsRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r UserPermissionsRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *UserPermissionsRepositoryBaseTx) Insert(ctx context.Context, e *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *UserPermissionsRepositoryBaseTx) Find(ctx context.Context, fe *UserPermissionsFindExpr) ([]*UserPermissionsEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *UserPermissionsRepositoryBaseTx) FindIter(ctx context.Context, fe *UserPermissionsFindExpr) (*UserPermissionsIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *UserPermissionsRepositoryBaseTx) UpdateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx context.Context, userPermissionsUserID int64, userPermissionsPermissionSubsystem string, userPermissionsPermissionModule string, userPermissionsPermissionAction string, p *UserPermissionsPatch) (*UserPermissionsEntity, error) {
	return r.base.updateOneByUserIDAndPermissionSubsystemAndPermissionModuleAndPermissionAction(ctx, r.tx, userPermissionsUserID, userPermissionsPermissionSubsystem, userPermissionsPermissionModule, userPermissionsPermissionAction, p)
}

func (r *UserPermissionsRepositoryBaseTx) Upsert(ctx context.Context, e *UserPermissionsEntity, p *UserPermissionsPatch, inf ...string) (*UserPermissionsEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *UserPermissionsRepositoryBaseTx) Count(ctx context.Context, exp *UserPermissionsCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

const (
	TableRefreshTokenConstraintTokenUnique         = "charon.refresh_token_token_key"
	TableRefreshTokenConstraintUserIDForeignKey    = "charon.refresh_token_user_id_fkey"
	TableRefreshTokenConstraintCreatedByForeignKey = "charon.refresh_token_created_by_fkey"
	TableRefreshTokenConstraintUpdatedByForeignKey = "charon.refresh_token_updated_by_fkey"
)

const (
	TableRefreshToken                 = "charon.refresh_token"
	TableRefreshTokenColumnCreatedAt  = "created_at"
	TableRefreshTokenColumnCreatedBy  = "created_by"
	TableRefreshTokenColumnExpireAt   = "expire_at"
	TableRefreshTokenColumnLastUsedAt = "last_used_at"
	TableRefreshTokenColumnNotes      = "notes"
	TableRefreshTokenColumnRevoked    = "revoked"
	TableRefreshTokenColumnToken      = "token"
	TableRefreshTokenColumnUpdatedAt  = "updated_at"
	TableRefreshTokenColumnUpdatedBy  = "updated_by"
	TableRefreshTokenColumnUserID     = "user_id"
)

var TableRefreshTokenColumns = []string{
	TableRefreshTokenColumnCreatedAt,
	TableRefreshTokenColumnCreatedBy,
	TableRefreshTokenColumnExpireAt,
	TableRefreshTokenColumnLastUsedAt,
	TableRefreshTokenColumnNotes,
	TableRefreshTokenColumnRevoked,
	TableRefreshTokenColumnToken,
	TableRefreshTokenColumnUpdatedAt,
	TableRefreshTokenColumnUpdatedBy,
	TableRefreshTokenColumnUserID,
}

// RefreshTokenEntity ...
type RefreshTokenEntity struct {
	// CreatedAt ...
	CreatedAt time.Time
	// CreatedBy ...
	CreatedBy ntypes.Int64
	// ExpireAt ...
	ExpireAt pq.NullTime
	// LastUsedAt ...
	LastUsedAt pq.NullTime
	// Notes ...
	Notes ntypes.String
	// Revoked ...
	Revoked bool
	// Token ...
	Token string
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

func (e *RefreshTokenEntity) Prop(cn string) (interface{}, bool) {
	switch cn {

	case TableRefreshTokenColumnCreatedAt:
		return &e.CreatedAt, true
	case TableRefreshTokenColumnCreatedBy:
		return &e.CreatedBy, true
	case TableRefreshTokenColumnExpireAt:
		return &e.ExpireAt, true
	case TableRefreshTokenColumnLastUsedAt:
		return &e.LastUsedAt, true
	case TableRefreshTokenColumnNotes:
		return &e.Notes, true
	case TableRefreshTokenColumnRevoked:
		return &e.Revoked, true
	case TableRefreshTokenColumnToken:
		return &e.Token, true
	case TableRefreshTokenColumnUpdatedAt:
		return &e.UpdatedAt, true
	case TableRefreshTokenColumnUpdatedBy:
		return &e.UpdatedBy, true
	case TableRefreshTokenColumnUserID:
		return &e.UserID, true
	default:
		return nil, false
	}
}

func (e *RefreshTokenEntity) Props(cns ...string) ([]interface{}, error) {
	if len(cns) == 0 {
		cns = TableRefreshTokenColumns
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

// ScanRefreshTokenRows helps to scan rows straight to the slice of entities.
func ScanRefreshTokenRows(rows Rows) (entities []*RefreshTokenEntity, err error) {
	for rows.Next() {
		var ent RefreshTokenEntity
		err = rows.Scan(
			&ent.CreatedAt,
			&ent.CreatedBy,
			&ent.ExpireAt,
			&ent.LastUsedAt,
			&ent.Notes,
			&ent.Revoked,
			&ent.Token,
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

// RefreshTokenIterator is not thread safe.
type RefreshTokenIterator struct {
	rows Rows
	cols []string
	expr *RefreshTokenFindExpr
}

func (i *RefreshTokenIterator) Next() bool {
	return i.rows.Next()
}

func (i *RefreshTokenIterator) Close() error {
	return i.rows.Close()
}

func (i *RefreshTokenIterator) Err() error {
	return i.rows.Err()
}

// Columns is wrapper around sql.Rows.Columns method, that also cache output inside iterator.
func (i *RefreshTokenIterator) Columns() ([]string, error) {
	if i.cols == nil {
		cols, err := i.rows.Columns()
		if err != nil {
			return nil, err
		}
		i.cols = cols
	}
	return i.cols, nil
}

// Ent is wrapper around RefreshToken method that makes iterator more generic.
func (i *RefreshTokenIterator) Ent() (interface{}, error) {
	return i.RefreshToken()
}

func (i *RefreshTokenIterator) RefreshToken() (*RefreshTokenEntity, error) {
	var ent RefreshTokenEntity
	cols, err := i.Columns()
	if err != nil {
		return nil, err
	}

	props, err := ent.Props(cols...)
	if err != nil {
		return nil, err
	}
	var prop []interface{}
	if i.expr.JoinUser != nil && i.expr.JoinUser.Kind.Actionable() && i.expr.JoinUser.Fetch {
		ent.User = &UserEntity{}
		if prop, err = ent.User.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinAuthor != nil && i.expr.JoinAuthor.Kind.Actionable() && i.expr.JoinAuthor.Fetch {
		ent.Author = &UserEntity{}
		if prop, err = ent.Author.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if i.expr.JoinModifier != nil && i.expr.JoinModifier.Kind.Actionable() && i.expr.JoinModifier.Fetch {
		ent.Modifier = &UserEntity{}
		if prop, err = ent.Modifier.Props(); err != nil {
			return nil, err
		}
		props = append(props, prop...)
	}
	if err := i.rows.Scan(props...); err != nil {
		return nil, err
	}
	return &ent, nil
}

type RefreshTokenCriteria struct {
	CreatedAt              *qtypes.Timestamp
	CreatedBy              *qtypes.Int64
	ExpireAt               *qtypes.Timestamp
	LastUsedAt             *qtypes.Timestamp
	Notes                  *qtypes.String
	Revoked                ntypes.Bool
	Token                  *qtypes.String
	UpdatedAt              *qtypes.Timestamp
	UpdatedBy              *qtypes.Int64
	UserID                 *qtypes.Int64
	operator               string
	child, sibling, parent *RefreshTokenCriteria
}

func RefreshTokenOperand(operator string, operands ...*RefreshTokenCriteria) *RefreshTokenCriteria {
	if len(operands) == 0 {
		return &RefreshTokenCriteria{operator: operator}
	}

	parent := &RefreshTokenCriteria{
		operator: operator,
		child:    operands[0],
	}

	for i := 0; i < len(operands); i++ {
		if i < len(operands)-1 {
			operands[i].sibling = operands[i+1]
		}
		operands[i].parent = parent
	}

	return parent
}

func RefreshTokenOr(operands ...*RefreshTokenCriteria) *RefreshTokenCriteria {
	return RefreshTokenOperand("OR", operands...)
}

func RefreshTokenAnd(operands ...*RefreshTokenCriteria) *RefreshTokenCriteria {
	return RefreshTokenOperand("AND", operands...)
}

type RefreshTokenFindExpr struct {
	Where         *RefreshTokenCriteria
	Offset, Limit int64
	Columns       []string
	OrderBy       []RowOrder
	JoinUser      *UserJoin
	JoinAuthor    *UserJoin
	JoinModifier  *UserJoin
}

type RefreshTokenJoin struct {
	On, Where    *RefreshTokenCriteria
	Fetch        bool
	Kind         JoinType
	JoinUser     *UserJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type RefreshTokenCountExpr struct {
	Where        *RefreshTokenCriteria
	JoinUser     *UserJoin
	JoinAuthor   *UserJoin
	JoinModifier *UserJoin
}

type RefreshTokenPatch struct {
	CreatedAt  pq.NullTime
	CreatedBy  ntypes.Int64
	ExpireAt   pq.NullTime
	LastUsedAt pq.NullTime
	Notes      ntypes.String
	Revoked    ntypes.Bool
	Token      ntypes.String
	UpdatedAt  pq.NullTime
	UpdatedBy  ntypes.Int64
	UserID     ntypes.Int64
}

type RefreshTokenRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Log     LogFunc
}

func (r *RefreshTokenRepositoryBase) Tx(tx *sql.Tx) (*RefreshTokenRepositoryBaseTx, error) {
	return &RefreshTokenRepositoryBaseTx{
		base: r,
		tx:   tx,
	}, nil
}

func (r *RefreshTokenRepositoryBase) BeginTx(ctx context.Context) (*RefreshTokenRepositoryBaseTx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Tx(tx)
}

func (r RefreshTokenRepositoryBase) RunInTransaction(ctx context.Context, fn func(rtx *RefreshTokenRepositoryBaseTx) error, attempts int) (err error) {
	return RunInTransaction(ctx, r.DB, func(tx *sql.Tx) error {
		rtx, err := r.Tx(tx)
		if err != nil {
			return err
		}
		return fn(rtx)
	}, attempts)
}

func (r *RefreshTokenRepositoryBase) InsertQuery(e *RefreshTokenEntity, read bool) (string, []interface{}, error) {
	insert := NewComposer(10)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnCreatedAt); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnCreatedBy); err != nil {
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

	if e.ExpireAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnExpireAt); err != nil {
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
		insert.Add(e.ExpireAt)
		insert.Dirty = true
	}

	if e.LastUsedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnLastUsedAt); err != nil {
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
		insert.Add(e.LastUsedAt)
		insert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnNotes); err != nil {
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
	insert.Add(e.Notes)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnRevoked); err != nil {
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
	insert.Add(e.Revoked)
	insert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnToken); err != nil {
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
	insert.Add(e.Token)
	insert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnUpdatedBy); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnUserID); err != nil {
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
				buf.WriteString("created_at, created_by, expire_at, last_used_at, notes, revoked, token, updated_at, updated_by, user_id")
			}
		}
	}
	return buf.String(), insert.Args(), nil
}

func (r *RefreshTokenRepositoryBase) insert(ctx context.Context, tx *sql.Tx, e *RefreshTokenEntity) (*RefreshTokenEntity, error) {
	query, args, err := r.InsertQuery(e, true)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.ExpireAt,
		&e.LastUsedAt,
		&e.Notes,
		&e.Revoked,
		&e.Token,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "insert", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "insert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *RefreshTokenRepositoryBase) Insert(ctx context.Context, e *RefreshTokenEntity) (*RefreshTokenEntity, error) {
	return r.insert(ctx, nil, e)
}

func RefreshTokenCriteriaWhereClause(comp *Composer, c *RefreshTokenCriteria, id int) error {
	if c.child == nil {
		return _RefreshTokenCriteriaWhereClause(comp, c, id)
	}
	node := c
	sibling := false
	for {
		if !sibling {
			if node.child != nil {
				if node.parent != nil {
					comp.WriteString("(")
				}
				node = node.child
				continue
			} else {
				comp.Dirty = false
				comp.WriteString("(")
				if err := _RefreshTokenCriteriaWhereClause(comp, node, id); err != nil {
					return err
				}
				comp.WriteString(")")
			}
		}
		if node.sibling != nil {
			sibling = false
			comp.WriteString(" ")
			comp.WriteString(node.parent.operator)
			comp.WriteString(" ")
			node = node.sibling
			continue
		}
		if node.parent != nil {
			sibling = true
			if node.parent.parent != nil {
				comp.WriteString(")")
			}
			node = node.parent
			continue
		}

		break
	}
	return nil
}

func _RefreshTokenCriteriaWhereClause(comp *Composer, c *RefreshTokenCriteria, id int) error {
	QueryTimestampWhereClause(c.CreatedAt, id, TableRefreshTokenColumnCreatedAt, comp, And)

	QueryInt64WhereClause(c.CreatedBy, id, TableRefreshTokenColumnCreatedBy, comp, And)

	QueryTimestampWhereClause(c.ExpireAt, id, TableRefreshTokenColumnExpireAt, comp, And)

	QueryTimestampWhereClause(c.LastUsedAt, id, TableRefreshTokenColumnLastUsedAt, comp, And)

	QueryStringWhereClause(c.Notes, id, TableRefreshTokenColumnNotes, comp, And)

	if c.Revoked.Valid {
		if comp.Dirty {
			if _, err := comp.WriteString(" AND "); err != nil {
				return err
			}
		}
		if err := comp.WriteAlias(id); err != nil {
			return err
		}
		if _, err := comp.WriteString(TableRefreshTokenColumnRevoked); err != nil {
			return err
		}
		if _, err := comp.WriteString("="); err != nil {
			return err
		}
		if err := comp.WritePlaceholder(); err != nil {
			return err
		}
		comp.Add(c.Revoked)
		comp.Dirty = true
	}

	QueryStringWhereClause(c.Token, id, TableRefreshTokenColumnToken, comp, And)

	QueryTimestampWhereClause(c.UpdatedAt, id, TableRefreshTokenColumnUpdatedAt, comp, And)

	QueryInt64WhereClause(c.UpdatedBy, id, TableRefreshTokenColumnUpdatedBy, comp, And)

	QueryInt64WhereClause(c.UserID, id, TableRefreshTokenColumnUserID, comp, And)

	return nil
}

func (r *RefreshTokenRepositoryBase) FindQuery(fe *RefreshTokenFindExpr) (string, []interface{}, error) {
	comp := NewComposer(10)
	buf := bytes.NewBufferString("SELECT ")
	if len(fe.Columns) == 0 {
		buf.WriteString("t0.created_at, t0.created_by, t0.expire_at, t0.last_used_at, t0.notes, t0.revoked, t0.token, t0.updated_at, t0.updated_by, t0.user_id")
	} else {
		buf.WriteString(strings.Join(fe.Columns, ", "))
	}
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
		buf.WriteString(", t1.confirmation_token, t1.created_at, t1.created_by, t1.first_name, t1.id, t1.is_active, t1.is_confirmed, t1.is_staff, t1.is_superuser, t1.last_login_at, t1.last_name, t1.password, t1.updated_at, t1.updated_by, t1.username")
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
		buf.WriteString(", t2.confirmation_token, t2.created_at, t2.created_by, t2.first_name, t2.id, t2.is_active, t2.is_confirmed, t2.is_staff, t2.is_superuser, t2.last_login_at, t2.last_name, t2.password, t2.updated_at, t2.updated_by, t2.username")
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
		buf.WriteString(", t3.confirmation_token, t3.created_at, t3.created_by, t3.first_name, t3.id, t3.is_active, t3.is_confirmed, t3.is_staff, t3.is_superuser, t3.last_login_at, t3.last_name, t3.password, t3.updated_at, t3.updated_by, t3.username")
	}
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" AS t0")
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() {
		joinClause(comp, fe.JoinUser.Kind, "charon.user AS t1 ON t0.user_id=t1.id")
		if fe.JoinUser.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinUser.On, 1); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() {
		joinClause(comp, fe.JoinAuthor.Kind, "charon.user AS t2 ON t0.created_by=t2.id")
		if fe.JoinAuthor.On != nil {
			comp.Dirty = true
			if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.On, 2); err != nil {
				return "", nil, err
			}
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() {
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
		if err := RefreshTokenCriteriaWhereClause(comp, fe.Where, 0); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinUser.Where, 1); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Where != nil {
		if err := UserCriteriaWhereClause(comp, fe.JoinAuthor.Where, 2); err != nil {
			return "", nil, err
		}
	}
	if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Where != nil {
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
		for _, order := range fe.OrderBy {
			for _, columnName := range TableRefreshTokenColumns {
				if order.Name == columnName {
					if i == 0 {
						comp.WriteString(" ORDER BY ")
					}
					if i > 0 {
						if _, err := comp.WriteString(", "); err != nil {
							return "", nil, err
						}
					}
					if _, err := comp.WriteString(order.Name); err != nil {
						return "", nil, err
					}
					if order.Descending {
						if _, err := comp.WriteString(" DESC"); err != nil {
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

func (r *RefreshTokenRepositoryBase) find(ctx context.Context, tx *sql.Tx, fe *RefreshTokenFindExpr) ([]*RefreshTokenEntity, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "find", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "find tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		entities []*RefreshTokenEntity
		props    []interface{}
	)
	for rows.Next() {
		var ent RefreshTokenEntity
		if props, err = ent.Props(); err != nil {
			return nil, err
		}
		var prop []interface{}
		if fe.JoinUser != nil && fe.JoinUser.Kind.Actionable() && fe.JoinUser.Fetch {
			ent.User = &UserEntity{}
			if prop, err = ent.User.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinAuthor != nil && fe.JoinAuthor.Kind.Actionable() && fe.JoinAuthor.Fetch {
			ent.Author = &UserEntity{}
			if prop, err = ent.Author.Props(); err != nil {
				return nil, err
			}
			props = append(props, prop...)
		}
		if fe.JoinModifier != nil && fe.JoinModifier.Kind.Actionable() && fe.JoinModifier.Fetch {
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
		r.Log(err, TableRefreshToken, "find", query, args...)
	}
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *RefreshTokenRepositoryBase) Find(ctx context.Context, fe *RefreshTokenFindExpr) ([]*RefreshTokenEntity, error) {
	return r.find(ctx, nil, fe)
}

func (r *RefreshTokenRepositoryBase) findIter(ctx context.Context, tx *sql.Tx, fe *RefreshTokenFindExpr) (*RefreshTokenIterator, error) {
	query, args, err := r.FindQuery(fe)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx == nil {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query, args...)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "find iter", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "find iter tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &RefreshTokenIterator{
		rows: rows,
		expr: fe,
		cols: fe.Columns,
	}, nil
}

func (r *RefreshTokenRepositoryBase) FindIter(ctx context.Context, fe *RefreshTokenFindExpr) (*RefreshTokenIterator, error) {
	return r.findIter(ctx, nil, fe)
}

func (r *RefreshTokenRepositoryBase) findOneByToken(ctx context.Context, tx *sql.Tx, refreshTokenToken string) (*RefreshTokenEntity, error) {
	find := NewComposer(10)
	find.WriteString("SELECT ")
	if len(r.Columns) == 0 {
		find.WriteString("created_at, created_by, expire_at, last_used_at, notes, revoked, token, updated_at, updated_by, user_id")
	} else {
		find.WriteString(strings.Join(r.Columns, ", "))
	}
	find.WriteString(" FROM ")
	find.WriteString(TableRefreshToken)
	find.WriteString(" WHERE ")
	find.WriteString(TableRefreshTokenColumnToken)
	find.WriteString("=")
	find.WritePlaceholder()
	find.Add(refreshTokenToken)

	var (
		ent RefreshTokenEntity
	)
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	} else {
		err = tx.QueryRowContext(ctx, find.String(), find.Args()...).Scan(props...)
	}
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *RefreshTokenRepositoryBase) FindOneByToken(ctx context.Context, refreshTokenToken string) (*RefreshTokenEntity, error) {
	return r.findOneByToken(ctx, nil, refreshTokenToken)
}

func (r *RefreshTokenRepositoryBase) UpdateOneByTokenQuery(refreshTokenToken string, p *RefreshTokenPatch) (string, []interface{}, error) {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(r.Table)
	update := NewComposer(1)
	if p.CreatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnCreatedAt); err != nil {
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
		if _, err := update.WriteString(TableRefreshTokenColumnCreatedBy); err != nil {
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

	if p.ExpireAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnExpireAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.ExpireAt)
		update.Dirty = true

	}
	if p.LastUsedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnLastUsedAt); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.LastUsedAt)
		update.Dirty = true

	}
	if p.Notes.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnNotes); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Notes)
		update.Dirty = true
	}

	if p.Revoked.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnRevoked); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Revoked)
		update.Dirty = true
	}

	if p.Token.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnToken); err != nil {
			return "", nil, err
		}
		if _, err := update.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := update.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		update.Add(p.Token)
		update.Dirty = true
	}

	if p.UpdatedAt.Valid {
		if update.Dirty {
			if _, err := update.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := update.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
		if _, err := update.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
		if _, err := update.WriteString(TableRefreshTokenColumnUpdatedBy); err != nil {
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
		if _, err := update.WriteString(TableRefreshTokenColumnUserID); err != nil {
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
		return "", nil, errors.New("refresh_token update failure, nothing to update")
	}
	buf.WriteString(" SET ")
	buf.ReadFrom(update)
	buf.WriteString(" WHERE ")
	update.WriteString(TableRefreshTokenColumnToken)
	update.WriteString("=")
	update.WritePlaceholder()
	update.Add(refreshTokenToken)
	buf.ReadFrom(update)
	buf.WriteString(" RETURNING ")
	if len(r.Columns) > 0 {
		buf.WriteString(strings.Join(r.Columns, ", "))
	} else {
		buf.WriteString("created_at, created_by, expire_at, last_used_at, notes, revoked, token, updated_at, updated_by, user_id")
	}
	return buf.String(), update.Args(), nil
}

func (r *RefreshTokenRepositoryBase) updateOneByToken(ctx context.Context, tx *sql.Tx, refreshTokenToken string, p *RefreshTokenPatch) (*RefreshTokenEntity, error) {
	query, args, err := r.UpdateOneByTokenQuery(refreshTokenToken, p)
	if err != nil {
		return nil, err
	}
	var ent RefreshTokenEntity
	props, err := ent.Props(r.Columns...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(props...)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "update one by unique", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "update one by unique tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *RefreshTokenRepositoryBase) UpdateOneByToken(ctx context.Context, refreshTokenToken string, p *RefreshTokenPatch) (*RefreshTokenEntity, error) {
	return r.updateOneByToken(ctx, nil, refreshTokenToken, p)
}

func (r *RefreshTokenRepositoryBase) UpsertQuery(e *RefreshTokenEntity, p *RefreshTokenPatch, inf ...string) (string, []interface{}, error) {
	upsert := NewComposer(20)
	columns := bytes.NewBuffer(nil)
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(r.Table)

	if !e.CreatedAt.IsZero() {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnCreatedAt); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnCreatedBy); err != nil {
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

	if e.ExpireAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnExpireAt); err != nil {
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
		upsert.Add(e.ExpireAt)
		upsert.Dirty = true
	}

	if e.LastUsedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnLastUsedAt); err != nil {
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
		upsert.Add(e.LastUsedAt)
		upsert.Dirty = true
	}

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnNotes); err != nil {
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
	upsert.Add(e.Notes)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnRevoked); err != nil {
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
	upsert.Add(e.Revoked)
	upsert.Dirty = true

	if columns.Len() > 0 {
		if _, err := columns.WriteString(", "); err != nil {
			return "", nil, err
		}
	}
	if _, err := columns.WriteString(TableRefreshTokenColumnToken); err != nil {
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
	upsert.Add(e.Token)
	upsert.Dirty = true

	if e.UpdatedAt.Valid {
		if columns.Len() > 0 {
			if _, err := columns.WriteString(", "); err != nil {
				return "", nil, err
			}
		}
		if _, err := columns.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnUpdatedBy); err != nil {
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
	if _, err := columns.WriteString(TableRefreshTokenColumnUserID); err != nil {
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
			if _, err := upsert.WriteString(TableRefreshTokenColumnCreatedAt); err != nil {
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
			if _, err := upsert.WriteString(TableRefreshTokenColumnCreatedBy); err != nil {
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

		if p.ExpireAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnExpireAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.ExpireAt)
			upsert.Dirty = true

		}
		if p.LastUsedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnLastUsedAt); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.LastUsedAt)
			upsert.Dirty = true

		}
		if p.Notes.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnNotes); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Notes)
			upsert.Dirty = true
		}

		if p.Revoked.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnRevoked); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Revoked)
			upsert.Dirty = true
		}

		if p.Token.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnToken); err != nil {
				return "", nil, err
			}
			if _, err := upsert.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := upsert.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			upsert.Add(p.Token)
			upsert.Dirty = true
		}

		if p.UpdatedAt.Valid {
			if upsert.Dirty {
				if _, err := upsert.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := upsert.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
			if _, err := upsert.WriteString(TableRefreshTokenColumnUpdatedAt); err != nil {
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
			if _, err := upsert.WriteString(TableRefreshTokenColumnUpdatedBy); err != nil {
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
			if _, err := upsert.WriteString(TableRefreshTokenColumnUserID); err != nil {
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
			buf.WriteString("created_at, created_by, expire_at, last_used_at, notes, revoked, token, updated_at, updated_by, user_id")
		}
	}
	return buf.String(), upsert.Args(), nil
}

func (r *RefreshTokenRepositoryBase) upsert(ctx context.Context, tx *sql.Tx, e *RefreshTokenEntity, p *RefreshTokenPatch, inf ...string) (*RefreshTokenEntity, error) {
	query, args, err := r.UpsertQuery(e, p, inf...)
	if err != nil {
		return nil, err
	}

	var row *sql.Row
	if tx == nil {
		row = r.DB.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&e.CreatedAt,
		&e.CreatedBy,
		&e.ExpireAt,
		&e.LastUsedAt,
		&e.Notes,
		&e.Revoked,
		&e.Token,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.UserID,
	)
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "upsert", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "upsert tx", query, args...)
		}
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *RefreshTokenRepositoryBase) Upsert(ctx context.Context, e *RefreshTokenEntity, p *RefreshTokenPatch, inf ...string) (*RefreshTokenEntity, error) {
	return r.upsert(ctx, nil, e, p, inf...)
}

func (r *RefreshTokenRepositoryBase) count(ctx context.Context, tx *sql.Tx, exp *RefreshTokenCountExpr) (int64, error) {
	query, args, err := r.FindQuery(&RefreshTokenFindExpr{
		Where:   exp.Where,
		Columns: []string{"COUNT(*)"},

		JoinUser:     exp.JoinUser,
		JoinAuthor:   exp.JoinAuthor,
		JoinModifier: exp.JoinModifier,
	})
	if err != nil {
		return 0, err
	}
	var count int64
	if tx == nil {
		err = r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if r.Log != nil {
		if tx == nil {
			r.Log(err, TableRefreshToken, "count", query, args...)
		} else {
			r.Log(err, TableRefreshToken, "count tx", query, args...)
		}
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RefreshTokenRepositoryBase) Count(ctx context.Context, exp *RefreshTokenCountExpr) (int64, error) {
	return r.count(ctx, nil, exp)
}

type RefreshTokenRepositoryBaseTx struct {
	base *RefreshTokenRepositoryBase
	tx   *sql.Tx
}

func (r RefreshTokenRepositoryBaseTx) Commit() error {
	return r.tx.Commit()
}

func (r RefreshTokenRepositoryBaseTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *RefreshTokenRepositoryBaseTx) Insert(ctx context.Context, e *RefreshTokenEntity) (*RefreshTokenEntity, error) {
	return r.base.insert(ctx, r.tx, e)
}

func (r *RefreshTokenRepositoryBaseTx) Find(ctx context.Context, fe *RefreshTokenFindExpr) ([]*RefreshTokenEntity, error) {
	return r.base.find(ctx, r.tx, fe)
}

func (r *RefreshTokenRepositoryBaseTx) FindIter(ctx context.Context, fe *RefreshTokenFindExpr) (*RefreshTokenIterator, error) {
	return r.base.findIter(ctx, r.tx, fe)
}

func (r *RefreshTokenRepositoryBaseTx) UpdateOneByToken(ctx context.Context, refreshTokenToken string, p *RefreshTokenPatch) (*RefreshTokenEntity, error) {
	return r.base.updateOneByToken(ctx, r.tx, refreshTokenToken, p)
}

func (r *RefreshTokenRepositoryBaseTx) Upsert(ctx context.Context, e *RefreshTokenEntity, p *RefreshTokenPatch, inf ...string) (*RefreshTokenEntity, error) {
	return r.base.upsert(ctx, r.tx, e, p, inf...)
}

func (r *RefreshTokenRepositoryBaseTx) Count(ctx context.Context, exp *RefreshTokenCountExpr) (int64, error) {
	return r.base.count(ctx, r.tx, exp)
}

const (
	JoinInner = iota
	JoinLeft
	JoinRight
	JoinCross
	JoinDoNot
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

// Actionable returns true if JoinType is one of the known type except JoinDoNot.
func (jt JoinType) Actionable() bool {
	switch jt {
	case JoinInner, JoinLeft, JoinRight, JoinCross:
		return true
	default:
		return false
	}
}

// ErrorConstraint returns the error constraint of err if it was produced by the pq library.
// Otherwise, it returns empty string.
func ErrorConstraint(err error) string {
	if err == nil {
		return ""
	}
	if pqerr, ok := err.(*pq.Error); ok {
		return pqerr.Constraint
	}

	return ""
}

type RowOrder struct {
	Name       string
	Descending bool
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

	var (
		tmp  []string
		srcs string
	)

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
-- sql schema beginning
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

	CONSTRAINT "charon.user_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "charon.user_username_key" UNIQUE (username),
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

CREATE TABLE IF NOT EXISTS charon.refresh_token (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	expire_at TIMESTAMPTZ,
	last_used_at TIMESTAMPTZ,
	notes TEXT,
	revoked BOOL DEFAULT false NOT NULL,
	token TEXT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	user_id BIGINT NOT NULL,

	CONSTRAINT "charon.refresh_token_token_key" UNIQUE (token),
	CONSTRAINT "charon.refresh_token_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.refresh_token_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.refresh_token_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id)
);

-- sql schema end
`
