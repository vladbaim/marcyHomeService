package domain

import (
	"context"
)

type Firmware struct {
	Key          string `json:"key" validate:"required"`
	FirmwareName string `json:"firmware_name" validate:"required"`
	Version      int    `json:"version" validate:"required"`
}

type GetFirmwareRequest struct {
	Key     string `json:"key" validate:"required"`
	Version int    `json:"version" validate:"required"`
}

type GetFirmwareResponse struct {
	Path string `json:"path" validate:"required"`
}

type FirmwareUsecase interface {
	GetNewFirmwarePath(context.Context, *GetFirmwareRequest) (GetFirmwareResponse, error)
}

type FirmwareRepository interface {
	GetNewFirmwarePath(context.Context, *GetFirmwareRequest) (GetFirmwareResponse, error)
}
