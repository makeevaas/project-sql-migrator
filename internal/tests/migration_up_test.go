package tests

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	mng "github.com/makeevaas/project/sql-migrator/internal"
	"github.com/makeevaas/project/sql-migrator/pkg/cfg"
	dbPkg "github.com/makeevaas/project/sql-migrator/pkg/db"
)

// TestAdd проверяет функцию Add.
func TestUpMigrations(t *testing.T) {
	migrationText := `# The SQL in the UP section is executed when applying a migration.
# The SQL in the DOWN section is executed when rolling back a migration.
up: CREATE TABLE IF NOT EXISTS users12 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL);
down: DROP TABLE IF EXISTS users12;
`
	f := "../migrations/20250907202020_test.yaml"
	migrateNum := "20250907202020"

	file, err := os.Create(f)
	if err != nil {
		t.Fatalf("failed to create migration file: %v", err)
	}
	defer file.Close()
	t.Log("Создан временный файл - ", f)

	_, err = os.Stat(f)
	if os.IsNotExist(err) {
		t.Fatalf("Файл %s не найден", file.Name())
	}

	_, err = os.ReadFile(f)
	if err != nil {
		t.Fatalf("Не удалось прочитать файл: %v", err)
	}

	// Записываем данные в файл
	_, err = file.WriteString(migrationText)
	if err != nil {
		t.Fatalf("error write comment string: %v", err)
	}

	dbConnPath := "postgres://postgres:pwd@localhost:5432/main_db?sslmode=disable"
	ctx := context.Background()
	db, err := dbPkg.Conn(ctx, dbConnPath)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	filesPath := "../migrations"
	filePath := filepath.Join(filesPath, "20250907202020_test.yaml")
	migrateFiles := []string{filePath}
	m := &mng.Management{
		Cfg: cfg.Config{Ctx: context.Background(), DB: db, MigrationFiles: migrateFiles, MigratePath: "../migrations"},
	}

	err = m.UpMigrations()
	if err != nil {
		t.Fatalf("Не удалось выполнить миграцию: %v", err)
	}

	rows, err := m.Cfg.DB.Conn.Query(context.Background(), mng.GetMigrateDataReq, migrateNum)
	if err != nil {
		t.Fatalf("failed to query data: %v", err)
	}
	defer rows.Close()

	var versionID string
	var isApplied bool
	var tstamp time.Time
	for rows.Next() {
		if err := rows.Scan(&versionID, &isApplied, &tstamp); err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("failed to scan row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("error during rows iteration: %v", err)
	}

	t.Log(isApplied, versionID, tstamp)

	if versionID != migrateNum || !isApplied {
		output := `Содержимое файла не совпадает. Ожидалось:
		version_id: 20250907202020
		is_applied: true
		получено: \nversion_id: %s
		is_applied: %v`
		t.Fatalf(output, versionID, isApplied)
	}

	// откатить или удалить данные миграции и удалить временный файл
	deleteMigrateReq := `DELETE from db_version 
	where version_id=$1;`
	_, err = m.Cfg.DB.Conn.Exec(context.Background(), deleteMigrateReq, migrateNum)
	if err != nil {
		t.Fatalf("failed to delete test migration: %v", err)
	}
	deleteMigrateReq = `DROP TABLE IF EXISTS users12;`
	_, err = m.Cfg.DB.Conn.Exec(context.Background(), deleteMigrateReq)
	if err != nil {
		t.Fatalf("failed to delete test migration: %v", err)
	}
	t.Log("Тестовая миграция удалена - ", migrateNum)
	os.Remove(file.Name())
	t.Log("Временный файл удален - ", f)
}
