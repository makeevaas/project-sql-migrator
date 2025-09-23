package mng

import (
	"fmt"
	"strings"
)

func (m *Management) UpMigrations() error {
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
		approve, err := m.CheckMigrateVersion("UP", idVersion)
		if err != nil {
			return fmt.Errorf("failed check migrate version: %w", err)
		}
		if !approve {
			fmt.Printf("not migrate\n")
			return nil
		}
		// Получить данные миграции
		dataMigrate, err := readFile(file)
		if err != nil {
			return err
		}
		// Выполнить миграцию
		if err := executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Up); err != nil {
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
