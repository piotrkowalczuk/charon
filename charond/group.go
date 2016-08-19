package charond

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/pqcomp"
)

func (ge *groupEntity) message() (*charon.Group, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ge.createdAt); err != nil {
		return nil, err
	}
	if ge.updatedAt != nil {
		if updatedAt, err = ptypes.TimestampProto(*ge.updatedAt); err != nil {
			return nil, err
		}
	}

	return &charon.Group{
		Id:          ge.id,
		Name:        ge.name,
		Description: ge.description.StringOr(""),
		CreatedAt:   createdAt,
		CreatedBy:   ge.createdBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ge.updatedBy,
	}, nil
}

type groupProvider interface {
	insert(entity *groupEntity) (*groupEntity, error)
	// findByUserID retrieves all groups for user represented by given id.
	findByUserID(int64) ([]*groupEntity, error)
	// findOneByID retrieves group for given id.
	findOneByID(int64) (*groupEntity, error)
	// find ...
	find(c *groupCriteria) ([]*groupEntity, error)
	// Create ...
	create(createdBy int64, name string, description *ntypes.String) (*groupEntity, error)
	// updateOneByID ...
	updateOneByID(id, updatedBy int64, name, description *ntypes.String) (*groupEntity, error)
	// DeleteByID ...
	deleteOneByID(id int64) (int64, error)
	// IsGranted ...
	isGranted(id int64, permission charon.Permission) (bool, error)
	// SetPermissions ...
	setPermissions(id int64, permissions ...charon.Permission) (int64, int64, error)
}

type groupRepository struct {
	groupRepositoryBase
}

func newGroupRepository(dbPool *sql.DB) groupProvider {
	return &groupRepository{
		groupRepositoryBase: groupRepositoryBase{
			db:      dbPool,
			table:   tableGroup,
			columns: tableGroupColumns,
		},
	}
}

func (gr *groupRepository) queryRow(query string, args ...interface{}) (*groupEntity, error) {
	var entity groupEntity
	err := gr.db.QueryRow(query, args...).Scan(
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

// findByUserID implements GroupRepository interface.
func (gr *groupRepository) findByUserID(userID int64) ([]*groupEntity, error) {
	query := `
		SELECT  ` + strings.Join(tableGroupColumns, ",") + `
		FROM ` + tableGroup + ` AS g
		JOIN ` + tableUserGroups + ` AS ug ON ug.group_id = g.id AND ug.user_id = $1
	`

	rows, err := gr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*groupEntity{}
	for rows.Next() {
		var g groupEntity
		err = rows.Scan(
			&g.createdAt,
			&g.createdBy,
			&g.description,
			&g.id,
			&g.name,
			&g.updatedAt,
			&g.updatedBy,
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

func (gr *groupRepository) create(createdBy int64, name string, description *ntypes.String) (ent *groupEntity, err error) {
	if description == nil {
		description = &ntypes.String{}
	}
	ent = &groupEntity{
		name:        name,
		description: description,
		createdBy:   &ntypes.Int64{Int64: createdBy, Valid: createdBy > 0},
	}

	return gr.insert(ent)
}

func (gr *groupRepository) updateOneByID(id, updatedBy int64, name, description *ntypes.String) (*groupEntity, error) {
	var (
		err    error
		query  string
		entity groupEntity
	)

	comp := pqcomp.New(2, 2)
	comp.AddArg(id)
	comp.AddArg(updatedBy)
	comp.AddExpr("g.name", pqcomp.Equal, name)
	comp.AddExpr("g.description", pqcomp.Equal, description)

	if comp.Len() == 0 {
		return nil, errors.New("nothing to update")
	}

	query = `UPDATE ` + tableGroup + ` SET `
	for comp.Next() {
		if !comp.First() {
			query += ", "
		}

		query += fmt.Sprintf("%s %s %s", comp.Key(), comp.Oper(), comp.PlaceHolder())
	}

	query += `
		, updated_by = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING ` + strings.Join(tableGroupColumns, ",") + `
	`

	err = gr.db.QueryRow(query, comp.Args()).Scan(
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

func (gr *groupRepository) isGranted(id int64, p charon.Permission) (bool, error) {
	var exists bool
	subsystem, module, action := p.Split()
	if err := gr.db.QueryRow(isGrantedQuery(
		tableGroupPermissions,
		tableGroupPermissionsColumnGroupID,
		tableGroupPermissionsColumnPermissionSubsystem,
		tableGroupPermissionsColumnPermissionModule,
		tableGroupPermissionsColumnPermissionAction,
	), id, subsystem, module, action).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (gr *groupRepository) setPermissions(id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(gr.db, tableGroupPermissions,
		tableUserPermissionsColumnUserID,
		tableUserPermissionsColumnPermissionSubsystem,
		tableUserPermissionsColumnPermissionModule,
		tableUserPermissionsColumnPermissionAction, id, p)
}
