package migrate

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type Migrator struct {
	path string
	dsn  string

	cleanDb []string
	logger  zerolog.Logger
}

func NewMigrator(path, dsn string, options ...OptionFn) *Migrator {
	m := &Migrator{
		path:   fmt.Sprintf("file://%s", strings.TrimPrefix(strings.TrimPrefix(path, "."), "/")),
		dsn:    dsn,
		logger: zerolog.Nop(),
	}

	for _, opt := range options {
		opt(m)
	}
	return m
}

func (m *Migrator) Run(ctx context.Context) error {
	driver, err := database.Open(m.dsn)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to init database driver")
		return err
	}
	defer driver.Close()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.dsn))
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to connect database")
		return err
	}
	defer client.Disconnect(ctx)

	if len(m.cleanDb) > 0 {
		for _, db := range m.cleanDb {
			if err := m.cleanDatabase(ctx, client, db); err != nil {
				return err
			}
		}
	}

	migration, err := migrate.NewWithDatabaseInstance(m.path, "", driver)
	if err != nil {
		return err
	}

	beforeVersion, dirty, err := migration.Version()
	if err != nil && beforeVersion != 0 {
		return err
	}

	m.logger.Info().Uint("version", beforeVersion).Msg("migration started")

	if dirty {
		m.logger.Warn().Msg("previous migration failed")
	}

	err = migration.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == migrate.ErrNoChange {
		m.logger.Info().Msg("no new database changes")
	}

	afterVersion, dirty, err := migration.Version()
	if err != nil && beforeVersion != 0 {
		return err
	}

	m.logger.Info().Uint("version", afterVersion).Msg("migration done")

	if dirty {
		m.logger.Warn().Msg("previous migration failed")
	}

	return nil
}

// Clean database
func (m *Migrator) cleanDatabase(ctx context.Context, client *mongo.Client, dbName string) error {
	m.logger.Info().Msgf("clean database %s\n", dbName)
	if err := client.Database(dbName).Drop(ctx); err != nil {
		return err
	}
	return nil
}
