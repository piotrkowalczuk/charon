package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
	"github.com/piotrkowalczuk/pqcnstr"
	"github.com/piotrkowalczuk/pqcomp"
	"github.com/piotrkowalczuk/protot"
)

const (
	//	tableGroup                                                 = "charon.group"

	tableGroupConstraintUniqueName          pqcnstr.Constraint = tableGroup + "_name_key"
	tableGroupConstraintForeignKeyCreatedBy pqcnstr.Constraint = tableGroup + "_created_by_fkey"
	tableGroupConstraintForeignKeyUpdatedBy pqcnstr.Constraint = tableGroup + "_updated_by_fkey"
	tableGroupCreate                                           = `
		CREATE TABLE IF NOT EXISTS ` + tableGroup + ` (
			id          SERIAL,
			name        TEXT                      NOT NULL,
			description TEXT,
			created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by  INTEGER,
			updated_at  TIMESTAMPTZ,
			updated_by  INTEGER,

			CONSTRAINT "` + tableGroupConstraintPrimaryKey + `" PRIMARY KEY (id),
			CONSTRAINT "` + tableGroupConstraintUniqueName + `" UNIQUE (name),
			CONSTRAINT "` + tableGroupConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableGroup + ` (id),
			CONSTRAINT "` + tableGroupConstraintForeignKeyUpdatedBy + `" FOREIGN KEY (updated_by) REFERENCES ` + tableGroup + ` (id)
		)
	`
)

func (ge *groupEntity) Message() *charon.Group {
	var createdAt, updatedAt *protot.Timestamp

	createdAt = protot.TimeToTimestamp(ge.CreatedAt)
	if ge.UpdatedAt != nil {
		updatedAt = protot.TimeToTimestamp(*ge.UpdatedAt)
	}

	return &charon.Group{
		Id:          ge.ID,
		Name:        ge.Name,
		Description: ge.Description.String,
		CreatedAt:   createdAt,
		CreatedBy:   &ge.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   &ge.UpdatedBy,
	}
}

// GroupRepository ...
type GroupRepository interface {
	// FindByUserID retrieves all groups for user represented by given id.
	FindByUserID(int64) ([]*groupEntity, error)
	// FindOneByID retrieves group for given id.
	FindOneByID(int64) (*groupEntity, error)
	// Find ...
	Find(c *groupCriteria) ([]*groupEntity, error)
	// Create ...
	Create(createdBy int64, name string, description *nilt.String) (*groupEntity, error)
	// UpdateOneByID ...
	UpdateOneByID(id, updatedBy int64, name, description *nilt.String) (*groupEntity, error)
	// DeleteOneByID ...
	DeleteOneByID(id int64) (int64, error)
}

func newGroupRepository(dbPool *sql.DB) GroupRepository {
	return &groupRepository{
		db:      dbPool,
		table:   tableGroup,
		columns: tableGroupColumns,
	}
}

func (gr *groupRepository) queryRow(query string, args ...interface{}) (*groupEntity, error) {
	var entity groupEntity
	err := gr.db.QueryRow(query, args...).Scan(
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

// FindByUserID implements GroupRepository interface.
func (gr *groupRepository) FindByUserID(userID int64) ([]*groupEntity, error) {
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

// Create implements GroupRepository interface.
func (gr *groupRepository) Create(createdBy int64, name string, description *nilt.String) (*groupEntity, error) {
	if description == nil {
		description = &nilt.String{}
	}
	entity := groupEntity{
		Name:        name,
		Description: *description,
		CreatedBy:   nilt.Int64{Int64: createdBy, Valid: createdBy > 0},
	}

	err := gr.insert(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (gr *groupRepository) insert(e *groupEntity) error {
	query := `
		INSERT INTO ` + tableGroup + ` (
			name string, description, created_at, created_by
		)
		VALUES ($1, $2, NOW(), $3)
		RETURNING id, created_at
	`
	return gr.db.QueryRow(
		query,
		e.CreatedBy,
		e.Description,
		e.Name,
	).Scan(&e.ID, &e.CreatedAt)
}

// UpdateOneByID implements GroupRepository interface.
func (gr *groupRepository) UpdateOneByID(id, updatedBy int64, name, description *nilt.String) (*groupEntity, error) {
	var (
		err    error
		entity groupEntity
		query  string
	)

	comp := pqcomp.New(2, 2)
	comp.AddArg(id)
	comp.AddArg(updatedBy)
	comp.AddExpr("g.name", pqcomp.E, name)
	comp.AddExpr("g.description", pqcomp.E, description)

	if comp.Len() == 0 {
		return nil, errors.New("charond: nothing to update")
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

// DeleteByUserID implements GroupRepository interface.
func (gr *groupRepository) DeleteOneByID(id int64) (int64, error) {
	query := `
		DELETE FROM ` + tableGroup + `
		WHERE id = $1
	`

	res, err := gr.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
