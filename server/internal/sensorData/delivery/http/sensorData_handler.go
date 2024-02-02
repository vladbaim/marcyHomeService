package http

import (
	"encoding/json"
	"log"
	"marcyHomeService/internal/domain"
	"marcyHomeService/pkg/common"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	getLastSensorData = "/api/sensor/last"
	sensorData        = "/api/sensor"
)

type SensorDataHandler struct {
	SUsecase domain.SensorDataUsecase
}

// NewSensorDataHandler will initialize the articles/ resources endpoint
func NewSensorDataHandler(us domain.SensorDataUsecase, router *httprouter.Router) {
	handler := &SensorDataHandler{
		SUsecase: us,
	}

	router.HandlerFunc(http.MethodGet, getLastSensorData, handler.GetLast)
	router.HandlerFunc(http.MethodPost, sensorData, handler.Store)
}

func (s *SensorDataHandler) GetLast(w http.ResponseWriter, r *http.Request) {
	sensorData, err := s.SUsecase.GetLast(r.Context())
	if err != nil {
		log.Printf("Get data error: %+v", err)
		http.Error(w, err.Error(), common.GetStatusCode(err))
		return
	}

	log.Printf("GetLastSensorData: %+v", sensorData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensorData)
}

func (s *SensorDataHandler) Store(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var sensorData domain.SensorData

	err := json.NewDecoder(r.Body).Decode(&sensorData)
	if err != nil {
		log.Printf("Set sensor data error: %+v", err)
		http.Error(w, err.Error(), common.GetStatusCode(err))
		return
	}

	log.Printf("StoreSensorData: %+v", sensorData)

	var ok bool
	if ok, err = common.IsRequestValid(&sensorData); !ok {
		log.Printf("Insert data error: %+v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = s.SUsecase.Store(r.Context(), &sensorData)
	if err != nil {
		log.Printf("Insert data error: %+v", err)
		http.Error(w, err.Error(), common.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
