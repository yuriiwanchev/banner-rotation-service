package bandit

import (
	"math/rand"

	"github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

type BannerStats struct {
	Clicks int
	Views  int
}

type MultiArmedBandit struct {
	Epsilon     float64
	BannerStats map[entities.SlotID]map[entities.BannerID]*BannerStats
}

func NewMultiArmedBandit(epsilon float64) *MultiArmedBandit {
	return &MultiArmedBandit{
		Epsilon:     epsilon,
		BannerStats: make(map[entities.SlotID]map[entities.BannerID]*BannerStats),
	}
}

func (mab *MultiArmedBandit) AddBanner(slotID entities.SlotID, bannerID entities.BannerID) {
	if _, exists := mab.BannerStats[slotID]; !exists {
		mab.BannerStats[slotID] = make(map[entities.BannerID]*BannerStats)
	}
	mab.BannerStats[slotID][bannerID] = &BannerStats{Clicks: 0, Views: 0}
}

func (mab *MultiArmedBandit) RemoveBanner(slotID entities.SlotID, bannerID entities.BannerID) {
	if _, exists := mab.BannerStats[slotID]; exists {
		delete(mab.BannerStats[slotID], bannerID)
	}
}

func (mab *MultiArmedBandit) RecordClick(slotID entities.SlotID, bannerID entities.BannerID) {
	if stats, exists := mab.BannerStats[slotID][bannerID]; exists {
		stats.Clicks++
	}
}

func (mab *MultiArmedBandit) RecordView(slotID entities.SlotID, bannerID entities.BannerID) {
	if stats, exists := mab.BannerStats[slotID][bannerID]; exists {
		stats.Views++
	}
}

func (mab *MultiArmedBandit) SelectBanner(slotID entities.SlotID) entities.BannerID {
	banners := mab.BannerStats[slotID]
	if rand.Float64() < mab.Epsilon {
		// Выбираем случайный баннер
		for bannerID := range banners {
			return bannerID
		}
	}

	// Выбираем баннер с максимальным CTR (click-through rate)
	var selectedBanner entities.BannerID
	maxCTR := -1.0
	for bannerID, stats := range banners {
		ctr := float64(stats.Clicks) / float64(stats.Views+1)
		if ctr > maxCTR {
			maxCTR = ctr
			selectedBanner = bannerID
		}
	}
	return selectedBanner
}
