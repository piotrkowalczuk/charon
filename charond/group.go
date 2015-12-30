package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
	"github.com/piotrkowalczuk/pqcnstr"
	"github.com/piotrkowalczuk/pqcomp"
	"github.com/piotrkowalczuk/protot"
)

const (
	tableGroup                                                 = "charon.group"
	tableGroupConstraintPrimaryKey          pqcnstr.Constraint = tableGroup + "_pkey"
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
			CONSTRAINT "` + tableUserConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id),
			CONSTRAINT "` + tableUserConstraintForeignKeyUpdatedBy + `" FOREIGN KEY (updated_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tableGroupColumns = `
		id, name, description, created_at, created_by, updated_at, updated_by
	`
)

type groupEntity struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   *time.Time
	CreatedBy   int64
	UpdatedAt   *time.Time
	UpdatedBy   nilt.Int64
}

func (ge *groupEntity) Message() *charon.Group {
	var createdAt, updatedAt *protot.Timestamp

	if ge.CreatedAt != nil {
		createdAt = protot.TimeToTimestamp(*ge.CreatedAt)
	}

	if ge.UpdatedAt != nil {
		updatedAt = protot.TimeToTimestamp(*ge.UpdatedAt)
	}

	return &charon.Group{
		Id:          ge.ID,
		Name:        ge.Name,
		Description: ge.Description,
		CreatedAt:   createdAt,
		CreatedBy:   ge.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ge.UpdatedBy.Int64Or(0),
	}
}

// GroupRepository ...
type GroupRepository interface {
	// FindByUserID retrieves all groups for user represented by given id.
	FindByUserID(int64) ([]*groupEntity, error)
	// FindOneByID retrieves group for given id.
	FindOneByID(int64) (*groupEntity, error)
	// Create ...
	Create(createdBy int64, name, description string) (*groupEntity, error)
	// UpdateOneByID ...
	UpdateOneByID(id, updatedBy int64, name, description *nilt.String) (*groupEntity, error)
	// DeleteOneByID ...
	DeleteOneByID(id int64) (int64, error)
}

type groupRepository struct {
	db *sql.DB
}

func newGroupRepository(dbPool *sql.DB) GroupRepository {
	return &groupRepository{
		db: dbPool,
	}
}

func (gr *groupRepository) queryRow(query string, args ...interface{}) (*groupEntity, error) {
	var entity groupEntity
	err := gr.db.QueryRow(query, args...).Scan(
		&entity.ID,
		&entity.Name,
		&entity.Description,
		&entity.CreatedAt,
		&entity.CreatedBy,
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
		SELECT  ` + tableGroupColumns + `
		FROM ` + tableGroup + ` AS g
		JOIN ` + tableUserGroups + ` AS ug ON ug.group_id = g.id AND ug.user_id = $1
	`

	rows, err := gr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	groups := []*groupEntity{}
	for rows.Next() {
		var g groupEntity
		err = rows.Scan(
			&g.ID,
			&g.Name,
			&g.Description,
			&g.CreatedAt,
			&g.CreatedBy,
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
func (gr *groupRepository) Create(createdBy int64, name, description string) (*groupEntity, error) {
	entity := groupEntity{
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
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
			name, description, created_at, created_by
		)
		VALUES ($1, $2, NOW(), $3)
		RETURNING id, created_at
	`
	return gr.db.QueryRow(
		query,
		e.Name,
		e.Description,
		e.CreatedBy,
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
		RETURNING ` + tableGroupColumns + `
	`

	err = gr.db.QueryRow(query, comp.Args()).Scan(
		&entity.ID,
		&entity.Name,
		&entity.Description,
		&entity.CreatedAt,
		&entity.CreatedBy,
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

// FindOneByID implements GroupRepository interface.
func (gr *groupRepository) FindOneByID(id int64) (*groupEntity, error) {
	query := `
		SELECT  ` + tableGroupColumns + `
		FROM ` + tableGroup + ` AS g
		WHERE g.id = $1 LIMIT 1
	`

	return gr.queryRow(query, id)
}
