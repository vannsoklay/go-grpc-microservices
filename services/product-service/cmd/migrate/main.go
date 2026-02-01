package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var (
		dbURL          = flag.String("db", "", "Database URL")
		migrationsPath = flag.String("path", "db/migrations", "Path to migrations directory")
		command        = flag.String("cmd", "up", "Migration command: up, down, version, force")
		version        = flag.Int("version", -1, "Migration version for force command")
	)
	flag.Parse()

	if *dbURL == "" {
		*dbURL = os.Getenv("DATABASE_URL")
	}

	if *dbURL == "" {
		log.Fatal("Database URL is required")
	}

	// Get absolute path to migrations
	absPath, err := filepath.Abs(*migrationsPath)
	if err != nil {
		log.Fatalf("Could not get absolute path: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", absPath)
	}

	log.Printf("Using migrations from: %s", absPath)

	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
	log.Println("Database connection successful")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create database driver: %v", err)
	}

	// Use absolute path for file URL
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	switch *command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Could not run up migrations: %v", err)
		}
		log.Println("✅ Migrations applied successfully!")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Could not run down migrations: %v", err)
		}
		log.Println("✅ Migrations rolled back successfully!")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Could not get version: %v", err)
		}
		log.Printf("Version: %d, Dirty: %t\n", v, dirty)
	case "force":
		if *version < 0 {
			log.Fatal("Version is required for force command")
		}
		if err := m.Force(*version); err != nil {
			log.Fatalf("Could not force version: %v", err)
		}
		log.Printf("✅ Forced version to %d\n", *version)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
