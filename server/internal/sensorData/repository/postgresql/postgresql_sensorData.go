package postgresql

import (
	"context"
	"log"
	"marcyHomeService/internal/domain"
	"marcyHomeService/pkg/client/postgresql"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresqlSensorDataRepository struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
}

func NewPostgresqlSensorDataRepository(client *pgxpool.Pool) domain.SensorDataRepository {
	return &postgresqlSensorDataRepository{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       client,
	}
}

const (
	scheme = "public"
	table  = "sensor_data"
)

func (p *postgresqlSensorDataRepository) GetLast(ctx context.Context) (sensorData domain.SensorData, err error) {
	query := p.queryBuilder.Select(
		"id",
		"version",
		"position",
		"temperature",
		"humidity",
		"carbon_dioxide",
		"created_at",
	).
		OrderBy("id DESC").
		Limit(1).
		From(scheme + "." + table)

	sql, args, err := query.ToSql()
	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		log.Print(err)
		return
	}

	row := p.client.QueryRow(ctx, sql, args...)

	if err = row.Scan(&sensorData.ID,
		&sensorData.Version,
		&sensorData.Position,
		&sensorData.Temperature,
		&sensorData.Humidity,
		&sensorData.CarbonDioxide,
		&sensorData.CreatedAt,
	); err != nil {
		err = postgresql.ErrScan(postgresql.ParsePgError(err))
		log.Print(err)
		return
	}

	return
}

func (p *postgresqlSensorDataRepository) Store(ctx context.Context, sensorData *domain.SensorData) (err error) {
	query := p.queryBuilder.Insert(scheme+"."+table).
		Columns(
			"version",
			"position",
			"temperature",
			"humidity",
			"carbon_dioxide",
		).
		Values(sensorData.Version, sensorData.Position, sensorData.Temperature, sensorData.Humidity, sensorData.CarbonDioxide).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		log.Print(err)
		return
	}

	_, err = p.client.Exec(ctx, sql, args...)

	return
}
