package mng

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (m *Management) DownMigrations() error {
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		filename := filepath.Base(file)
		idVersion := strings.Split(filename, "_")[0]
		approve, err := m.CheckMigrateVersion("down", idVersion)
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
		// Проверить данные миграции
		err = verifyDataInMigration([]byte(dataMigrate.Up))
		if err != nil {
			continue
		}
		// Выполнить миграцию
		if err := executeMigration(m.Cfg.Ctx, m.Cfg.DB, file, dataMigrate.Down); err != nil {
			return fmt.Errorf("migration failed for %s: %w", file, err)
		}
		// зафиксировать версию
		err = m.CommitMigrateVersion(false, idVersion)
		if err != nil {
			return fmt.Errorf("failed commit migrate version: %w", err)
		}
	}

	return nil
}
