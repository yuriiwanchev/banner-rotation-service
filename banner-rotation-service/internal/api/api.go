package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	e "github.com/yuriiwanchev/banner-rotation-service/internal/entities"
	"github.com/yuriiwanchev/banner-rotation-service/internal/kafka"
	"github.com/yuriiwanchev/banner-rotation-service/internal/logic/bandit"
	m "github.com/yuriiwanchev/banner-rotation-service/internal/models"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository/slotbannersrepository"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository/statisticrepository"
	"github.com/yuriiwanchev/banner-rotation-service/internal/repository/usergrouprepository"
)

var (
	banditService         *bandit.MultiArmedBandit
	kafkaProducer         *kafka.Producer
	slotBannersRepository slotbannersrepository.PgSlotBannerRepository
	statisticRepository   statisticrepository.PgStatisticRepository
	userGroupRepository   usergrouprepository.PgUserGroupRepository
)

func InitKafkaProducer(brokers []string, topic string) {
	kafkaProducer = kafka.NewKafkaProducer(brokers, topic)
}

func InitRepositories() {
	slotBannersRepository = slotbannersrepository.PgSlotBannerRepository{DB: repository.GetDB()}
	statisticRepository = statisticrepository.PgStatisticRepository{DB: repository.GetDB()}
	userGroupRepository = usergrouprepository.PgUserGroupRepository{DB: repository.GetDB()}
}

func InitRotationAlgorithm() {
	banditService = bandit.NewMultiArmedBandit()
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

	if err := slotBannersRepository.AddBannerToSlot(request.SlotID, request.BannerID); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add banner to slot to db"})
		return
	}

	userGroupIDs, err := userGroupRepository.GetAllUserGroupsIDs()
	if err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get userGroupIds from db"})
		return
	}

	if err := statisticRepository.CreateStartStatisticsForBannerInSlot(request.SlotID,
		request.BannerID, userGroupIDs); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError,
			map[string]string{"error": "Failed to CreateStartStatisticsForBannerInSlot"})
		return
	}

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

	if err := slotBannersRepository.RemoveBannerFromSlot(request.SlotID, request.BannerID); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add banner to db"})
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
	if kafkaProducer != nil {
		kafkaProducer.PublishMessage(slotIDBytes, eventBytes)
	}

	if err := statisticRepository.IncrementClick(request.SlotID, request.BannerID, request.UserGroupID); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to record click"})
		return
	}

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
	if kafkaProducer != nil {
		kafkaProducer.PublishMessage(slotIDBytes, eventBytes)
	}

	if err := statisticRepository.IncrementView(request.SlotID, response.BannerID, request.UserGroupID); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to record view"})
		return
	}

	jsonResponse(w, http.StatusOK, response)
}

func idToBytes(id int) []byte {
	slotIDString := strconv.Itoa(id)
	idBytes := []byte(slotIDString)
	return idBytes
}
