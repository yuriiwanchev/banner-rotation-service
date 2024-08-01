package usergrouprepository

import (
	"database/sql"
	"fmt"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type UserGroupRepository interface {
	GetUserGroupByID(id e.UserGroupID) (*e.UserGroup, error)
	CreateUserGroup(group *e.UserGroup) (e.UserGroupID, error)
	GetAllUserGroupsIds() ([]e.UserGroupID, error)
}

type PgUserGroupRepository struct {
	DB *sql.DB
}

func (r *PgUserGroupRepository) GetUserGroupByID(id e.UserGroupID) (*e.UserGroup, error) {
	group := &e.UserGroup{}
	err := r.DB.QueryRow("SELECT id, description FROM user_groups WHERE id = $1", id).Scan(&group.ID, &group.Description)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *PgUserGroupRepository) CreateUserGroup(group *e.UserGroup) (e.UserGroupID, error) {
	var id e.UserGroupID
	err := r.DB.QueryRow("INSERT INTO user_groups (description) VALUES ($1) RETURNING id", group.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PgUserGroupRepository) GetAllUserGroupsIDs() ([]e.UserGroupID, error) {
	sql := `SELECT id FROM user_groups`
	rows, err := r.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []e.UserGroupID

	for rows.Next() {
		var id e.UserGroupID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return ids, nil
}
