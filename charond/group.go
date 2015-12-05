package main

import (
	"database/sql"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqcnstr"
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
	UpdatedBy   int64
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
		UpdatedBy:   ge.UpdatedBy,
	}
}

type GroupRepository interface {
	FindByUserID(int64) ([]*groupEntity, error)
}

type groupRepository struct {
	db *sql.DB
}

func newGroupRepository(dbPool *sql.DB) GroupRepository {
	return &groupRepository{
		db: dbPool,
	}
}

// FindOneByID retrieves all permissions for user represented by given id.
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
