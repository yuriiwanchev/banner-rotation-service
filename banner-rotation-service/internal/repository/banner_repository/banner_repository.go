package banner_repository

import (
	"database/sql"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository"
)

type BannerRepository interface {
	GetBannerByID(id e.BannerID) (*e.Banner, error)
	CreateBanner(banner *e.Banner) (e.BannerID, error)
}

type PgBannerRepository struct {
	DB *sql.DB
}

func (r *PgBannerRepository) GetBannerByID(id e.BannerID) (*e.Banner, error) {
	db := repository.GetDB()

	banner := &e.Banner{}
	err := db.QueryRow("SELECT id, description FROM banners WHERE id = $1", id).Scan(&banner.ID, &banner.Description)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *PgBannerRepository) CreateBanner(banner *e.Banner) (e.BannerID, error) {
	db := repository.GetDB()

	var id e.BannerID
	err := db.QueryRow("INSERT INTO banners (description) VALUES ($1) RETURNING id", banner.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
