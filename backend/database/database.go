package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens a PostgreSQL connection using the pgx driver,
// validates it with a ping, and sets sensible pool limits.
func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	return db, nil
}

// RunMigrations reads SQL files from the embedded filesystem and executes
// any that haven't been applied yet, tracking versions in a schema_migrations table.
func RunMigrations(db *sql.DB, migrationsFS fs.FS) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, "db/migrations")
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, err := parseVersion(entry.Name())
		if err != nil {
			return fmt.Errorf("parsing migration version from %s: %w", entry.Name(), err)
		}

		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking migration %d: %w", version, err)
		}
		if exists {
			continue
		}

		content, err := fs.ReadFile(migrationsFS, "db/migrations/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", entry.Name(), err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("beginning transaction for migration %d: %w", version, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", entry.Name(), err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording migration %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %d: %w", version, err)
		}
	}

	return nil
}

// RunSeeds reads SQL files from the embedded filesystem and executes
// any that haven't been applied yet, tracking versions in a schema_seeds table.
func RunSeeds(db *sql.DB, seedsFS fs.FS) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_seeds (
			version    INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_seeds table: %w", err)
	}

	entries, err := fs.ReadDir(seedsFS, "db/seeds")
	if err != nil {
		return fmt.Errorf("reading seeds directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, err := parseVersion(entry.Name())
		if err != nil {
			return fmt.Errorf("parsing seed version from %s: %w", entry.Name(), err)
		}

		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_seeds WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking seed %d: %w", version, err)
		}
		if exists {
			continue
		}

		content, err := fs.ReadFile(seedsFS, "db/seeds/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading seed %s: %w", entry.Name(), err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("beginning transaction for seed %d: %w", version, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing seed %s: %w", entry.Name(), err)
		}

		if _, err := tx.Exec("INSERT INTO schema_seeds (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording seed %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing seed %d: %w", version, err)
		}
	}

	return nil
}

// parseVersion extracts the numeric prefix from a migration filename like "001_create_zones.sql".
func parseVersion(filename string) (int, error) {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid migration filename: %s", filename)
	}
	return strconv.Atoi(parts[0])
}
