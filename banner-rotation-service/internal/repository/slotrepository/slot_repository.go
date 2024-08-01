package slotrepository

import (
	"database/sql"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type SlotRepository interface {
	GetSlotByID(id e.SlotID) (*e.Slot, error)
	CreateSlot(slot *e.Slot) (e.SlotID, error)
	GetAllSlots() ([]*e.Slot, error)
}

type PgSlotRepository struct {
	DB *sql.DB
}

func (r *PgSlotRepository) GetSlotByID(id e.SlotID) (*e.Slot, error) {
	slot := &e.Slot{}
	err := r.DB.QueryRow("SELECT id, description FROM slots WHERE id = $1", id).Scan(&slot.ID, &slot.Description)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r *PgSlotRepository) CreateSlot(slot *e.Slot) (e.SlotID, error) {
	var id e.SlotID
	err := r.DB.QueryRow("INSERT INTO slots (description) VALUES ($1) RETURNING id", slot.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PgSlotRepository) GetAllSlots() ([]*e.Slot, error) {
	sql := `SELECT id, description FROM slots`
	rows, err := r.DB.Query(sql)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*e.Slot
	for rows.Next() {
		slot := &e.Slot{}
		if err := rows.Scan(&slot.ID, &slot.Description); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}
