package api

import (
	"encoding/json"
	"net/http"

	"github.com/yuriiwanchev/banner-rotation-service/internal/logic/bandit"
	m "github.com/yuriiwanchev/banner-rotation-service/internal/models"
)

var banditService = bandit.NewMultiArmedBandit()

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
	// jsonResponse(w, http.StatusOK, nil)
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

	banditService.RemoveBanner(request.SlotID, request.BannerID)
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

	banditService.RecordClick(request.SlotID, request.BannerID, request.UserGroupID)
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

	// jsonResponse(w, http.StatusOK, map[string]string{"banner_id": bannerID})
	jsonResponse(w, http.StatusOK, response)
}
