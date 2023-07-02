package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/chema0/url-shortener/config"
)

var (
	Global      *sql.DB
	connectOnce sync.Once
)

func TestContext(t testing.TB) context.Context {
	os.Setenv("env", config.Test)

	connectOnce.Do(func() {
		// Loads configuration and configures logger
		cfg := config.LoadConfig()

		db, err := Open(BuildConnString(&cfg.Database, true))
		if err != nil {
			t.Fatal("failed to connect to db server")
		}

		// Creates test database, drops the previous one before
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", cfg.Database.Name))
		if err != nil {
			t.Errorf("failed to drop testing db: %s", err)
		}
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", cfg.Database.Name))
		if err != nil {
			t.Errorf("failed to create testing db: %s", err)
		}

		Global, err = NewDatabase(&cfg.Database)
		if err != nil {
			t.Fatal("failed to load config")
		}
	})

	emptyDBPreserveSchema(t, Global)

	return context.TODO()
}

func emptyDBPreserveSchema(t testing.TB, d *sql.DB) {
	_, err := d.Exec(`SELECT * FROM schema_migrations`)
	if err != nil {
		t.Fatalf("Table schema_migrations not found: %v", err)
	}

	rows, err := d.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'schema_migrations'")
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			t.Fatal(err)
		}
		tables = append(tables, table)
	}
	if err = rows.Close(); err != nil {
		t.Fatal(err)
	}
	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}
	if testing.Verbose() {
		t.Logf("Truncating all %d tables", len(tables))
	}
	_, err = d.Exec("TRUNCATE " + strings.Join(tables, ", ") + " RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}
}
