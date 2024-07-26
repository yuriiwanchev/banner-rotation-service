package entities

type Slot struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Banner struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type UserGroup struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Event struct {
	Type        string `json:"type"`
	SlotID      string `json:"slot_id"`
	BannerID    string `json:"banner_id"`
	UserGroupID string `json:"user_group_id"`
}
