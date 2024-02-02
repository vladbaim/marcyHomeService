package file

import (
	"context"
	"marcyHomeService/internal/domain"
	"path/filepath"
	"strings"
)

type fileFirmwareRepository struct {
}

func getFirmwareList() []domain.Firmware {
	return []domain.Firmware{
		{
			Key:          "cabinet",
			Version:      1,
			FirmwareName: "sensor.ino.bin",
		},
	}
}

func NewFileFirmwareRepository() domain.FirmwareRepository {
	return &fileFirmwareRepository{}
}

func (p *fileFirmwareRepository) GetNewFirmwarePath(ctx context.Context, firmwareRequest *domain.GetFirmwareRequest) (response domain.GetFirmwareResponse, err error) {
	var firmwareList = getFirmwareList()
	for _, firmware := range firmwareList {
		if strings.Compare(firmware.Key, firmwareRequest.Key) == 0 && firmware.Version > firmwareRequest.Version {
			response.Path = "static/builds/" + firmware.Key + "/" + firmware.FirmwareName
			_, err = filepath.Abs(response.Path)
			if err != nil {
				return
			}
			return
		}
	}

	err = domain.ErrNotFound

	return
}
