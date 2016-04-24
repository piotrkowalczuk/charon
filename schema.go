package charon

import (
	"database/sql"
	"errors"
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
	where.AddExpr(tableUserColumnConfirmationToken, pqcomp.E, c.confirmationToken)

	if c.createdAt != nil && c.createdAt.Valid {
		createdAtt1 := c.createdAt.Value()
		if createdAtt1 != nil {
			createdAt1, err := ptypes.Timestamp(createdAtt1)
			if err != nil {
				return nil, err
			}

			switch c.createdAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tableUserColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.E, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.NE, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GT, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GTE, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LT, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LTE, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserColumnCreatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.GT, c.createdBy.Values[0])
			where.AddExpr(tableUserColumnCreatedBy, pqcomp.LT, c.createdBy.Values[1])
		}
	}

	if c.firstName != nil && c.firstName.Valid {
		switch c.firstName.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.firstName.Negation {
				where.AddExpr(tableUserColumnFirstName, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserColumnFirstName, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnFirstName, pqcomp.E, c.firstName.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnFirstName, pqcomp.LIKE, "%"+c.firstName.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnFirstName, pqcomp.LIKE, c.firstName.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnFirstName, pqcomp.LIKE, "%"+c.firstName.Value())
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserColumnID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.E, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.NE, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnID, pqcomp.GT, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.GTE, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnID, pqcomp.LT, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnID, pqcomp.LTE, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tableUserColumnID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnID, pqcomp.GT, c.id.Values[0])
			where.AddExpr(tableUserColumnID, pqcomp.LT, c.id.Values[1])
		}
	}

	where.AddExpr(tableUserColumnIsActive, pqcomp.E, c.isActive)
	where.AddExpr(tableUserColumnIsConfirmed, pqcomp.E, c.isConfirmed)
	where.AddExpr(tableUserColumnIsStaff, pqcomp.E, c.isStaff)
	where.AddExpr(tableUserColumnIsSuperuser, pqcomp.E, c.isSuperuser)

	if c.lastLoginAt != nil && c.lastLoginAt.Valid {
		lastLoginAtt1 := c.lastLoginAt.Value()
		if lastLoginAtt1 != nil {
			lastLoginAt1, err := ptypes.Timestamp(lastLoginAtt1)
			if err != nil {
				return nil, err
			}

			switch c.lastLoginAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.E, lastLoginAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.NE, lastLoginAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GT, lastLoginAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GTE, lastLoginAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LT, lastLoginAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LTE, lastLoginAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnLastLoginAt, pqcomp.IN, lastLoginAt1)
			case protot.NumericQueryType_BETWEEN:
				lastLoginAtt2 := c.lastLoginAt.Values[1]
				if lastLoginAtt2 != nil {
					lastLoginAt2, err := ptypes.Timestamp(lastLoginAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GT, lastLoginAt1)
					where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LT, lastLoginAt2)
				}
			}
		}
	}

	if c.lastName != nil && c.lastName.Valid {
		switch c.lastName.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.lastName.Negation {
				where.AddExpr(tableUserColumnLastName, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserColumnLastName, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnLastName, pqcomp.E, c.lastName.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnLastName, pqcomp.LIKE, "%"+c.lastName.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnLastName, pqcomp.LIKE, c.lastName.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnLastName, pqcomp.LIKE, "%"+c.lastName.Value())
		}
	}

	where.AddExpr(tableUserColumnPassword, pqcomp.E, c.password)

	if c.updatedAt != nil && c.updatedAt.Valid {
		updatedAtt1 := c.updatedAt.Value()
		if updatedAtt1 != nil {
			updatedAt1, err := ptypes.Timestamp(updatedAtt1)
			if err != nil {
				return nil, err
			}

			switch c.updatedAt.Type {
			case protot.NumericQueryType_NOT_A_NUMBER:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LT, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.E, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.NE, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GT, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GTE, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LT, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LTE, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserColumnUpdatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.GT, c.updatedBy.Values[0])
			where.AddExpr(tableUserColumnUpdatedBy, pqcomp.LT, c.updatedBy.Values[1])
		}
	}

	if c.username != nil && c.username.Valid {
		switch c.username.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.username.Negation {
				where.AddExpr(tableUserColumnUsername, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserColumnUsername, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserColumnUsername, pqcomp.E, c.username.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserColumnUsername, pqcomp.LIKE, "%"+c.username.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserColumnUsername, pqcomp.LIKE, c.username.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserColumnUsername, pqcomp.LIKE, "%"+c.username.Value())
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.ConfirmationToken,
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
	update.AddExpr(tableUserColumnID, pqcomp.E, id)
	update.AddExpr(tableUserColumnConfirmationToken, pqcomp.E, confirmationToken)
	if createdAt != nil {
		update.AddExpr(tableUserColumnCreatedAt, pqcomp.E, createdAt)

	}
	update.AddExpr(tableUserColumnCreatedBy, pqcomp.E, createdBy)
	update.AddExpr(tableUserColumnFirstName, pqcomp.E, firstName)

	update.AddExpr(tableUserColumnIsActive, pqcomp.E, isActive)

	update.AddExpr(tableUserColumnIsConfirmed, pqcomp.E, isConfirmed)

	update.AddExpr(tableUserColumnIsStaff, pqcomp.E, isStaff)

	update.AddExpr(tableUserColumnIsSuperuser, pqcomp.E, isSuperuser)
	update.AddExpr(tableUserColumnLastLoginAt, pqcomp.E, lastLoginAt)
	update.AddExpr(tableUserColumnLastName, pqcomp.E, lastName)
	update.AddExpr(tableUserColumnPassword, pqcomp.E, password)
	if updatedAt != nil {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.E, updatedAt)
	} else {
		update.AddExpr(tableUserColumnUpdatedAt, pqcomp.E, "NOW()")
	}
	update.AddExpr(tableUserColumnUpdatedBy, pqcomp.E, updatedBy)
	update.AddExpr(tableUserColumnUsername, pqcomp.E, username)

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
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.E, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.NE, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GT, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GTE, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LT, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LTE, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableGroupColumnCreatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.GT, c.createdBy.Values[0])
			where.AddExpr(tableGroupColumnCreatedBy, pqcomp.LT, c.createdBy.Values[1])
		}
	}

	if c.description != nil && c.description.Valid {
		switch c.description.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.description.Negation {
				where.AddExpr(tableGroupColumnDescription, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableGroupColumnDescription, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupColumnDescription, pqcomp.E, c.description.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupColumnDescription, pqcomp.LIKE, "%"+c.description.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupColumnDescription, pqcomp.LIKE, c.description.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupColumnDescription, pqcomp.LIKE, "%"+c.description.Value())
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupColumnID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.E, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.NE, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnID, pqcomp.GT, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.GTE, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnID, pqcomp.LT, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnID, pqcomp.LTE, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tableGroupColumnID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnID, pqcomp.GT, c.id.Values[0])
			where.AddExpr(tableGroupColumnID, pqcomp.LT, c.id.Values[1])
		}
	}

	if c.name != nil && c.name.Valid {
		switch c.name.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.name.Negation {
				where.AddExpr(tableGroupColumnName, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableGroupColumnName, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupColumnName, pqcomp.E, c.name.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupColumnName, pqcomp.LIKE, "%"+c.name.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupColumnName, pqcomp.LIKE, c.name.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupColumnName, pqcomp.LIKE, "%"+c.name.Value())
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
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LT, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.E, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.NE, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GT, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GTE, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LT, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LTE, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.GT, c.updatedBy.Values[0])
			where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.LT, c.updatedBy.Values[1])
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
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
	update.AddExpr(tableGroupColumnID, pqcomp.E, id)
	if createdAt != nil {
		update.AddExpr(tableGroupColumnCreatedAt, pqcomp.E, createdAt)

	}
	update.AddExpr(tableGroupColumnCreatedBy, pqcomp.E, createdBy)
	update.AddExpr(tableGroupColumnDescription, pqcomp.E, description)
	update.AddExpr(tableGroupColumnName, pqcomp.E, name)
	if updatedAt != nil {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.E, updatedAt)
	} else {
		update.AddExpr(tableGroupColumnUpdatedAt, pqcomp.E, "NOW()")
	}
	update.AddExpr(tableGroupColumnUpdatedBy, pqcomp.E, updatedBy)

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
				where.AddExpr(tablePermissionColumnAction, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tablePermissionColumnAction, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnAction, pqcomp.E, c.action.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnAction, pqcomp.LIKE, "%"+c.action.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnAction, pqcomp.LIKE, c.action.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnAction, pqcomp.LIKE, "%"+c.action.Value())
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
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.id != nil && c.id.Valid {
		switch c.id.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tablePermissionColumnID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.E, c.id.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.NE, c.id.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tablePermissionColumnID, pqcomp.GT, c.id.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.GTE, c.id.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tablePermissionColumnID, pqcomp.LT, c.id.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tablePermissionColumnID, pqcomp.LTE, c.id.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.id.Values {
				where.AddExpr(tablePermissionColumnID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tablePermissionColumnID, pqcomp.GT, c.id.Values[0])
			where.AddExpr(tablePermissionColumnID, pqcomp.LT, c.id.Values[1])
		}
	}

	if c.module != nil && c.module.Valid {
		switch c.module.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.module.Negation {
				where.AddExpr(tablePermissionColumnModule, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tablePermissionColumnModule, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnModule, pqcomp.E, c.module.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnModule, pqcomp.LIKE, "%"+c.module.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnModule, pqcomp.LIKE, c.module.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnModule, pqcomp.LIKE, "%"+c.module.Value())
		}
	}

	if c.subsystem != nil && c.subsystem.Valid {
		switch c.subsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.subsystem.Negation {
				where.AddExpr(tablePermissionColumnSubsystem, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tablePermissionColumnSubsystem, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.E, c.subsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.LIKE, "%"+c.subsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.LIKE, c.subsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tablePermissionColumnSubsystem, pqcomp.LIKE, "%"+c.subsystem.Value())
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
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LT, updatedAt2)
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.Action,
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
	update.AddExpr(tablePermissionColumnID, pqcomp.E, id)
	update.AddExpr(tablePermissionColumnAction, pqcomp.E, action)
	if createdAt != nil {
		update.AddExpr(tablePermissionColumnCreatedAt, pqcomp.E, createdAt)

	}
	update.AddExpr(tablePermissionColumnModule, pqcomp.E, module)
	update.AddExpr(tablePermissionColumnSubsystem, pqcomp.E, subsystem)
	if updatedAt != nil {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.E, updatedAt)
	} else {
		update.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.E, "NOW()")
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
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.E, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.NE, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GT, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GTE, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LT, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LTE, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.GT, c.createdBy.Values[0])
			where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.LT, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.E, c.groupID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.NE, c.groupID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GT, c.groupID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GTE, c.groupID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LT, c.groupID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LTE, c.groupID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.groupID.Values {
				where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.GT, c.groupID.Values[0])
			where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.LT, c.groupID.Values[1])
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
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LT, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.E, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.NE, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GTE, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LTE, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Values[0])
			where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.E, c.userID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.NE, c.userID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GT, c.userID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GTE, c.userID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LT, c.userID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LTE, c.userID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.userID.Values {
				where.AddExpr(tableUserGroupsColumnUserID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.GT, c.userID.Values[0])
			where.AddExpr(tableUserGroupsColumnUserID, pqcomp.LT, c.userID.Values[1])
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
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
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.E, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.NE, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GT, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GTE, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LT, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LTE, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.GT, c.createdBy.Values[0])
			where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.LT, c.createdBy.Values[1])
		}
	}

	if c.groupID != nil && c.groupID.Valid {
		switch c.groupID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.E, c.groupID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.NE, c.groupID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GT, c.groupID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GTE, c.groupID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LT, c.groupID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LTE, c.groupID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.groupID.Values {
				where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.GT, c.groupID.Values[0])
			where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.LT, c.groupID.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionAction.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.E, c.permissionAction.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.LIKE, "%"+c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.LIKE, c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionAction, pqcomp.LIKE, "%"+c.permissionAction.Value())
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionModule.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.E, c.permissionModule.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.LIKE, "%"+c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.LIKE, c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionModule, pqcomp.LIKE, "%"+c.permissionModule.Value())
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionSubsystem.Negation {
				where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.E, c.permissionSubsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.LIKE, "%"+c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.LIKE, c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableGroupPermissionsColumnPermissionSubsystem, pqcomp.LIKE, "%"+c.permissionSubsystem.Value())
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
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LT, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.E, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.NE, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GTE, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LTE, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Values[0])
			where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Values[1])
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
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
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.E, createdAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.NE, createdAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GT, createdAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GTE, createdAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LT, createdAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LTE, createdAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.IN, createdAt1)
			case protot.NumericQueryType_BETWEEN:
				createdAtt2 := c.createdAt.Values[1]
				if createdAtt2 != nil {
					createdAt2, err := ptypes.Timestamp(createdAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GT, createdAt1)
					where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LT, createdAt2)
				}
			}
		}
	}

	if c.createdBy != nil && c.createdBy.Valid {
		switch c.createdBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.E, c.createdBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.NE, c.createdBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GT, c.createdBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GTE, c.createdBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LT, c.createdBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LTE, c.createdBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.createdBy.Values {
				where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.GT, c.createdBy.Values[0])
			where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.LT, c.createdBy.Values[1])
		}
	}

	if c.permissionAction != nil && c.permissionAction.Valid {
		switch c.permissionAction.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionAction.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.E, c.permissionAction.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.LIKE, "%"+c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.LIKE, c.permissionAction.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionAction, pqcomp.LIKE, "%"+c.permissionAction.Value())
		}
	}

	if c.permissionModule != nil && c.permissionModule.Valid {
		switch c.permissionModule.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionModule.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.E, c.permissionModule.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.LIKE, "%"+c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.LIKE, c.permissionModule.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionModule, pqcomp.LIKE, "%"+c.permissionModule.Value())
		}
	}

	if c.permissionSubsystem != nil && c.permissionSubsystem.Valid {
		switch c.permissionSubsystem.Type {
		case protot.TextQueryType_NOT_A_TEXT:
			if c.permissionSubsystem.Negation {
				where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.IS, pqcomp.NOT_NULL)
			} else {
				where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.IS, pqcomp.NULL)
			}
		case protot.TextQueryType_EXACT:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.E, c.permissionSubsystem.Value())
		case protot.TextQueryType_SUBSTRING:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.LIKE, "%"+c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_PREFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.LIKE, c.permissionSubsystem.Value()+"%")
		case protot.TextQueryType_HAS_SUFFIX:
			where.AddExpr(tableUserPermissionsColumnPermissionSubsystem, pqcomp.LIKE, "%"+c.permissionSubsystem.Value())
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
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.IS, pqcomp.NULL)
			case protot.NumericQueryType_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.E, updatedAt1)
			case protot.NumericQueryType_NOT_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.NE, updatedAt1)
			case protot.NumericQueryType_GREATER:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GT, updatedAt1)
			case protot.NumericQueryType_GREATER_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GTE, updatedAt1)
			case protot.NumericQueryType_LESS:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LT, updatedAt1)
			case protot.NumericQueryType_LESS_EQUAL:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LTE, updatedAt1)
			case protot.NumericQueryType_IN:
				where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.IN, updatedAt1)
			case protot.NumericQueryType_BETWEEN:
				updatedAtt2 := c.updatedAt.Values[1]
				if updatedAtt2 != nil {
					updatedAt2, err := ptypes.Timestamp(updatedAtt2)
					if err != nil {
						return nil, err
					}

					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GT, updatedAt1)
					where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LT, updatedAt2)
				}
			}
		}
	}

	if c.updatedBy != nil && c.updatedBy.Valid {
		switch c.updatedBy.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.E, c.updatedBy.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.NE, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GTE, c.updatedBy.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LTE, c.updatedBy.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.updatedBy.Values {
				where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.GT, c.updatedBy.Values[0])
			where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.LT, c.updatedBy.Values[1])
		}
	}

	if c.userID != nil && c.userID.Valid {
		switch c.userID.Type {
		case protot.NumericQueryType_NOT_A_NUMBER:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.IS, pqcomp.NULL)
		case protot.NumericQueryType_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.E, c.userID.Value())
		case protot.NumericQueryType_NOT_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.NE, c.userID.Value())
		case protot.NumericQueryType_GREATER:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GT, c.userID.Value())
		case protot.NumericQueryType_GREATER_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GTE, c.userID.Value())
		case protot.NumericQueryType_LESS:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LT, c.userID.Value())
		case protot.NumericQueryType_LESS_EQUAL:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LTE, c.userID.Value())
		case protot.NumericQueryType_IN:
			for _, v := range c.userID.Values {
				where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.IN, v)
			}
		case protot.NumericQueryType_BETWEEN:
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.GT, c.userID.Values[0])
			where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.LT, c.userID.Values[1])
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
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
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
