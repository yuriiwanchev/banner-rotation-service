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
	Type        EventType   `json:"type"`
	SlotID      SlotID      `json:"slotId"`
	BannerID    BannerID    `json:"bannerId"`
	UserGroupID UserGroupID `json:"userGroupId"`
}

type EventType string

const (
	Click EventType = "Click"
	View  EventType = "View"
)

type Statistics struct {
	ID          int         `json:"id"`
	SlotID      SlotID      `json:"slotId"`
	BannerID    BannerID    `json:"bannerId"`
	UserGroupID UserGroupID `json:"userGroupId"`
	Clicks      int         `json:"clicks"`
	Views       int         `json:"views"`
}
