package main

import (
	"database/sql"
	"errors"
	"strings"
	"time"

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
	CreatedBy         nilt.Int64
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
	UpdatedBy         nilt.Int64
	Username          string
	Author            *userEntity
	Modifier          *userEntity
	Group             []*groupEntity
	Permission        []*permissionEntity
}
type userCriteria struct {
	offset, limit     int64
	sort              map[string]bool
	confirmationToken []byte
	createdAt         protot.TimestampRange
	createdBy         nilt.Int64
	firstName         nilt.String
	id                nilt.Int64
	isActive          nilt.Bool
	isConfirmed       nilt.Bool
	isStaff           nilt.Bool
	isSuperuser       nilt.Bool
	lastLoginAt       protot.TimestampRange
	lastName          nilt.String
	password          []byte
	updatedAt         protot.TimestampRange
	updatedBy         nilt.Int64
	username          nilt.String
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
	if c.createdAt.From != nil {
		where.AddExpr(tableUserColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tableUserColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tableUserColumnCreatedBy, pqcomp.E, c.createdBy)
	where.AddExpr(tableUserColumnFirstName, pqcomp.E, c.firstName)
	where.AddExpr(tableUserColumnID, pqcomp.E, c.id)
	where.AddExpr(tableUserColumnIsActive, pqcomp.E, c.isActive)
	where.AddExpr(tableUserColumnIsConfirmed, pqcomp.E, c.isConfirmed)
	where.AddExpr(tableUserColumnIsStaff, pqcomp.E, c.isStaff)
	where.AddExpr(tableUserColumnIsSuperuser, pqcomp.E, c.isSuperuser)
	if c.lastLoginAt.From != nil {
		where.AddExpr(tableUserColumnLastLoginAt, pqcomp.GT, c.lastLoginAt.From.Time())
	}
	if c.lastLoginAt.To != nil {
		where.AddExpr(tableUserColumnLastLoginAt, pqcomp.LT, c.lastLoginAt.To.Time())
	}

	where.AddExpr(tableUserColumnLastName, pqcomp.E, c.lastName)
	where.AddExpr(tableUserColumnPassword, pqcomp.E, c.password)
	if c.updatedAt.From != nil {
		where.AddExpr(tableUserColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tableUserColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
	}

	where.AddExpr(tableUserColumnUpdatedBy, pqcomp.E, c.updatedBy)
	where.AddExpr(tableUserColumnUsername, pqcomp.E, c.username)

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
	insert.AddExpr(tableUserColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableUserColumnFirstName, "", e.FirstName)
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
	createdBy nilt.Int64,
	firstName nilt.String,
	isActive nilt.Bool,
	isConfirmed nilt.Bool,
	isStaff nilt.Bool,
	isSuperuser nilt.Bool,
	lastLoginAt *time.Time,
	lastName nilt.String,
	password []byte,
	updatedAt *time.Time,
	updatedBy nilt.Int64,
	username nilt.String,
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
		return nil, errors.New("main: user update failure, nothing to update")
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
	CreatedBy   nilt.Int64
	Description nilt.String
	ID          int64
	Name        string
	UpdatedAt   *time.Time
	UpdatedBy   nilt.Int64
	Author      *userEntity
	Modifier    *userEntity
	Users       []*userEntity
	Permission  []*permissionEntity
}
type groupCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     protot.TimestampRange
	createdBy     nilt.Int64
	description   nilt.String
	id            nilt.Int64
	name          nilt.String
	updatedAt     protot.TimestampRange
	updatedBy     nilt.Int64
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
	if c.createdAt.From != nil {
		where.AddExpr(tableGroupColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tableGroupColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tableGroupColumnCreatedBy, pqcomp.E, c.createdBy)
	where.AddExpr(tableGroupColumnDescription, pqcomp.E, c.description)
	where.AddExpr(tableGroupColumnID, pqcomp.E, c.id)
	where.AddExpr(tableGroupColumnName, pqcomp.E, c.name)
	if c.updatedAt.From != nil {
		where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tableGroupColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
	}

	where.AddExpr(tableGroupColumnUpdatedBy, pqcomp.E, c.updatedBy)

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
	createdBy nilt.Int64,
	description nilt.String,
	name nilt.String,
	updatedAt *time.Time,
	updatedBy nilt.Int64,
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
		return nil, errors.New("main: group update failure, nothing to update")
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
	action        nilt.String
	createdAt     protot.TimestampRange
	id            nilt.Int64
	module        nilt.String
	subsystem     nilt.String
	updatedAt     protot.TimestampRange
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
	where.AddExpr(tablePermissionColumnAction, pqcomp.E, c.action)
	if c.createdAt.From != nil {
		where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tablePermissionColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tablePermissionColumnID, pqcomp.E, c.id)
	where.AddExpr(tablePermissionColumnModule, pqcomp.E, c.module)
	where.AddExpr(tablePermissionColumnSubsystem, pqcomp.E, c.subsystem)
	if c.updatedAt.From != nil {
		where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tablePermissionColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
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
	action nilt.String,
	createdAt *time.Time,
	module nilt.String,
	subsystem nilt.String,
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
		return nil, errors.New("main: permission update failure, nothing to update")
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
	CreatedBy nilt.Int64
	GroupID   int64
	UpdatedAt *time.Time
	UpdatedBy nilt.Int64
	UserID    int64
	User      *userEntity
	Group     *groupEntity
	Author    *userEntity
	Modifier  *userEntity
}
type userGroupsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     protot.TimestampRange
	createdBy     nilt.Int64
	groupID       nilt.Int64
	updatedAt     protot.TimestampRange
	updatedBy     nilt.Int64
	userID        nilt.Int64
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
	if c.createdAt.From != nil {
		where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tableUserGroupsColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tableUserGroupsColumnCreatedBy, pqcomp.E, c.createdBy)
	where.AddExpr(tableUserGroupsColumnGroupID, pqcomp.E, c.groupID)
	if c.updatedAt.From != nil {
		where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tableUserGroupsColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
	}

	where.AddExpr(tableUserGroupsColumnUpdatedBy, pqcomp.E, c.updatedBy)
	where.AddExpr(tableUserGroupsColumnUserID, pqcomp.E, c.userID)

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
	tableGroupPermissions                                 = "charon.group_permissions"
	tableGroupPermissionsColumnCreatedAt                  = "created_at"
	tableGroupPermissionsColumnCreatedBy                  = "created_by"
	tableGroupPermissionsColumnGroupID                    = "group_id"
	tableGroupPermissionsColumnPermissionID               = "permission_id"
	tableGroupPermissionsColumnUpdatedAt                  = "updated_at"
	tableGroupPermissionsColumnUpdatedBy                  = "updated_by"
	tableGroupPermissionsConstraintCreatedByForeignKey    = "charon.group_permissions_created_by_fkey"
	tableGroupPermissionsConstraintUpdatedByForeignKey    = "charon.group_permissions_updated_by_fkey"
	tableGroupPermissionsConstraintGroupIDForeignKey      = "charon.group_permissions_group_id_fkey"
	tableGroupPermissionsConstraintPermissionIDForeignKey = "charon.group_permissions_permission_id_fkey"
)

var (
	tableGroupPermissionsColumns = []string{
		tableGroupPermissionsColumnCreatedAt,
		tableGroupPermissionsColumnCreatedBy,
		tableGroupPermissionsColumnGroupID,
		tableGroupPermissionsColumnPermissionID,
		tableGroupPermissionsColumnUpdatedAt,
		tableGroupPermissionsColumnUpdatedBy,
	}
)

type groupPermissionsEntity struct {
	CreatedAt    time.Time
	CreatedBy    nilt.Int64
	GroupID      int64
	PermissionID int64
	UpdatedAt    *time.Time
	UpdatedBy    nilt.Int64
	Group        *groupEntity
	Permission   *permissionEntity
	Author       *userEntity
	Modifier     *userEntity
}
type groupPermissionsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     protot.TimestampRange
	createdBy     nilt.Int64
	groupID       nilt.Int64
	permissionID  nilt.Int64
	updatedAt     protot.TimestampRange
	updatedBy     nilt.Int64
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

	where := comp.Compose(6)
	if c.createdAt.From != nil {
		where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tableGroupPermissionsColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tableGroupPermissionsColumnCreatedBy, pqcomp.E, c.createdBy)
	where.AddExpr(tableGroupPermissionsColumnGroupID, pqcomp.E, c.groupID)
	where.AddExpr(tableGroupPermissionsColumnPermissionID, pqcomp.E, c.permissionID)
	if c.updatedAt.From != nil {
		where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tableGroupPermissionsColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
	}

	where.AddExpr(tableGroupPermissionsColumnUpdatedBy, pqcomp.E, c.updatedBy)

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
			&entity.PermissionID,
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
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tableGroupPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableGroupPermissionsColumnGroupID, "", e.GroupID)
	insert.AddExpr(tableGroupPermissionsColumnPermissionID, "", e.PermissionID)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableGroupPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.GroupID,
		&e.PermissionID,
		&e.UpdatedAt,
		&e.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

const (
	tableUserPermissions                                 = "charon.user_permissions"
	tableUserPermissionsColumnCreatedAt                  = "created_at"
	tableUserPermissionsColumnCreatedBy                  = "created_by"
	tableUserPermissionsColumnPermissionID               = "permission_id"
	tableUserPermissionsColumnUpdatedAt                  = "updated_at"
	tableUserPermissionsColumnUpdatedBy                  = "updated_by"
	tableUserPermissionsColumnUserID                     = "user_id"
	tableUserPermissionsConstraintCreatedByForeignKey    = "charon.user_permissions_created_by_fkey"
	tableUserPermissionsConstraintUpdatedByForeignKey    = "charon.user_permissions_updated_by_fkey"
	tableUserPermissionsConstraintUserIDForeignKey       = "charon.user_permissions_user_id_fkey"
	tableUserPermissionsConstraintPermissionIDForeignKey = "charon.user_permissions_permission_id_fkey"
)

var (
	tableUserPermissionsColumns = []string{
		tableUserPermissionsColumnCreatedAt,
		tableUserPermissionsColumnCreatedBy,
		tableUserPermissionsColumnPermissionID,
		tableUserPermissionsColumnUpdatedAt,
		tableUserPermissionsColumnUpdatedBy,
		tableUserPermissionsColumnUserID,
	}
)

type userPermissionsEntity struct {
	CreatedAt    time.Time
	CreatedBy    nilt.Int64
	PermissionID int64
	UpdatedAt    *time.Time
	UpdatedBy    nilt.Int64
	UserID       int64
	User         *userEntity
	Permission   *permissionEntity
	Author       *userEntity
	Modifier     *userEntity
}
type userPermissionsCriteria struct {
	offset, limit int64
	sort          map[string]bool
	createdAt     protot.TimestampRange
	createdBy     nilt.Int64
	permissionID  nilt.Int64
	updatedAt     protot.TimestampRange
	updatedBy     nilt.Int64
	userID        nilt.Int64
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

	where := comp.Compose(6)
	if c.createdAt.From != nil {
		where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.GT, c.createdAt.From.Time())
	}
	if c.createdAt.To != nil {
		where.AddExpr(tableUserPermissionsColumnCreatedAt, pqcomp.LT, c.createdAt.To.Time())
	}

	where.AddExpr(tableUserPermissionsColumnCreatedBy, pqcomp.E, c.createdBy)
	where.AddExpr(tableUserPermissionsColumnPermissionID, pqcomp.E, c.permissionID)
	if c.updatedAt.From != nil {
		where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.GT, c.updatedAt.From.Time())
	}
	if c.updatedAt.To != nil {
		where.AddExpr(tableUserPermissionsColumnUpdatedAt, pqcomp.LT, c.updatedAt.To.Time())
	}

	where.AddExpr(tableUserPermissionsColumnUpdatedBy, pqcomp.E, c.updatedBy)
	where.AddExpr(tableUserPermissionsColumnUserID, pqcomp.E, c.userID)

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
			&entity.PermissionID,
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
	insert := pqcomp.New(0, 6)
	insert.AddExpr(tableUserPermissionsColumnCreatedBy, "", e.CreatedBy)
	insert.AddExpr(tableUserPermissionsColumnPermissionID, "", e.PermissionID)
	insert.AddExpr(tableUserPermissionsColumnUpdatedAt, "", e.UpdatedAt)
	insert.AddExpr(tableUserPermissionsColumnUpdatedBy, "", e.UpdatedBy)
	insert.AddExpr(tableUserPermissionsColumnUserID, "", e.UserID)
	err := insertQueryComp(r.db, r.table, insert, r.columns).Scan(&e.CreatedAt,
		&e.CreatedBy,
		&e.PermissionID,
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
	CONSTRAINT "charon.user_groups_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id)
);

CREATE TABLE IF NOT EXISTS charon.group_permissions (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	group_id BIGINT NOT NULL,
	permission_id BIGINT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,

	CONSTRAINT "charon.group_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.group_permissions_group_id_fkey" FOREIGN KEY (group_id) REFERENCES charon.group (id),
	CONSTRAINT "charon.group_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charon.permission (id)
);

CREATE TABLE IF NOT EXISTS charon.user_permissions (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by BIGINT,
	permission_id BIGINT NOT NULL,
	updated_at TIMESTAMPTZ,
	updated_by BIGINT,
	user_id BIGINT NOT NULL,

	CONSTRAINT "charon.user_permissions_created_by_fkey" FOREIGN KEY (created_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_updated_by_fkey" FOREIGN KEY (updated_by) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_user_id_fkey" FOREIGN KEY (user_id) REFERENCES charon.user (id),
	CONSTRAINT "charon.user_permissions_permission_id_fkey" FOREIGN KEY (permission_id) REFERENCES charon.permission (id)
);

`
