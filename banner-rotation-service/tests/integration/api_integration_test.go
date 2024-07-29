//go:build integration
// +build integration

package api_integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/kafka"
	brokers "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/yuriiwanchev/banner-rotation-service/internal/api"
	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	m "github.com/yuriiwanchev/banner-rotation-service/internal/models"
)

func sendAddBannerRequest(t *testing.T, slotID e.SlotID, bannerID e.BannerID) {
	t.Helper()
	addBannerRequest := m.AddBannerRequest{
		SlotID:   slotID,
		BannerID: bannerID,
	}
	requestBody, _ := json.Marshal(addBannerRequest)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/add-banner", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.AddBannerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAddBannerHandler(t *testing.T) {
	sendAddBannerRequest(t, 1, 1)
}

func TestRemoveBannerHandler(t *testing.T) {
	slotID := e.SlotID(2)
	bannerID := e.BannerID(2)

	sendAddBannerRequest(t, slotID, bannerID)

	requestBody, _ := json.Marshal(m.RemoveBannerRequest{SlotID: slotID, BannerID: bannerID})

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/remove-banner", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.RemoveBannerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func sendRecordClickRequest(t *testing.T, slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) {
	t.Helper()
	recordClickRequest := m.RecordClickRequest{
		SlotID:      slotID,
		BannerID:    bannerID,
		UserGroupID: userGroupID,
	}
	requestBody, _ := json.Marshal(recordClickRequest)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/record-click", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.RecordClickHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRecordClickHandler(t *testing.T) {
	kafkaTopic := "banner_events"

	container, err := gnomock.Start(
		kafka.Preset(kafka.WithTopics(kafkaTopic)),
		gnomock.WithDebugMode(), gnomock.WithLogWriter(os.Stdout),
		gnomock.WithContainerName("kafka"),
	)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, gnomock.Stop(container))
	}()

	kafkaBrokers := container.Address(kafka.BrokerPort)

	api.InitKafkaProducer([]string{kafkaBrokers}, "banner_events")

	reader := brokers.NewReader(brokers.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   kafkaTopic,
		GroupID: "consumer-group-id",
	})
	defer reader.Close()

	// Добавить баннер
	slotID := e.SlotID(2)
	bannerID := e.BannerID(2)
	sendAddBannerRequest(t, slotID, bannerID)

	userGroupID := e.UserGroupID(1)
	sendRecordClickRequest(t, slotID, bannerID, userGroupID)

	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var event e.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e.Click, event.Type)
	assert.Equal(t, slotID, event.SlotID)
	assert.Equal(t, bannerID, event.BannerID)
	assert.Equal(t, userGroupID, event.UserGroupID)
}

func TestSelectBannerHandler(t *testing.T) {
	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)
	sendAddBannerRequest(t, slotID, bannerID)

	userGroupID := e.UserGroupID(1)
	sendRecordClickRequest(t, slotID, bannerID, userGroupID)

	requestBody, _ := json.Marshal(m.SelectBannerRequest{SlotID: slotID, UserGroupID: userGroupID})

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/select-banner", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.SelectBannerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"bannerId":1}`
	assert.JSONEq(t, expected, rr.Body.String())
}
