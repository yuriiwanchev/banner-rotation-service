package slotbannersrepository

import (
	"database/sql"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository"
)

type SlotBannerRepository interface {
	AddBannerToSlot(slotID e.SlotID, bannerID e.BannerID) error
	RemoveBannerFromSlot(slotID e.SlotID, bannerID e.BannerID) error
	GetBannersForSlot(slotID e.SlotID) ([]*e.Banner, error)
}

type PgSlotBannerRepository struct {
	DB *sql.DB
}

func (r *PgSlotBannerRepository) AddBannerToSlot(slotID e.SlotID, bannerID e.BannerID) error {
	db := repository.GetDB()

	_, err := db.Exec("INSERT INTO slot_banners (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
	return err
}

func (r *PgSlotBannerRepository) RemoveBannerFromSlot(slotID e.SlotID, bannerID e.BannerID) error {
	_, err := r.DB.Exec("DELETE FROM slot_banners WHERE slot_id = $1 AND banner_id = $2", slotID, bannerID)
	return err
}

func (r *PgSlotBannerRepository) GetBannersForSlot(slotID e.SlotID) ([]*e.Banner, error) {
	sql := `SELECT b.id, b.description 
			FROM banners b 
			INNER JOIN slot_banners sb ON b.id = sb.banner_id 
			WHERE sb.slot_id = $1`
	rows, err := r.DB.Query(sql, slotID)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*e.Banner
	for rows.Next() {
		banner := &e.Banner{}
		if err := rows.Scan(&banner.ID, &banner.Description); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
