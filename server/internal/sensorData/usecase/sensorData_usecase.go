package usecase

import (
	"context"
	"marcyHomeService/internal/domain"
)

type sensorDataUsecase struct {
	sensorDataRepo domain.SensorDataRepository
}

func NewSensorDataUsecase(sensorDataRepo domain.SensorDataRepository) domain.SensorDataUsecase {
	return &sensorDataUsecase{
		sensorDataRepo: sensorDataRepo,
	}
}

func (s *sensorDataUsecase) GetLast(c context.Context) (sensorData domain.SensorData, err error) {
	sensorData, err = s.sensorDataRepo.GetLast(c)
	if err != nil {
		return
	}
	return
}

func (s *sensorDataUsecase) Store(ctx context.Context, sensorData *domain.SensorData) (err error) {
	err = s.sensorDataRepo.Store(ctx, sensorData)
	return
}
