package common

import (
	"log"
	"marcyHomeService/internal/domain"
	"net/http"
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	log.Printf("Data http handler error: %+v", err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
