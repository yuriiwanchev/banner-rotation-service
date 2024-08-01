package bandit

import (
	"fmt"
	"log"
	"math"
	"sync"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type GroupStats struct {
	Views  int
	Clicks int
}

type Slot struct {
	Banners   map[e.BannerID]e.Banner
	GroupData map[e.UserGroupID]map[e.BannerID]*GroupStats
}

type MultiArmedBandit struct {
	slots map[e.SlotID]*Slot
	mu    sync.Mutex
}

func NewMultiArmedBandit(slots map[e.SlotID]*Slot) *MultiArmedBandit {
	return &MultiArmedBandit{
		slots: slots,
	}
}

func (mab *MultiArmedBandit) AddBanner(slotID e.SlotID, bannerID e.BannerID) {
	mab.mu.Lock()
	defer mab.mu.Unlock()

	slot, exists := mab.slots[slotID]
	if !exists {
		slot = &Slot{
			Banners:   make(map[e.BannerID]e.Banner),
			GroupData: make(map[e.UserGroupID]map[e.BannerID]*GroupStats),
		}
		mab.slots[slotID] = slot
	}

	slot.Banners[bannerID] = e.Banner{ID: bannerID}
}

func (mab *MultiArmedBandit) RemoveBanner(slotID e.SlotID, bannerID e.BannerID) error {
	mab.mu.Lock()
	defer mab.mu.Unlock()

	slot, exists := mab.slots[slotID]
	if !exists {
		return fmt.Errorf("slot %d does not exist", slotID)
	}

	delete(slot.Banners, bannerID)
	for _, groupStats := range slot.GroupData {
		delete(groupStats, bannerID)
	}

	return nil
}

func (mab *MultiArmedBandit) RecordClick(slotID e.SlotID, bannerID e.BannerID, groupID e.UserGroupID) error {
	mab.mu.Lock()
	defer mab.mu.Unlock()

	slot, exists := mab.slots[slotID]
	if !exists {
		return fmt.Errorf("slot %d does not exist", slotID)
	}

	groupStats, exists := slot.GroupData[groupID]
	if !exists {
		groupStats = make(map[e.BannerID]*GroupStats)
		slot.GroupData[groupID] = groupStats
	}

	stats, exists := groupStats[bannerID]
	if !exists {
		stats = &GroupStats{}
		groupStats[bannerID] = stats
	}

	stats.Clicks++

	return nil
}

func (mab *MultiArmedBandit) SelectBanner(slotID e.SlotID, groupID e.UserGroupID) e.BannerID {
	mab.mu.Lock()
	defer mab.mu.Unlock()

	slot, exists := mab.slots[slotID]
	if !exists {
		log.Printf("SelectBanner: slot %d does not exist", slotID)
		return 0
	}

	groupStats, exists := slot.GroupData[groupID]
	if !exists {
		groupStats = make(map[e.BannerID]*GroupStats)
		slot.GroupData[groupID] = groupStats
	}

	var selectedBanner e.BannerID
	maxUCB := -1.0

	for bannerID := range slot.Banners {
		stats, exists := groupStats[bannerID]
		if !exists {
			stats = &GroupStats{}
			groupStats[bannerID] = stats
		}

		ucb := calculateUCB(stats.Clicks, stats.Views, totalViews(groupStats))

		if ucb > maxUCB {
			maxUCB = ucb
			selectedBanner = bannerID
		}
	}

	if selectedBanner != 0 {
		groupStats[selectedBanner].Views++
	}

	return selectedBanner
}

func calculateUCB(clicks, views, totalViews int) float64 {
	if views == 0 {
		return 1e6
	}
	return float64(clicks)/float64(views) + 2.0*math.Sqrt(math.Log(float64(totalViews))/float64(views))
}

func totalViews(groupStats map[e.BannerID]*GroupStats) int {
	total := 0
	for _, stats := range groupStats {
		total += stats.Views
	}
	return total
}
