package model

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

// GroupProvider ...
type GroupProvider interface {
	Insert(context.Context, *GroupEntity) (*GroupEntity, error)
	// FindByUserID retrieves all groups for user represented by given id.
	FindByUserID(context.Context, int64) ([]*GroupEntity, error)
	// FindOneByID retrieves group for given id.
	FindOneByID(context.Context, int64) (*GroupEntity, error)
	// find ...
	Find(context.Context, *GroupFindExpr) ([]*GroupEntity, error)
	// Create ...
	Create(ctx context.Context, createdBy int64, name string, description *ntypes.String) (*GroupEntity, error)
	// updateOneByID ...
	UpdateOneByID(context.Context, int64, *GroupPatch) (*GroupEntity, error)
	// DeleteByID ...
	DeleteOneByID(context.Context, int64) (int64, error)
	// IsGranted ...
	IsGranted(context.Context, int64, charon.Permission) (bool, error)
	// SetPermissions ...
	SetPermissions(context.Context, int64, ...charon.Permission) (int64, int64, error)
}

// GroupRepository extends GroupRepositoryBase
type GroupRepository struct {
	GroupRepositoryBase
	UserGroups UserGroupsRepositoryBase
}

// NewGroupRepository ...
func NewGroupRepository(dbPool *sql.DB) GroupProvider {
	return &GroupRepository{
		GroupRepositoryBase: GroupRepositoryBase{
			DB:      dbPool,
			Table:   TableGroup,
			Columns: TableGroupColumns,
		},
		UserGroups: UserGroupsRepositoryBase{
			DB:      dbPool,
			Table:   TableUserGroups,
			Columns: TableUserGroupsColumns,
		},
	}
}

// FindByUserID implements GroupProvider interface.
func (gr *GroupRepository) FindByUserID(ctx context.Context, userID int64) ([]*GroupEntity, error) {
	userGroups, err := gr.UserGroups.Find(ctx, &UserGroupsFindExpr{
		//Columns: TableGroupColumns,
		Where: &UserGroupsCriteria{
			UserID: qtypes.EqualInt64(userID),
		},
		JoinGroup: &GroupJoin{
			Fetch: true,
			Kind:  JoinRight,
		},
	})
	if err != nil {
		return nil, err
	}

	groups := make([]*GroupEntity, 0, len(userGroups))
	for _, userGroup := range userGroups {
		groups = append(groups, userGroup.Group)
	}

	return groups, nil
}

// Create ...
func (gr *GroupRepository) Create(ctx context.Context, createdBy int64, name string, description *ntypes.String) (ent *GroupEntity, err error) {
	if description == nil {
		description = &ntypes.String{}
	}
	ent = &GroupEntity{
		Name:        name,
		Description: *description,
		CreatedBy:   ntypes.Int64{Int64: createdBy, Valid: createdBy > 0},
	}

	return gr.Insert(ctx, ent)
}

// IsGranted ...
func (gr *GroupRepository) IsGranted(ctx context.Context, id int64, p charon.Permission) (bool, error) {
	var exists bool
	subsystem, module, action := p.Split()
	if err := gr.DB.QueryRowContext(ctx, isGrantedQuery(
		TableGroupPermissions,
		TableGroupPermissionsColumnGroupID,
		TableGroupPermissionsColumnPermissionSubsystem,
		TableGroupPermissionsColumnPermissionModule,
		TableGroupPermissionsColumnPermissionAction,
	), id, subsystem, module, action).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// SetPermissions ...
func (gr *GroupRepository) SetPermissions(ctx context.Context, id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(gr.DB, ctx, TableGroupPermissions,
		TableGroupPermissionsColumnGroupID,
		TableGroupPermissionsColumnPermissionSubsystem,
		TableGroupPermissionsColumnPermissionModule,
		TableGroupPermissionsColumnPermissionAction, id, p)
}
