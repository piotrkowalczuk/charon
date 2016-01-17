package main

import (
	"time"

	"github.com/piotrkowalczuk/nilt"
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
	tableUserConstraintUsernameUnique      = "charon.user_username_key"
	tableUserConstraintPrimaryKey          = "charon.user_id_pkey"
	tableUserConstraintCreatedByForeignKey = "charon.user_created_by_fkey"
	tableUserConstraintUpdatedByForeignKey = "charon.user_updated_by_fkey"
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
	Password          []byte
	Username          string
	FirstName         string
	LastName          string
	IsSuperuser       bool
	IsActive          bool
	IsStaff           bool
	IsConfirmed       bool
	ConfirmationToken []byte
	LastLoginAt       *time.Time
	ID                int64
	CreatedBy         nilt.Int64
	UpdatedBy         nilt.Int64
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	Author            *userEntity
	Modifier          *userEntity
	Group             []*groupEntity
	Permission        []*permissionEntity
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
	tableGroupConstraintNameUnique          = "charon.group_name_key"
	tableGroupConstraintPrimaryKey          = "charon.group_id_pkey"
	tableGroupConstraintCreatedByForeignKey = "charon.group_created_by_fkey"
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
	Name        string
	Description nilt.String
	ID          int64
	CreatedBy   nilt.Int64
	UpdatedBy   nilt.Int64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Author      *groupEntity
	Modifier    *groupEntity
	Users       []*userEntity
	Permission  []*permissionEntity
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
	Subsystem string
	Module    string
	Action    string
	ID        int64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Groups    []*groupEntity
	Users     []*userEntity
}

const (
	tableUserGroups                              = "charon.user_groups"
	tableUserGroupsColumnCreatedAt               = "created_at"
	tableUserGroupsColumnCreatedBy               = "created_by"
	tableUserGroupsColumnGroupID                 = "group_id"
	tableUserGroupsColumnUpdatedAt               = "updated_at"
	tableUserGroupsColumnUpdatedBy               = "updated_by"
	tableUserGroupsColumnUserID                  = "user_id"
	tableUserGroupsConstraintUserIDForeignKey    = "charon.user_groups_user_id_fkey"
	tableUserGroupsConstraintGroupIDForeignKey   = "charon.user_groups_group_id_fkey"
	tableUserGroupsConstraintCreatedByForeignKey = "charon.user_groups_created_by_fkey"
	tableUserGroupsConstraintUpdatedByForeignKey = "charon.user_groups_updated_by_fkey"
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
	UserID    int64
	GroupID   int64
	CreatedBy nilt.Int64
	UpdatedBy nilt.Int64
	CreatedAt time.Time
	UpdatedAt *time.Time
	User      *userEntity
	Group     *groupEntity
	Author    *userGroupsEntity
	Modifier  *userGroupsEntity
}

const (
	tableGroupPermissions                                 = "charon.group_permissions"
	tableGroupPermissionsColumnCreatedAt                  = "created_at"
	tableGroupPermissionsColumnCreatedBy                  = "created_by"
	tableGroupPermissionsColumnGroupID                    = "group_id"
	tableGroupPermissionsColumnPermissionID               = "permission_id"
	tableGroupPermissionsColumnUpdatedAt                  = "updated_at"
	tableGroupPermissionsColumnUpdatedBy                  = "updated_by"
	tableGroupPermissionsConstraintGroupIDForeignKey      = "charon.group_permissions_group_id_fkey"
	tableGroupPermissionsConstraintPermissionIDForeignKey = "charon.group_permissions_permission_id_fkey"
	tableGroupPermissionsConstraintCreatedByForeignKey    = "charon.group_permissions_created_by_fkey"
	tableGroupPermissionsConstraintUpdatedByForeignKey    = "charon.group_permissions_updated_by_fkey"
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
	GroupID      int64
	PermissionID int64
	CreatedBy    nilt.Int64
	UpdatedBy    nilt.Int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Group        *groupEntity
	Permission   *permissionEntity
	Author       *groupPermissionsEntity
	Modifier     *groupPermissionsEntity
}

const (
	tableUserPermissions                                 = "charon.user_permissions"
	tableUserPermissionsColumnCreatedAt                  = "created_at"
	tableUserPermissionsColumnCreatedBy                  = "created_by"
	tableUserPermissionsColumnPermissionID               = "permission_id"
	tableUserPermissionsColumnUpdatedAt                  = "updated_at"
	tableUserPermissionsColumnUpdatedBy                  = "updated_by"
	tableUserPermissionsColumnUserID                     = "user_id"
	tableUserPermissionsConstraintUserIDForeignKey       = "charon.user_permissions_user_id_fkey"
	tableUserPermissionsConstraintPermissionIDForeignKey = "charon.user_permissions_permission_id_fkey"
	tableUserPermissionsConstraintCreatedByForeignKey    = "charon.user_permissions_created_by_fkey"
	tableUserPermissionsConstraintUpdatedByForeignKey    = "charon.user_permissions_updated_by_fkey"
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
	UserID       int64
	PermissionID int64
	CreatedBy    nilt.Int64
	UpdatedBy    nilt.Int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	User         *userEntity
	Permission   *permissionEntity
	Author       *userPermissionsEntity
	Modifier     *userPermissionsEntity
}
