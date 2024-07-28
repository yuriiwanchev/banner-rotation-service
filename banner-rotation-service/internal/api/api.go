package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	"github.com/yuriiwanchev/banner-rotation-service/internal/kafka"
	"github.com/yuriiwanchev/banner-rotation-service/internal/logic/bandit"
	m "github.com/yuriiwanchev/banner-rotation-service/internal/models"
)

var (
	banditService = bandit.NewMultiArmedBandit()
	kafkaProducer *kafka.Producer
)

func InitKafkaProducer(brokers []string, topic string) {
	kafkaProducer = kafka.NewKafkaProducer(brokers, topic)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func AddBannerHandler(w http.ResponseWriter, r *http.Request) {
	var request m.AddBannerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if request.SlotID == 0 || request.BannerID == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "SlotID and BannerID are required"})
		return
	}

	banditService.AddBanner(request.SlotID, request.BannerID)
	jsonResponse(w, http.StatusOK, nil)
}

func RemoveBannerHandler(w http.ResponseWriter, r *http.Request) {
	var request m.RemoveBannerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if request.SlotID == 0 || request.BannerID == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "SlotID and BannerID are required"})
		return
	}

	err := banditService.RemoveBanner(request.SlotID, request.BannerID)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

func RecordClickHandler(w http.ResponseWriter, r *http.Request) {
	var request m.RecordClickRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if request.SlotID == 0 || request.BannerID == 0 || request.UserGroupID == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "SlotID, BannerID, and UserGroup are required"})
		return
	}

	err := banditService.RecordClick(request.SlotID, request.BannerID, request.UserGroupID)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	event := e.Event{
		Type:        e.Click,
		SlotID:      request.SlotID,
		BannerID:    request.BannerID,
		UserGroupID: request.UserGroupID,
	}

	eventBytes, _ := json.Marshal(event)
	slotIDBytes := idToBytes(int(request.SlotID))
	kafkaProducer.PublishMessage(slotIDBytes, eventBytes)

	jsonResponse(w, http.StatusOK, nil)
}

func SelectBannerHandler(w http.ResponseWriter, r *http.Request) {
	var request m.SelectBannerRequest
	var response m.SelectBannerResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if request.SlotID == 0 || request.UserGroupID == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "SlotID and UserGroup are required"})
		return
	}

	response.BannerID = banditService.SelectBanner(request.SlotID, request.UserGroupID)
	if response.BannerID == 0 {
		jsonResponse(w, http.StatusNotFound,
			map[string]string{"error": "No banner available for the given slot and user group"})
		return
	}

	event := e.Event{
		Type:        e.View,
		SlotID:      request.SlotID,
		BannerID:    response.BannerID,
		UserGroupID: request.UserGroupID,
	}

	eventBytes, _ := json.Marshal(event)
	slotIDBytes := idToBytes(int(request.SlotID))
	kafkaProducer.PublishMessage(slotIDBytes, eventBytes)

	jsonResponse(w, http.StatusOK, response)
}

func idToBytes(id int) []byte {
	slotIDString := strconv.Itoa(id)
	idBytes := []byte(slotIDString)
	return idBytes
}
