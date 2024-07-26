package models

type AddBannerRequest struct {
	SlotID   string `json:"slot_id"`
	BannerID string `json:"banner_id"`
}

type RemoveBannerRequest struct {
	SlotID   string `json:"slot_id"`
	BannerID string `json:"banner_id"`
}

type RecordClickRequest struct {
	SlotID    string `json:"slot_id"`
	BannerID  string `json:"banner_id"`
	UserGroup string `json:"user_group"`
}

type RecordViewRequest struct {
	SlotID    string `json:"slot_id"`
	BannerID  string `json:"banner_id"`
	UserGroup string `json:"user_group"`
}

type SelectBannerRequest struct {
	SlotID      string `json:"slot_id"`
	UserGroupID string `json:"user_group_id"`
}

type SelectBannerResponse struct {
	BannerID string `json:"banner_id"`
}
