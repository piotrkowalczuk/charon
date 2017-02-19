package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/pqcomp"
)

// Message maps entity into protobuf message.
func (ge *GroupEntity) Message() (*charonrpc.Group, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ge.CreatedAt); err != nil {
		return nil, err
	}
	if ge.UpdatedAt.Valid {
		if updatedAt, err = ptypes.TimestampProto(ge.UpdatedAt.Time); err != nil {
			return nil, err
		}
	}

	return &charonrpc.Group{
		Id:          ge.ID,
		Name:        ge.Name,
		Description: ge.Description.Chars,
		CreatedAt:   createdAt,
		CreatedBy:   &ge.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   &ge.UpdatedBy,
	}, nil
}

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
	UpdateOneByID(ctx context.Context, id, updatedBy int64, name, description *ntypes.String) (*GroupEntity, error)
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
}

// NewGroupRepository ...
func NewGroupRepository(dbPool *sql.DB) GroupProvider {
	return &GroupRepository{
		GroupRepositoryBase: GroupRepositoryBase{
			DB:      dbPool,
			Table:   TableGroup,
			Columns: TableGroupColumns,
		},
	}
}

func (gr *GroupRepository) queryRow(ctx context.Context, query string, args ...interface{}) (*GroupEntity, error) {
	var entity GroupEntity
	err := gr.DB.QueryRowContext(ctx, query, args...).Scan(
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

// FindByUserID implements GroupProvider interface.
func (gr *GroupRepository) FindByUserID(ctx context.Context, userID int64) ([]*GroupEntity, error) {
	query := `
		SELECT  ` + strings.Join(TableGroupColumns, ",") + `
		FROM ` + TableGroup + ` AS g
		JOIN ` + TableUserGroups + ` AS ug ON ug.group_id = g.ID AND ug.user_id = $1
	`

	rows, err := gr.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*GroupEntity{}
	for rows.Next() {
		var g GroupEntity
		err = rows.Scan(
			&g.CreatedAt,
			&g.CreatedBy,
			&g.Description,
			&g.ID,
			&g.Name,
			&g.UpdatedAt,
			&g.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &g)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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

// UpdateOneByID ...
func (gr *GroupRepository) UpdateOneByID(ctx context.Context, id, updatedBy int64, name, description *ntypes.String) (*GroupEntity, error) {
	var (
		err    error
		query  string
		entity GroupEntity
	)

	comp := pqcomp.New(2, 2)
	comp.AddArg(id)
	comp.AddArg(updatedBy)
	comp.AddExpr("g.Name", pqcomp.Equal, name)
	comp.AddExpr("g.description", pqcomp.Equal, description)

	if comp.Len() == 0 {
		return nil, errors.New("nothing to update")
	}

	query = `UPDATE ` + TableGroup + ` SET `
	for comp.Next() {
		if !comp.First() {
			query += ", "
		}

		query += fmt.Sprintf("%s %s %s", comp.Key(), comp.Oper(), comp.PlaceHolder())
	}

	query += `
		, updated_by = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING ` + strings.Join(TableGroupColumns, ",") + `
	`

	err = gr.DB.QueryRowContext(ctx, query, comp.Args()).Scan(
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
