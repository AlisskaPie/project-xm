package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/AlisskaPie/project-xm/internal/config"

	"github.com/golang-migrate/migrate/v4"
)

func addMigrations(conf config.Config) {
	m, err := migrate.New("file://migrations", conf.DB.DSN)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create migration handler: %w", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(fmt.Errorf("failed to apply migration: %w", err))
	}
}
