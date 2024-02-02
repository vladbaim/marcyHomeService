package domain

import (
	"context"
	"time"
)

type SensorData struct {
	ID            int       `json:"id"`
	Version       int       `json:"version"`
	Position      string    `json:"position" validate:"required"`
	Humidity      float32   `json:"humidity" validate:"required"`
	Temperature   float32   `json:"temperature" validate:"required"`
	CarbonDioxide int       `json:"carbon_dioxide" validate:"required"`
	CreatedAt     time.Time `json:"created_at"`
}

type SensorDataUsecase interface {
	GetLast(ctx context.Context) (SensorData, error)
	Store(context.Context, *SensorData) error
}

type SensorDataRepository interface {
	GetLast(ctx context.Context) (SensorData, error)
	Store(ctx context.Context, a *SensorData) error
}
