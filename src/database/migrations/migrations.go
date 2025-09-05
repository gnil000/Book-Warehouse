package migrations

import (
	"embed"
	"fmt"

	"gin_main/config"
	"gin_main/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	glog "go.finelli.dev/gooseloggers/zerolog"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(cfg *config.Config, logger zerolog.Logger) {
	err := migrate(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate")
		//panic(fmt.Errorf("failed to migrate: %w", err))
	}
}

func migrate(cfg *config.Config, log zerolog.Logger) error {
	db, err := goose.OpenDBWithDriver("pgx", cfg.Database.BookDB)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	goose.SetTableName("book_warehouse_api_migrations")
	goose.SetBaseFS(embedMigrations)

	log = logger.WithModule(log, "migrations")
	goose.SetLogger(glog.GooseZerologLogger(&log))

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to make migrations: %w", err)
	}

	return nil
}
