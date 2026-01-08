package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/truegul/api-server/internal/migrations"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	m, err := createMigrator(databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid number of steps: %v", err)
			}
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Printf("Migration down %d step(s) completed successfully", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("No migrations applied yet")
				os.Exit(0)
			}
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current version: %d (dirty: %v)", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force command requires a version number")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		log.Printf("Forced version to %d", version)

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
		log.Println("All tables dropped successfully")

	default:
		log.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func createMigrator(databaseURL string) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(migrations.FS, migrations.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}

func printUsage() {
	fmt.Println("Usage: migrate <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up              Apply all pending migrations")
	fmt.Println("  down [N]        Rollback N migrations (default: 1)")
	fmt.Println("  version         Show current migration version")
	fmt.Println("  force VERSION   Force set migration version (use with caution)")
	fmt.Println("  drop            Drop all tables in database")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  DATABASE_URL    PostgreSQL connection string (required)")
}
