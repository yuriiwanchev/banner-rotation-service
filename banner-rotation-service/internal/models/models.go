package models

import (
	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type AddBannerRequest struct {
	SlotID   e.SlotID   `json:"slot_id"`
	BannerID e.BannerID `json:"banner_id"`
}

type RemoveBannerRequest struct {
	SlotID   e.SlotID   `json:"slot_id"`
	BannerID e.BannerID `json:"banner_id"`
}

type RecordClickRequest struct {
	SlotID      e.SlotID      `json:"slot_id"`
	BannerID    e.BannerID    `json:"banner_id"`
	UserGroupID e.UserGroupID `json:"user_group_id"`
}

type SelectBannerRequest struct {
	SlotID      e.SlotID      `json:"slot_id"`
	UserGroupID e.UserGroupID `json:"user_group_id"`
}

type SelectBannerResponse struct {
	BannerID e.BannerID `json:"banner_id"`
}
