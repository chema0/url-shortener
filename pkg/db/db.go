package db

import (
	"context"
	"fmt"
	"os"

	"github.com/chema0/url-shortener/config"

	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	Conn       *pgx.Conn
	ConnString string
}

func NewDatabase(config *config.DatabaseConfig) (*Database, error) {
	// To use a connection pool replace the import github.com/jackc/pgx/v5 with github.com/jackc/pgx/v5/pgxpool
	// and connect with pgxpool.New() instead of pgx.Connect()
	connString := buildConnString(config)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	return &Database{
		Conn:       conn,
		ConnString: connString,
	}, nil
}

func (db *Database) RunMigrations() {
	fmt.Println("Running UP database migrations..")
	m, err := migrate.New("file://migrations/", db.ConnString)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}
	if err := m.Up(); err != nil {
		fmt.Println("No changes after running migrations")
	}
}

func (db *Database) RunMigrationsDown() {
	fmt.Println("Running DOWN database migrations..")
	m, err := migrate.New("file://migrations/", db.ConnString)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}
	if err := m.Down(); err != nil {
		fmt.Println("No changes after running down migrations")
	}
}

func buildConnString(config *config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.Name)
}
