package statistic_repository

import (
	"database/sql"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type StatisticRepository interface {
	GetStatistics(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) (*e.Statistics, error)
	GetStatisticsForSlotAndBanner(slotID e.SlotID, bannerID e.BannerID) (*e.Statistics, error)
	UpdateStatistics(stat *e.Statistics) error
	IncrementClick(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) error
	IncrementView(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) error
}

type PgStatisticRepository struct {
	DB *sql.DB
}

func (r *PgStatisticRepository) GetStatistics(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) (*e.Statistics, error) {
	sql := `SELECT id, slot_id, banner_id, user_group_id, clicks, views 
			FROM statistics 
			WHERE slot_id = $1 
				AND banner_id = $2 
				AND user_group_id = $3`

	stat := &e.Statistics{}
	err := r.DB.QueryRow(sql, slotID, bannerID, userGroupID).Scan(&stat.ID, &stat.SlotID, &stat.BannerID,
		&stat.UserGroupID, &stat.Clicks, &stat.Views)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (r *PgStatisticRepository) GetStatisticsForSlotAndBanner(slotID e.SlotID, bannerID e.BannerID) (*e.Statistics, error) {
	sql := `SELECT id, slot_id, banner_id, user_group_id, clicks, views 
			FROM statistics 
			WHERE slot_id = $1 
				AND banner_id = $2`

	stat := &e.Statistics{}
	err := r.DB.QueryRow(sql, slotID, bannerID).Scan(&stat.ID, &stat.SlotID, &stat.BannerID,
		&stat.UserGroupID, &stat.Clicks, &stat.Views)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (r *PgStatisticRepository) UpdateStatistics(stat *e.Statistics) error {
	sql := `UPDATE statistics 
			SET clicks = $1, views = $2 
			WHERE slot_id = $3 AND banner_id = $4 AND user_group_id = $5`
	_, err := r.DB.Exec(sql,
		stat.Clicks, stat.Views, stat.SlotID, stat.BannerID, stat.UserGroupID)
	return err
}

func (r *PgStatisticRepository) IncrementClick(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) error {
	sql := `UPDATE statistics 
			SET clicks = clicks + 1 
			WHERE slot_id = $1 AND banner_id = $2 AND user_group_id = $3`
	_, err := r.DB.Exec(sql, slotID, bannerID, userGroupID)
	return err
}

func (r *PgStatisticRepository) IncrementView(slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) error {
	sql := `UPDATE statistics 
			SET views = views + 1 
			WHERE slot_id = $1 AND banner_id = $2 AND user_group_id = $3`
	_, err := r.DB.Exec(sql, slotID, bannerID, userGroupID)
	return err
}
