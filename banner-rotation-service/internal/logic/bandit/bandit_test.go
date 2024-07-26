package bandit

import (
	"sync"
	"testing"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
)

func TestAddBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)

	mab.AddBanner(slotID, bannerID)

	if _, exists := mab.slots[slotID]; !exists {
		t.Errorf("Slot %d was not created", slotID)
	}

	if _, exists := mab.slots[slotID].Banners[bannerID]; !exists {
		t.Errorf("Banner %d was not added to slot %d", bannerID, slotID)
	}
}

func TestRemoveBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)

	mab.AddBanner(slotID, bannerID)
	mab.RemoveBanner(slotID, bannerID)

	if _, exists := mab.slots[slotID].Banners[bannerID]; exists {
		t.Errorf("Banner %d was not removed from slot %d", bannerID, slotID)
	}
}

func TestRecordClick(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID)
	mab.RecordClick(slotID, bannerID, groupID)

	if mab.slots[slotID].GroupData[groupID][bannerID].Clicks != 1 {
		t.Errorf("Click was not recorded for banner %d in slot %d for group %d", bannerID, slotID, groupID)
	}
}

func TestSelectBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID1 := e.BannerID(1)
	bannerID2 := e.BannerID(2)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID1)
	mab.AddBanner(slotID, bannerID2)

	selectedBanner := mab.SelectBanner(slotID, groupID)
	if selectedBanner != bannerID1 && selectedBanner != bannerID2 {
		t.Errorf("Selected banner %d is not one of the added banners", selectedBanner)
	}

	mab.RecordClick(slotID, selectedBanner, groupID)
	selectedBanner = mab.SelectBanner(slotID, groupID)
	if selectedBanner != bannerID1 && selectedBanner != bannerID2 {
		t.Errorf("Selected banner %d is not one of the added banners", selectedBanner)
	}
}

func TestSelectBanner_NoBanners(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	groupID := e.UserGroupID(1)

	selectedBanner := mab.SelectBanner(slotID, groupID)
	if selectedBanner != 0 {
		t.Errorf("Expected no banner to be selected, got %d", selectedBanner)
	}
}

// Edge Case Tests

func TestAddBanner_NonExistentSlot(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(2)
	bannerID := e.BannerID(3)

	mab.AddBanner(slotID, bannerID)
	if _, exists := mab.slots[slotID]; !exists {
		t.Errorf("Slot %d was not created for banner %d", slotID, bannerID)
	}
}

func TestRemoveBanner_NonExistentSlot(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(2)
	bannerID := e.BannerID(3)

	mab.RemoveBanner(slotID, bannerID)
	// Slot doesn't exist, no action should be taken
}

func TestRecordClick_NonExistentBannerOrSlot(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(3)
	bannerID := e.BannerID(4)
	groupID := e.UserGroupID(2)

	mab.RecordClick(slotID, bannerID, groupID)
	// No banner or slot exists, no action should be taken
}

func TestSelectBanner_EmptySlot(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(4)
	groupID := e.UserGroupID(3)

	selectedBanner := mab.SelectBanner(slotID, groupID)
	if selectedBanner != 0 {
		t.Errorf("Expected no banner to be selected from empty slot, got %d", selectedBanner)
	}
}

func TestSelectBanner_NoViews(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(5)
	bannerID1 := e.BannerID(1)
	bannerID2 := e.BannerID(2)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID1)
	mab.AddBanner(slotID, bannerID2)

	selectedBanner := mab.SelectBanner(slotID, groupID)
	if selectedBanner != bannerID1 && selectedBanner != bannerID2 {
		t.Errorf("Selected banner %d is not one of the added banners", selectedBanner)
	}
}

func TestSelectBanner_UCBAlgorithm(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(6)
	bannerID1 := e.BannerID(1)
	bannerID2 := e.BannerID(2)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID1)
	mab.AddBanner(slotID, bannerID2)

	// Simulate clicks and views
	mab.RecordClick(slotID, bannerID1, groupID)
	mab.RecordClick(slotID, bannerID2, groupID)
	mab.slots[slotID].GroupData[groupID][bannerID1].Views = 10
	mab.slots[slotID].GroupData[groupID][bannerID2].Views = 20

	selectedBanner := mab.SelectBanner(slotID, groupID)
	if selectedBanner != bannerID1 && selectedBanner != bannerID2 {
		t.Errorf("Selected banner %d is not one of the added banners", selectedBanner)
	}
}

// Concurrency Tests

func TestConcurrentAddBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bannerID := e.BannerID(i)
			mab.AddBanner(slotID, bannerID)
		}(i)
	}
	wg.Wait()

	if len(mab.slots[slotID].Banners) != 100 {
		t.Errorf("Expected 100 banners, but got %d", len(mab.slots[slotID].Banners))
	}
}

func TestConcurrentRemoveBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)

	for i := 0; i < 100; i++ {
		bannerID := e.BannerID(i)
		mab.AddBanner(slotID, bannerID)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bannerID := e.BannerID(i)
			mab.RemoveBanner(slotID, bannerID)
		}(i)
	}
	wg.Wait()

	if len(mab.slots[slotID].Banners) != 0 {
		t.Errorf("Expected 0 banners, but got %d", len(mab.slots[slotID].Banners))
	}
}

func TestConcurrentRecordClick(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mab.RecordClick(slotID, bannerID, groupID)
		}()
	}
	wg.Wait()

	if mab.slots[slotID].GroupData[groupID][bannerID].Clicks != 1000 {
		t.Errorf("Expected 1000 clicks, but got %d", mab.slots[slotID].GroupData[groupID][bannerID].Clicks)
	}
}

func TestConcurrentSelectBanner(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(1)
	bannerID1 := e.BannerID(1)
	bannerID2 := e.BannerID(2)
	groupID := e.UserGroupID(1)

	mab.AddBanner(slotID, bannerID1)
	mab.AddBanner(slotID, bannerID2)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mab.SelectBanner(slotID, groupID)
		}()
	}
	wg.Wait()

	totalViews := mab.slots[slotID].GroupData[groupID][bannerID1].Views + mab.slots[slotID].GroupData[groupID][bannerID2].Views
	if totalViews != 1000 {
		t.Errorf("Expected 1000 views, but got %d", totalViews)
	}
}

// Additional Tests for Specific Scenarios

func TestExhaustiveSelection(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(7)
	groupID := e.UserGroupID(1)

	// Add 10 banners
	for i := 1; i <= 10; i++ {
		mab.AddBanner(slotID, e.BannerID(i))
	}

	// Simulate 10000 views
	for i := 0; i < 10000; i++ {
		mab.SelectBanner(slotID, groupID)
	}

	// Check that each banner was shown at least once
	for i := 1; i <= 10; i++ {
		if mab.slots[slotID].GroupData[groupID][e.BannerID(i)].Views == 0 {
			t.Errorf("Banner %d was not shown even once", i)
		}
	}
}

func TestPopularBannerSelection(t *testing.T) {
	mab := NewMultiArmedBandit()
	slotID := e.SlotID(8)
	groupID := e.UserGroupID(1)

	bannerID1 := e.BannerID(1)
	bannerID2 := e.BannerID(2)
	bannerID3 := e.BannerID(3)

	mab.AddBanner(slotID, bannerID1)
	mab.AddBanner(slotID, bannerID2)
	mab.AddBanner(slotID, bannerID3)

	// Simulate 10000 views and clicks for bannerID1
	for i := 0; i < 10000; i++ {
		selectedBanner := mab.SelectBanner(slotID, groupID)
		if selectedBanner == bannerID1 {
			mab.RecordClick(slotID, bannerID1, groupID)
		}
	}

	views1 := mab.slots[slotID].GroupData[groupID][bannerID1].Views
	views2 := mab.slots[slotID].GroupData[groupID][bannerID2].Views
	views3 := mab.slots[slotID].GroupData[groupID][bannerID3].Views

	threshold := 0.5

	if !isSignificantlyBigger(views1, views2, threshold) ||
		!isSignificantlyBigger(views1, views3, threshold) {
		t.Errorf("Banner %d does not have significantly more views than others: %d, %d, %d", bannerID1, views1, views2, views3)
	}
}

func isSignificantlyBigger(num1, num2 int, threshold float64) bool {
	if num2 == 0 {
		return num1 != 0
	}

	relativeDifference := float64((num1 - num2)) / float64(num2)

	return relativeDifference >= threshold
}
