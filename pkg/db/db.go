package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chema0/url-shortener/config"
	"github.com/chema0/url-shortener/pkg/utils"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	Conn       *pgx.Conn
	ConnString string
}

// DefaultMaxOpenConnections the default value for max open connections in the
// PostgreSQL connection pool
const DefaultMaxOpenConnections = 10

func NewDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	// Force PostfreSQL session timezone to UTC
	if v, ok := os.LookupEnv("PGTZ"); ok && strings.ToLower(v) != "utc" {
		log.Warn().Str("ingoredPGTZ", v).Msg("Ingoring PGTZ environment variable; using PGTZ=UTC.")
	}
	if err := os.Setenv("PGTZ", "UTC"); err != nil {
		return nil, errors.New("error setting PGTZ=UTC")
	}

	db, err := openDBWithStartupWait(buildConnString(cfg, false))
	if err != nil {
		println(err.Error())
		return nil, errors.New("DB not available")
	}
	configureConnectionPool(db)

	log.Info().Msg("DoMigrate")

	if err := DoMigrate(NewMigrate(db)); err != nil {
		return nil, errors.New("failed to migrate the DB")
	}

	log.Info().Msg("DoMigrate complete")

	return db, nil
}

var startupTimeout = func() time.Duration {
	str := utils.Get("DB_STARTUP_TIMEOUT", "10s")
	d, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal().Err(err).Msg("db startup timeout")
	}
	return d
}()

func openDBWithStartupWait(connString string) (db *sql.DB, err error) {
	// Allow the DB to take up to 10s while it reports "pq: the database system is starting up".
	startupDeadline := time.Now().Add(startupTimeout)
	for {
		if time.Now().After(startupDeadline) {
			return nil, fmt.Errorf("database did not startup within %s (%v)", startupTimeout, err)
		}
		db, err = Open(connString)
		if err == nil {
			err = db.Ping()
		}
		if err != nil && isDatabaseLikelyStartingUp(err) {
			time.Sleep(startupTimeout / 10)
			continue
		}
		return db, err
	}
}

// isDatabaseLikelyStartingUp returns whether the err likely just means the PostgreSQL database is
// starting up, and it should not be treated as a fatal error during program initialization.
func isDatabaseLikelyStartingUp(err error) bool {
	if strings.Contains(err.Error(), "pq: the database system is starting up") {
		// Wait for DB to start up.
		return true
	}
	if e, ok := err.(net.Error); ok && strings.Contains(e.Error(), "connection refused") {
		// Wait for DB to start listening.
		return true
	}
	return false
}

// Open creates a new DB handle with the given schema by connecting to
// the database identified by dataSource (e.g., "dbname=mypgdb" or
// blank to use the PG* env vars).
//
// Open assumes that the database already exists.
func Open(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, errors.New("postgresql open")
	}
	return db, nil
}

func buildConnString(config *config.DatabaseConfig, connectToServer bool) string {
	if connectToServer {
		return fmt.Sprintf("postgres://%s:%s@%s:%v/?sslmode=disable",
			config.Username, config.Password, config.Host, config.Port)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.Name)
}

// Ping attempts to contact the database and returns a non-nil error upon failure. It is intended to
// be used by health checks.
//func Ping(ctx context.Context) error { return Global.PingContext(ctx) }

// configureConnectionPool sets reasonable sizes on the built in DB queue. By
// default the connection pool is unbounded, which leads to the error `pq:
// sorry too many clients already`.
func configureConnectionPool(db *sql.DB) {
	var err error
	maxOpen := DefaultMaxOpenConnections
	if e := os.Getenv("SRC_PGSQL_MAX_OPEN"); e != "" {
		maxOpen, err = strconv.Atoi(e)
		if err != nil {
			log.Fatal().Err(err).Msg("SRC_PGSQL_MAX_OPEN is not an int")
		}
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxOpen)
	db.SetConnMaxLifetime(time.Minute)
}
