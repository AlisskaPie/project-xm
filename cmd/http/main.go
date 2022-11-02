package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	delivery "github.com/AlisskaPie/project-xm/internal/company/delivery/http"
	"github.com/AlisskaPie/project-xm/internal/company/event_sender/noop"
	"github.com/AlisskaPie/project-xm/internal/company/repository/postgres"
	"github.com/AlisskaPie/project-xm/internal/company/usecase"
	"github.com/AlisskaPie/project-xm/internal/config/viper"
	"github.com/AlisskaPie/project-xm/internal/user/delivery/http/middleware"
)

func main() {
	ctx := context.Background()

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	conf, err := viper.GetConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get config: %w", err))
	}

	dbConn, err := sqlx.Open("postgres", conf.DB.DSN)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to sql.Open: %w", err))
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal(fmt.Errorf("failed to Ping db: %w", err))
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(fmt.Errorf("failed to close db:  %w", err))
		}
	}()

	addMigrations(conf)

	e := echo.New()
	e.Use(emiddleware.Logger())

	auth := middleware.KeyAuth([]byte(conf.Auth.JWTKey))

	companyRepo := postgres.NewCompanyRepository(ctx, dbConn)
	if conf.EventSender {
		companyRepo = postgres.NewEventSenderWrapper(
			companyRepo,
			noop.NewCompanyEventSenderNoop(logger),
		)
	}

	companyUsecase := usecase.NewCompanyUsecase(companyRepo)
	delivery.NewCompanyHandler(e, companyUsecase, auth, logger)

	e.Logger.Fatal(e.Start(conf.HTTP.ListenHostPort))
}
