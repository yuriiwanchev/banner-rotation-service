package entities

type (
	SlotID      int
	BannerID    int
	UserGroupID int
)

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
	SlotID      string    `json:"slotId"`
	BannerID    string    `json:"bannerId"`
	UserGroupID string    `json:"userGroupId"`
}

type EventType string

const (
	Click EventType = "Click"
	View  EventType = "View"
)
