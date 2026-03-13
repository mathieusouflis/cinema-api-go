package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|create> [name] [steps]")
	}

	switch os.Args[1] {
	case "up":
		runUp()
	case "down":
		steps := 1
		if len(os.Args) >= 3 {
			var err error
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("invalid steps: %v", err)
			}
		}
		runDown(steps)
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate create <name>")
		}
		runCreate(os.Args[2])
	default:
		log.Fatalf("unknown command %q — use up, down, or create", os.Args[1])
	}
}

func newMigrate() *migrate.Migrate {
	dbURL := mustEnv("DATABASE_URL")
	migrationsPath := mustEnv("MIGRATIONS_PATH")

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}
	return m
}

func runUp() {
	m := newMigrate()
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate up: %v", err)
	}

	v, _, _ := m.Version()
	fmt.Printf("ok  migrate up  (version %d)\n", v)
}

func runDown(steps int) {
	m := newMigrate()
	defer m.Close()

	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate down: %v", err)
	}

	v, _, _ := m.Version()
	fmt.Printf("ok  migrate down %d step(s)  (version %d)\n", steps, v)
}

func runCreate(name string) {
	migrationsPath := mustEnv("MIGRATIONS_PATH")

	version, err := nextVersion(migrationsPath)
	if err != nil {
		log.Fatalf("nextVersion: %v", err)
	}

	base := filepath.Join(migrationsPath, fmt.Sprintf("%06d_%s", version, name))
	upFile := base + ".up.sql"
	downFile := base + ".down.sql"

	for _, f := range []string{upFile, downFile} {
		if err := os.WriteFile(f, []byte(""), 0o644); err != nil {
			log.Fatalf("create %s: %v", f, err)
		}
		fmt.Println("created", f)
	}
}

// nextVersion returns the next sequential migration version number.
func nextVersion(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}

	max := 0
	for _, e := range entries {
		parts := strings.SplitN(e.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}
		v, err := strconv.Atoi(parts[0])
		if err == nil && v > max {
			max = v
		}
	}
	return max + 1, nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env %s is not set", key)
	}
	return v
}
