package models

import (
	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type AddBannerRequest struct {
	SlotID   e.SlotID   `json:"slotId"`
	BannerID e.BannerID `json:"bannerId"`
}

type RemoveBannerRequest struct {
	SlotID   e.SlotID   `json:"slotId"`
	BannerID e.BannerID `json:"bannerId"`
}

type RecordClickRequest struct {
	SlotID      e.SlotID      `json:"slotId"`
	BannerID    e.BannerID    `json:"bannerId"`
	UserGroupID e.UserGroupID `json:"userGroupId"`
}

type SelectBannerRequest struct {
	SlotID      e.SlotID      `json:"slotId"`
	UserGroupID e.UserGroupID `json:"userGroupId"`
}

type SelectBannerResponse struct {
	BannerID e.BannerID `json:"bannerId"`
}
