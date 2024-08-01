//go:build integration
// +build integration

package api_integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	brokers "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	api "github.com/yuriiwanchev/banner-rotation-service/internal/api"
	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	"github.com/yuriiwanchev/banner-rotation-service/internal/kafka"
	m "github.com/yuriiwanchev/banner-rotation-service/internal/models"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository"
)

var db *sql.DB
var reader *brokers.Reader

func TestMain(m *testing.M) {
	exec.Command("docker-compose", "-f", "../../docker-compose.integrational.yml", "up", "-d").Run()

	time.Sleep(1 * time.Second)

	repository.InitDB("postgres://user:password@localhost:5432/banner_rotation_db?sslmode=disable")
	repository.InitSchema()
	db = repository.GetDB()

	kafkaBrokers := "localhost:9092"
	kafkaTopic := "banner_events"

	api.InitKafkaProducer([]string{kafkaBrokers}, kafkaTopic)
	api.InitRepositories()

	reader = brokers.NewReader(brokers.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   kafkaTopic,
		GroupID: "consumer-group-id",
	})

	// tryBrocker([]string{kafkaBrokers}, kafkaTopic)

	// time.Sleep(5 * time.Second)

	// tryBrocker([]string{kafkaBrokers}, kafkaTopic)

	clearDatabase()
	fillDatabase()

	code := m.Run()

	reader.Close()
	repository.CloseDB()

	os.Exit(code)
}

func clearDatabase() error {
	tables := []string{"slot_banners", "statistics"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			return err
		}
	}
	api.InitRotationAlgorithm()
	return nil
}

func fillDatabase() error {
	_, err := db.Exec("INSERT INTO banners (id, description) VALUES (1, 'соки'), (2, 'сладости'), (3, 'мясо');")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO slots (id, description) VALUES (1, 'вверху'), (2, 'внизу'), (3, 'слева');")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_groups (id, description) VALUES (1, 'дети'), (2, 'взрослые');")
	if err != nil {
		return err
	}

	return nil
}

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

func tryBrocker(brokers []string, topic string) {
	event := e.Event{
		Type:        e.Click,
		SlotID:      e.SlotID(1),
		BannerID:    e.BannerID(1),
		UserGroupID: e.UserGroupID(1),
	}

	eventBytes, _ := json.Marshal(event)
	slotIDString := strconv.Itoa(1)
	slotIDBytes := []byte(slotIDString)

	kafkaProducer := kafka.NewKafkaProducer(brokers, topic)
	err := kafkaProducer.PublishMessage(slotIDBytes, eventBytes)

	if err != nil {
		log.Printf("Failed to publish message in tryBrocker: %v\n", err)
	}

	reader.ReadMessage(context.Background())
}

func bannerInSlotExists(slotID e.SlotID, bannerID e.BannerID) bool {
	sql := `SELECT EXISTS(
		SELECT 1
		FROM slot_banners sb
		WHERE sb.slot_id = $1
			AND sb.banner_id = $2
	);`

	var exists bool
	err := db.QueryRow(sql, slotID, bannerID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func TestAddBannerHandler(t *testing.T) {
	clearDatabase()

	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)

	exists := bannerInSlotExists(slotID, bannerID)
	assert.False(t, exists)

	sendAddBannerRequest(t, slotID, bannerID)

	exists = bannerInSlotExists(slotID, bannerID)
	assert.True(t, exists)
}

func TestAddBannerAlreadyExistsHandler(t *testing.T) {
	clearDatabase()

	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)

	exists := bannerInSlotExists(slotID, bannerID)
	assert.False(t, exists)

	_, err := db.Exec("INSERT INTO slot_banners (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
	if err != nil {
		t.Fatal(err)
	}

	exists = bannerInSlotExists(slotID, bannerID)
	assert.True(t, exists)

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
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRemoveBannerHandler(t *testing.T) {
	clearDatabase()

	slotID := e.SlotID(2)
	bannerID := e.BannerID(2)

	exists := bannerInSlotExists(slotID, bannerID)
	assert.False(t, exists)

	sendAddBannerRequest(t, slotID, bannerID)

	exists = bannerInSlotExists(slotID, bannerID)
	assert.True(t, exists)

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

	exists = bannerInSlotExists(slotID, bannerID)
	assert.False(t, exists)
}

func TestRemoveNonExistingBannerHandler(t *testing.T) {
	requestBody, _ := json.Marshal(m.RemoveBannerRequest{SlotID: 100, BannerID: 100})

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/remove-banner", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.RemoveBannerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
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

func getClicks(t *testing.T, slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) int {
	t.Helper()
	sql := `SELECT clicks
		FROM statistics
		WHERE slot_id = $1
			AND banner_id = $2
			AND user_group_id = $3;`

	var clicks int
	err := db.QueryRow(sql, slotID, bannerID, userGroupID).Scan(&clicks)
	if err != nil {
		t.Fatal(err)
	}
	return clicks
}

func getViews(t *testing.T, slotID e.SlotID, bannerID e.BannerID, userGroupID e.UserGroupID) int {
	t.Helper()
	sql := `SELECT views
		FROM statistics
		WHERE slot_id = $1
			AND banner_id = $2
			AND user_group_id = $3;`

	var views int
	err := db.QueryRow(sql, slotID, bannerID, userGroupID).Scan(&views)
	if err != nil {
		t.Fatal(err)
	}
	return views
}

func readEventFromKafka(t *testing.T) e.Event {
	t.Helper()
	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var event e.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		t.Fatal(err)
	}

	return event
}

func TestRecordClickHandler(t *testing.T) {
	clearDatabase()

	slotID := e.SlotID(2)
	bannerID := e.BannerID(2)
	userGroupID := e.UserGroupID(1)

	sendAddBannerRequest(t, slotID, bannerID)

	clicksBefore := getClicks(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 0, clicksBefore)

	sendRecordClickRequest(t, slotID, bannerID, userGroupID)

	clicksAfter := getClicks(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 1, clicksAfter)

	event := readEventFromKafka(t)

	assert.Equal(t, e.Click, event.Type)
	assert.Equal(t, slotID, event.SlotID)
	assert.Equal(t, bannerID, event.BannerID)
	assert.Equal(t, userGroupID, event.UserGroupID)
}

func sendSelectBannerRequest(t *testing.T, slotID e.SlotID, userGroupID e.UserGroupID) {
	t.Helper()
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

func TestSelectBannerHandler(t *testing.T) {
	clearDatabase()

	slotID := e.SlotID(1)
	bannerID := e.BannerID(1)
	userGroupID := e.UserGroupID(1)

	sendAddBannerRequest(t, slotID, bannerID)

	clicksBefore := getClicks(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 0, clicksBefore)

	sendRecordClickRequest(t, slotID, bannerID, userGroupID)

	clicksAfter := getClicks(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 1, clicksAfter)

	event := readEventFromKafka(t)

	assert.Equal(t, e.Click, event.Type)
	assert.Equal(t, slotID, event.SlotID)
	assert.Equal(t, bannerID, event.BannerID)
	assert.Equal(t, userGroupID, event.UserGroupID)

	viewsBefore := getViews(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 0, viewsBefore)

	sendSelectBannerRequest(t, slotID, userGroupID)

	viewsAfter := getViews(t, slotID, bannerID, userGroupID)
	assert.Equal(t, 1, viewsAfter)

	event = readEventFromKafka(t)

	assert.Equal(t, e.View, event.Type)
	assert.Equal(t, slotID, event.SlotID)
	assert.Equal(t, bannerID, event.BannerID)
	assert.Equal(t, userGroupID, event.UserGroupID)
}
