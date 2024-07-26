package entities

type SlotID int
type BannerID int
type UserGroupID int

type Slot struct {
	ID          SlotID `json:"id"`
	Description string `json:"description"`
}

type Banner struct {
	ID          BannerID `json:"id"`
	Description string   `json:"description"`
}

type UserGroup struct {
	ID          UserGroupID `json:"id"`
	Description string      `json:"description"`
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
