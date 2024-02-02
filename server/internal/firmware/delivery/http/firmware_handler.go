package http

import (
	"encoding/json"
	"log"
	"marcyHomeService/internal/domain"
	"marcyHomeService/pkg/common"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	getFirmware = "/api/firmware"
)

type firmwareHandler struct {
	firmwareUsecase domain.FirmwareUsecase
}

func NewFirmwareHandlerHandler(firmwareUsecase domain.FirmwareUsecase, router *httprouter.Router) {
	handler := &firmwareHandler{
		firmwareUsecase: firmwareUsecase,
	}

	router.HandlerFunc(http.MethodGet, getFirmware, handler.getFirmware)
}

func (handler *firmwareHandler) getFirmware(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var firmwareRequest domain.GetFirmwareRequest
	queryValues := r.URL.Query()

	var err error
	firmwareRequest.Key = queryValues.Get("key")
	firmwareRequest.Version, err = strconv.Atoi(queryValues.Get("version"))
	if err != nil {
		log.Printf("FirmwareRequest data error: %+v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	log.Printf("firmwareRequest: %+v", firmwareRequest)

	var ok bool
	if ok, err = common.IsRequestValid(&firmwareRequest); !ok {
		log.Printf("FirmwareRequest data error: %+v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response, err := handler.firmwareUsecase.GetNewFirmwarePath(r.Context(), &firmwareRequest)
	if err != nil {
		log.Printf("FirmwareRequest data error: %+v", err)
		http.Error(w, err.Error(), common.GetStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
