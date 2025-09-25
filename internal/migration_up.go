package mng

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (m *Management) UpMigrations() error {
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		filename := filepath.Base(file)
		idVersion := strings.Split(filename, "_")[0]
		approve, err := m.CheckMigrateVersion("up", idVersion)
		if err != nil {
			return fmt.Errorf("failed check migrate version: %w", err)
		}
		if !approve {
			continue
		}
		// Получить данные миграции
		dataMigrate, err := readFile(file)
		if err != nil {
			return fmt.Errorf("failed to get migration file data: %w", err)
		}
		// Выполнить миграцию
		if err := executeMigration(m.Cfg.Ctx, m.Cfg.DB, file, dataMigrate.Up); err != nil {
			return fmt.Errorf("migration failed for %s: %w", file, err)
		}
		// зафиксировать версию
		err = m.CommitMigrateVersion(true, idVersion)
		if err != nil {
			return fmt.Errorf("failed commit migrate version: %w", err)
		}
	}
	return nil
}
