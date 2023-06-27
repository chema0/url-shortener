package db

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
)

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
