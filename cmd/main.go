package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	mng "github.com/makeevaas/project/sql-migrator/internal"
	"github.com/makeevaas/project/sql-migrator/pkg/cfg"
	"github.com/makeevaas/project/sql-migrator/pkg/db"
	l "github.com/sirupsen/logrus"
)

var up, down, create, redo, status bool

func init() {
	flag.BoolVar(&up, "up", false, "up migrations")
	flag.BoolVar(&down, "down", false, "down migrations")
	flag.BoolVar(&create, "create", false, "create file migration")
	flag.BoolVar(&redo, "redo", false, "repeat the last migration")
	flag.BoolVar(&status, "status", false, "return migrations status")
}

func main() {
	// Database connection parameters
	flag.Parse()
	ctx := context.Background()
	// "postgres://postgres:pwd@localhost:5432/main_db?sslmode=disable"
	dbConnPath := os.Getenv("DB_CONNECTION_PATH")
	if dbConnPath == "" {
		l.Fatal("DB_CONNECTION_PATH environment variable not set")
	}
	migratePath := os.Getenv("MIGRATIONS_PATH")
	// ./migrations/
	if migratePath == "" {
		l.Fatal("MIGRATIONS_PATH environment variable not set")
	}
	// Подключение к базе данных
	db, err := db.Conn(ctx, dbConnPath)
	if err != nil {
		l.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Conn.Close(ctx)
	// файлы миграции в директории
	files, err := os.ReadDir(migratePath)
	if err != nil {
		l.Fatalf("error open directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			filePath := filepath.Join(migratePath, file.Name())
			migrationFiles = append(migrationFiles, filePath)
		}
	}

	config := &cfg.Config{
		Ctx:            ctx,
		DB:             db,
		MigrationFiles: migrationFiles,
		MigratePath:    migratePath,
	}

	management := &mng.Management{
		Cfg: *config,
	}

	// Выполнение команд утилиты
	if create {
		if _, err := management.CreateFileMigration(); err != nil {
			l.Fatalf("failed to run create migrations: %v", err)
		}
	}
	if up {
		if err := management.UpMigrations(); err != nil {
			l.Fatalf("failed to run up migrations: %v", err)
		}
	}

	if down {
		if err := management.DownMigrations(); err != nil {
			l.Fatalf("failed to run down migrations: %v", err)
		}
	}
	if redo {
		if err := management.RedoMigrations(); err != nil {
			l.Fatalf("failed to run redo migrations: %v", err)
		}
	}
	if status {
		statusMigrate, err := management.StatusMigrations()
		if err != nil {
			l.Fatalf("failed to run status migrations: %v", err)
		}
		for _, r := range statusMigrate {
			l.Info(r)
		}
	}
}
