package db

import (
	"database/sql"
	"os"
	"time"

	"github.com/chema0/url-shortener/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

func NewMigrate(db *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load driver")
	}

	d, err := iofs.New(migrations.MigrationsFs, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read assets from iofs")
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migration from go-bindata")
	}
	// In case another process was faster and runs migrations, we will wait
	// this long
	m.LockTimeout = 5 * time.Minute

	return m
}

func DoMigrate(m *migrate.Migrate) (err error) {
	err = m.Up()
	if err == nil || err == migrate.ErrNoChange {
		return nil
	}

	if os.IsNotExist(err) {
		// This should only happen if the DB is ahead of the migrations available
		version, _, verr := m.Version()
		if verr != nil {
			return verr
		}
		log.Warn().Uint("db_version", version).Msg("WARNING: Detected an old version of database.")
		return nil
	}
	return err
}

// func (db *Database) RunMigrations() {
// 	fmt.Println("Running UP database migrations..")
// 	m, err := migrate.New("file://migrations/", db.ConnString)
// 	if err != nil {
// 		panic(fmt.Errorf("unable to connect to database: %v", err))
// 	}
// 	if err := m.Up(); err != nil {
// 		fmt.Println("No changes after running migrations")
// 	}
// }

// func (db *Database) RunMigrationsDown() {
// 	fmt.Println("Running DOWN database migrations..")
// 	m, err := migrate.New("file://migrations/", db.ConnString)
// 	if err != nil {
// 		panic(fmt.Errorf("unable to connect to database: %v", err))
// 	}
// 	if err := m.Down(); err != nil {
// 		fmt.Println("No changes after running down migrations")
// 	}
// }
