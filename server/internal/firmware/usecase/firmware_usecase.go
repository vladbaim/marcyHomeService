package usecase

import (
	"context"
	"marcyHomeService/internal/domain"
)

type firmwareUsecase struct {
	firmwareRepo domain.FirmwareRepository
}

func NewFirmwareUsecase(firmwareRepo domain.FirmwareRepository) domain.FirmwareUsecase {
	return &firmwareUsecase{
		firmwareRepo: firmwareRepo,
	}
}

func (f *firmwareUsecase) GetNewFirmwarePath(c context.Context, firmwareRequest *domain.GetFirmwareRequest) (response domain.GetFirmwareResponse, err error) {
	return f.firmwareRepo.GetNewFirmwarePath(c, firmwareRequest)
}
