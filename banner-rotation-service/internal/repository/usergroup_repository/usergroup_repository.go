package usergroup_repository

import (
	"database/sql"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type UserGroupRepository interface {
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
