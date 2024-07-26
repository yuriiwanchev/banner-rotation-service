package models

type AddBannerRequest struct {
	SlotID   int `json:"slot_id"`
	BannerID int `json:"banner_id"`
}

type RemoveBannerRequest struct {
	SlotID   int `json:"slot_id"`
	BannerID int `json:"banner_id"`
}

type RecordClickRequest struct {
	SlotID      int `json:"slot_id"`
	BannerID    int `json:"banner_id"`
	UserGroupID int `json:"user_group_id"`
}

type RecordViewRequest struct {
	SlotID      int `json:"slot_id"`
	BannerID    int `json:"banner_id"`
	UserGroupID int `json:"user_group_id"`
}

type SelectBannerRequest struct {
	SlotID      int `json:"slot_id"`
	UserGroupID int `json:"user_group_id"`
}

type SelectBannerResponse struct {
	BannerID string `json:"banner_id"`
}
