package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"marcyHomeService/internal/config"
	"marcyHomeService/pkg/client/postgresql"
	"marcyHomeService/pkg/client/telegram"
	"marcyHomeService/pkg/metric"
	"net"
	"net/http"
	"time"

	_ "marcyHomeService/docs"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"

	_firmwareDelivery "marcyHomeService/internal/firmware/delivery/http"
	_firmwareRepo "marcyHomeService/internal/firmware/repository/file"
	_firmwareUcase "marcyHomeService/internal/firmware/usecase"

	_sensorDataDelivery "marcyHomeService/internal/sensorData/delivery/http"
	_sensorDataRepo "marcyHomeService/internal/sensorData/repository/postgresql"
	_sensorDataUcase "marcyHomeService/internal/sensorData/usecase"

	_telegramBotDelivery "marcyHomeService/internal/telegramBot/delivery/http"
	_telegramBotUcase "marcyHomeService/internal/telegramBot/usecase"
)

type App struct {
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
	pgClient   *pgxpool.Pool
}

func NewApp(cfg *config.Config) (App, error) {
	log.Print("router init")

	router := httprouter.New()

	Migrate(cfg)

	log.Print("postgres connect")
	pgConfig := postgresql.NewPgConfig(
		cfg.PostgresSQL.Username,
		cfg.PostgresSQL.Password,
		cfg.PostgresSQL.Host,
		cfg.PostgresSQL.Port,
		cfg.PostgresSQL.Database,
	)
	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("swagger docs init")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusPermanentRedirect))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	log.Print("heartbeat metric init")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	firmwareRepo := _firmwareRepo.NewFileFirmwareRepository()
	sensorDataRepo := _sensorDataRepo.NewPostgresqlSensorDataRepository(pgClient)

	sensorDataUseCase := _sensorDataUcase.NewSensorDataUsecase(sensorDataRepo)
	_sensorDataDelivery.NewSensorDataHandler(sensorDataUseCase, router)

	telegramBot, whUrl, err := telegram.NewTelegramBotClient(cfg.App.Url, cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}
	telegramBotUsecase := _telegramBotUcase.NewTelegramBotUsecase(telegramBot, sensorDataRepo)
	_telegramBotDelivery.NewTelegramBotHandler(telegramBotUsecase, router, whUrl)

	_firmwareUcase := _firmwareUcase.NewFirmwareUsecase(firmwareRepo)
	_firmwareDelivery.NewFirmwareHandlerHandler(_firmwareUcase, router)

	return App{
		cfg:      cfg,
		router:   router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHttp(ctx)
	})
	log.Print("app init and started")
	return grp.Wait()
}

func (a *App) startHttp(ctx context.Context) error {
	log.Print("start HTTP")
	log.Printf("bind application to host: %s and port %s", a.cfg.HTTP.IP, a.cfg.HTTP.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		log.Fatal(err)
	}

	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler: handler,
	}

	log.Print("application init and started!")

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Print("server shutdown")
		default:
			log.Fatal(err)
		}
	}

	if err = a.httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	return err
}
