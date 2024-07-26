package entities

type Slot struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Banner struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type UserGroup struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Event struct {
	Type        EventType `json:"type"`
	SlotID      string    `json:"slot_id"`
	BannerID    string    `json:"banner_id"`
	UserGroupID string    `json:"user_group_id"`
}

type EventType string

const (
	Click EventType = "Click"
	View  EventType = "View"
)
