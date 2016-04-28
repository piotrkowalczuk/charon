package charon

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/nilt"
	"github.com/piotrkowalczuk/pqcomp"
	"github.com/piotrkowalczuk/protot"
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
	CreatedBy         *nilt.Int64
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
	UpdatedBy         *nilt.Int64
	Username          string
	Author            *userEntity
	Modifier          *userEntity
	Permission        []*permissionEntity
	Group             []*groupEntity
}
type userCriteria struct {
	offset, limit     int64
	sort              map[string]bool
	confirmationToken []byte
	createdAt         *protot.QueryTimestamp
	createdBy         *protot.QueryInt64
	firstName         *protot.QueryString
	id                *protot.QueryInt64
	isActive          *nilt.Bool
	isConfirmed       *nilt.Bool
	isStaff           *nilt.Bool
	isSuperuser       *nilt.Bool
	lastLoginAt       *protot.QueryTimestamp
	lastName          *protot.QueryString
	password          []byte
	updatedAt         *protot.QueryTimestamp
	updatedBy         *protot.QueryInt64
	username          *protot.QueryString
}

type userRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userRepository) Find(c *userCriteria) ([]*userEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(15)
	where.AddExpr(tableUserColumnConfirmationToken, pqcomp.Equal, c.confirmationToken)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tableUserColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tableUserColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.createdBy.Negation {
				where.AddExpr(tableUserColumnCreatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnCreatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.Equal, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.NotEqual, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GreaterThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LessThan, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LessThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserColumnCreatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Values[0])
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LessThan, c.createdBy.Values[1])
		}
	}

	if c.firstName != nil && c.firstName.Valid {
		switch c.firstName.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.firstName.Negation {
				where.AddExpr(tableUserColumnFirstName, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnFirstName, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnFirstName, pqcomp.Equal, c.firstName.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnFirstName, pqcomp.Like, "%"+c.firstName.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnFirstName, pqcomp.Like, c.firstName.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnFirstName, pqcomp.Like, "%"+c.firstName.Value())
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.id.Negation {
				where.AddExpr(tableUserColumnID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.Equal, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.NotEqual, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnID, pqcomp.GreaterThan, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.GreaterThanOrEqual, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnID, pqcomp.LessThan, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.LessThanOrEqual, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tableUserColumnID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnID, pqcomp.GreaterThan, c.id.Values[0])
			where.AddExpr(tableUserColumnID, pqcomp.LessThan, c.id.Values[1])
		}
	}

	where.AddExpr(tableUserColumnIsActive, pqcomp.Equal, c.isActive)
	where.AddExpr(tableUserColumnIsConfirmed, pqcomp.Equal, c.isConfirmed)
	where.AddExpr(tableUserColumnIsStaff, pqcomp.Equal, c.isStaff)
	where.AddExpr(tableUserColumnIsSuperuser, pqcomp.Equal, c.isSuperuser)

	if c.lastLoginAt != nil && c.lastLoginAt.Valid {
		lastLoginAtt1 := c.lastLoginAt.Value()
		if lastLoginAtt1 != nil {
			lastLoginAt1, err := ptypes.Timestamp(lastLoginAtt1)
			if err != nil {
				return nil, err
			}

			switch c.lastLoginAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.lastLoginAt.Negation {
					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.Equal, lastLoginAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.NotEqual, lastLoginAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GreaterThan, lastLoginAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GreaterThanOrEqual, lastLoginAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LessThan, lastLoginAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LessThanOrEqual, lastLoginAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.In, lastLoginAt1)
			case protot.NumericQueryType_BETWEEN:
				lastLoginAtt2 := c.lastLoginAt.Values[1]
				if lastLoginAtt2 != nil {
					lastLoginAt2, err := ptypes.Timestamp(lastLoginAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GreaterThan, lastLoginAt1)
					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LessThan, lastLoginAt2)
				}
			}
		}
	}

	if c.lastName != nil && c.lastName.Valid {
		switch c.lastName.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.lastName.Negation {
				where.AddExpr(tableUserColumnLastName, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnLastName, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnLastName, pqcomp.Equal, c.lastName.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnLastName, pqcomp.Like, "%"+c.lastName.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnLastName, pqcomp.Like, c.lastName.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnLastName, pqcomp.Like, "%"+c.lastName.Value())
		}
	}

	where.AddExpr(tableUserColumnPassword, pqcomp.Equal, c.password)

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.updatedBy.Negation {
				where.AddExpr(tableUserColumnUpdatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnUpdatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.Equal, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.NotEqual, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GreaterThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LessThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserColumnUpdatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Values[0])
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Values[1])
		}
	}

	if c.username != nil && c.username.Valid {
		switch c.username.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.username.Negation {
				where.AddExpr(tableUserColumnUsername, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserColumnUsername, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnUsername, pqcomp.Equal, c.username.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnUsername, pqcomp.Like, "%"+c.username.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnUsername, pqcomp.Like, c.username.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnUsername, pqcomp.Like, "%"+c.username.Value())
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
	createdBy *nilt.Int64,
	firstName *nilt.String,
	isActive *nilt.Bool,
	isConfirmed *nilt.Bool,
	isStaff *nilt.Bool,
	isSuperuser *nilt.Bool,
	lastLoginAt *time.Time,
	lastName *nilt.String,
	password []byte,
	updatedAt *time.Time,
	updatedBy *nilt.Int64,
	username *nilt.String,
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
	CreatedBy   *nilt.Int64
	Description *nilt.String
	ID          int64
	Name        string
	UpdatedAt   *time.Time
	UpdatedBy   *nilt.Int64
	Author      *userEntity
	Modifier    *userEntity
	Permission  []*permissionEntity
	Users       []*userEntity
}
type groupCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *protot.QueryTimestamp
	createdBy     *protot.QueryInt64
	description   *protot.QueryString
	id            *protot.QueryInt64
	name          *protot.QueryString
	updatedAt     *protot.QueryTimestamp
	updatedBy     *protot.QueryInt64
}

type groupRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *groupRepository) Find(c *groupCriteria) ([]*groupEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(7)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.createdBy.Negation {
				where.AddExpr(tableGroupColumnCreatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupColumnCreatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.Equal, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.NotEqual, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GreaterThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LessThan, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LessThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableGroupColumnCreatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Values[0])
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LessThan, c.createdBy.Values[1])
		}
	}

	if c.description != nil && c.description.Valid {
		switch c.description.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.description.Negation {
				where.AddExpr(tableGroupColumnDescription, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupColumnDescription, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupColumnDescription, pqcomp.Equal, c.description.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupColumnDescription, pqcomp.Like, "%"+c.description.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupColumnDescription, pqcomp.Like, c.description.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupColumnDescription, pqcomp.Like, "%"+c.description.Value())
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.id.Negation {
				where.AddExpr(tableGroupColumnID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupColumnID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.Equal, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.NotEqual, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnID, pqcomp.GreaterThan, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.GreaterThanOrEqual, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnID, pqcomp.LessThan, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.LessThanOrEqual, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tableGroupColumnID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnID, pqcomp.GreaterThan, c.id.Values[0])
			where.AddExpr(tableGroupColumnID, pqcomp.LessThan, c.id.Values[1])
		}
	}

	if c.name != nil && c.name.Valid {
		switch c.name.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.name.Negation {
				where.AddExpr(tableGroupColumnName, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupColumnName, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupColumnName, pqcomp.Equal, c.name.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupColumnName, pqcomp.Like, "%"+c.name.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupColumnName, pqcomp.Like, c.name.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupColumnName, pqcomp.Like, "%"+c.name.Value())
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.updatedBy.Negation {
				where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.Equal, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.NotEqual, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GreaterThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LessThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Values[0])
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Values[1])
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
	createdBy *nilt.Int64,
	description *nilt.String,
	name *nilt.String,
	updatedAt *time.Time,
	updatedBy *nilt.Int64,
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
	action        *protot.QueryString
	createdAt     *protot.QueryTimestamp
	id            *protot.QueryInt64
	module        *protot.QueryString
	subsystem     *protot.QueryString
	updatedAt     *protot.QueryTimestamp
}

type permissionRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *permissionRepository) Find(c *permissionCriteria) ([]*permissionEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(6)

	if c.action != nil && c.action.Valid {
		switch c.action.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.action.Negation {
				where.AddExpr(tablePermissionColumnAction, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tablePermissionColumnAction, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnAction, pqcomp.Equal, c.action.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnAction, pqcomp.Like, "%"+c.action.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnAction, pqcomp.Like, c.action.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnAction, pqcomp.Like, "%"+c.action.Value())
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.id.Negation {
				where.AddExpr(tablePermissionColumnID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tablePermissionColumnID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.Equal, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.NotEqual, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tablePermissionColumnID, pqcomp.GreaterThan, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.GreaterThanOrEqual, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tablePermissionColumnID, pqcomp.LessThan, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.LessThanOrEqual, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tablePermissionColumnID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tablePermissionColumnID, pqcomp.GreaterThan, c.id.Values[0])
			where.AddExpr(tablePermissionColumnID, pqcomp.LessThan, c.id.Values[1])
		}
	}

	if c.module != nil && c.module.Valid {
		switch c.module.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.module.Negation {
				where.AddExpr(tablePermissionColumnModule, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tablePermissionColumnModule, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnModule, pqcomp.Equal, c.module.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnModule, pqcomp.Like, "%"+c.module.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnModule, pqcomp.Like, c.module.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnModule, pqcomp.Like, "%"+c.module.Value())
		}
	}

	if c.subsystem != nil && c.subsystem.Valid {
		switch c.subsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.subsystem.Negation {
				where.AddExpr(tablePermissionColumnSubsystem, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tablePermissionColumnSubsystem, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.Equal, c.subsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.Like, "%"+c.subsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.Like, c.subsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.Like, "%"+c.subsystem.Value())
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
	action *nilt.String,
	createdAt *time.Time,
	module *nilt.String,
	subsystem *nilt.String,
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
	CreatedBy *nilt.Int64
	GroupID   int64
	UpdatedAt *time.Time
	UpdatedBy *nilt.Int64
	UserID    int64
	User      *userEntity
	Group     *groupEntity
	Author    *userEntity
	Modifier  *userEntity
}
type userGroupsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     *protot.QueryTimestamp
	createdBy     *protot.QueryInt64
	groupID       *protot.QueryInt64
	updatedAt     *protot.QueryTimestamp
	updatedBy     *protot.QueryInt64
	userID        *protot.QueryInt64
}

type userGroupsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userGroupsRepository) Find(c *userGroupsCriteria) ([]*userGroupsEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(6)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.createdBy.Negation {
				where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.Equal, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.NotEqual, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GreaterThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LessThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Values[0])
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.groupID.Negation {
				where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.Equal, c.groupID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.NotEqual, c.groupID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GreaterThan, c.groupID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GreaterThanOrEqual, c.groupID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LessThan, c.groupID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LessThanOrEqual, c.groupID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.groupID.Values {
				where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GreaterThan, c.groupID.Values[0])
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LessThan, c.groupID.Values[1])
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.updatedBy.Negation {
				where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.Equal, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.NotEqual, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GreaterThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LessThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Values[0])
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.userID.Negation {
				where.AddExpr(tableUserGroupsColumnUserID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserGroupsColumnUserID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.Equal, c.userID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.NotEqual, c.userID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GreaterThan, c.userID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GreaterThanOrEqual, c.userID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LessThan, c.userID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LessThanOrEqual, c.userID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.userID.Values {
				where.AddExpr(tableUserGroupsColumnUserID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GreaterThan, c.userID.Values[0])
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LessThan, c.userID.Values[1])
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
	CreatedBy           *nilt.Int64
	GroupID             int64
	PermissionAction    string
	PermissionModule    string
	PermissionSubsystem string
	UpdatedAt           *time.Time
	UpdatedBy           *nilt.Int64
	Group               *groupEntity
	Permission          *permissionEntity
	Author              *userEntity
	Modifier            *userEntity
}
type groupPermissionsCriteria struct {
	offset, limit       int64
	sort                map[string]bool
	createdAt           *protot.QueryTimestamp
	createdBy           *protot.QueryInt64
	groupID             *protot.QueryInt64
	permissionAction    *protot.QueryString
	permissionModule    *protot.QueryString
	permissionSubsystem *protot.QueryString
	updatedAt           *protot.QueryTimestamp
	updatedBy           *protot.QueryInt64
}

type groupPermissionsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *groupPermissionsRepository) Find(c *groupPermissionsCriteria) ([]*groupPermissionsEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(8)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.createdBy.Negation {
				where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.Equal, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.NotEqual, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GreaterThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LessThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Values[0])
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.groupID.Negation {
				where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.Equal, c.groupID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.NotEqual, c.groupID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GreaterThan, c.groupID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GreaterThanOrEqual, c.groupID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LessThan, c.groupID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LessThanOrEqual, c.groupID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.groupID.Values {
				where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GreaterThan, c.groupID.Values[0])
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LessThan, c.groupID.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionAction.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.Equal, c.permissionAction.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.Like, "%"+c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.Like, c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.Like, "%"+c.permissionAction.Value())
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionModule.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.Equal, c.permissionModule.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.Like, "%"+c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.Like, c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.Like, "%"+c.permissionModule.Value())
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionSubsystem.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.Equal, c.permissionSubsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.Like, "%"+c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.Like, c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.Like, "%"+c.permissionSubsystem.Value())
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.updatedBy.Negation {
				where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.Equal, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.NotEqual, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GreaterThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LessThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Values[0])
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Values[1])
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
	CreatedBy           *nilt.Int64
	PermissionAction    string
	PermissionModule    string
	PermissionSubsystem string
	UpdatedAt           *time.Time
	UpdatedBy           *nilt.Int64
	UserID              int64
	User                *userEntity
	Permission          *permissionEntity
	Author              *userEntity
	Modifier            *userEntity
}
type userPermissionsCriteria struct {
	offset, limit       int64
	sort                map[string]bool
	createdAt           *protot.QueryTimestamp
	createdBy           *protot.QueryInt64
	permissionAction    *protot.QueryString
	permissionModule    *protot.QueryString
	permissionSubsystem *protot.QueryString
	updatedAt           *protot.QueryTimestamp
	updatedBy           *protot.QueryInt64
	userID              *protot.QueryInt64
}

type userPermissionsRepository struct {
	table   string
	columns []string
	db      *sql.DB
}

func (r *userPermissionsRepository) Find(c *userPermissionsCriteria) ([]*userPermissionsEntity, error) {
	comp := pqcomp.New(2, 0, 1)
	comp.AddArg(c.offset)
	comp.AddArg(c.limit)

	where := comp.Compose(8)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.createdAt.Negation {
					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.Equal, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.NotEqual, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GreaterThanOrEqual, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LessThan, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LessThanOrEqual, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.In, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GreaterThan, createdAt1)
					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LessThan, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.createdBy.Negation {
				where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.Equal, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.NotEqual, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GreaterThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LessThanOrEqual, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GreaterThan, c.createdBy.Values[0])
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LessThan, c.createdBy.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionAction.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.Equal, c.permissionAction.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.Like, "%"+c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.Like, c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.Like, "%"+c.permissionAction.Value())
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionModule.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.Equal, c.permissionModule.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.Like, "%"+c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.Like, c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.Like, "%"+c.permissionModule.Value())
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionSubsystem.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.IsNull, "")
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.Equal, c.permissionSubsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.Like, "%"+c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.Like, c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.Like, "%"+c.permissionSubsystem.Value())
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
			case protot.NumericQueryType_NOT_A_NUMBER:
				if c.updatedAt.Negation {
					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.IsNotNull, "")
				} else {
					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.IsNull, "")
				}
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.Equal, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.NotEqual, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GreaterThanOrEqual, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LessThan, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LessThanOrEqual, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.In, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GreaterThan, updatedAt1)
					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LessThan, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.updatedBy.Negation {
				where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.Equal, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.NotEqual, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GreaterThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LessThanOrEqual, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GreaterThan, c.updatedBy.Values[0])
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LessThan, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			if c.userID.Negation {
				where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.IsNotNull, "")
			} else {
				where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.IsNull, "")
			}
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.Equal, c.userID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.NotEqual, c.userID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GreaterThan, c.userID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GreaterThanOrEqual, c.userID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LessThan, c.userID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LessThanOrEqual, c.userID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.userID.Values {
				where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.In, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GreaterThan, c.userID.Values[0])
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LessThan, c.userID.Values[1])
		}
	}

	rows, err := findQueryComp(r.db, r.table, comp, where, c.sort, r.columns)
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
