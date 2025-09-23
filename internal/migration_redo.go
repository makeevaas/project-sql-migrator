package mng

import (
	"fmt"
	"strings"
)

func (m *Management) RedoMigrations() error {
	// получить последний файл миграции
	file := m.Cfg.MigrationFiles[len(m.Cfg.MigrationFiles)-1]
	// проверить версию миграции
	idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
	approve, err := m.CheckMigrateVersion("DOWN", idVersion)
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
	// Выполнить откат
	if err := executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Down); err != nil {
		return fmt.Errorf("migration failed for %s: %w", file, err)
	}
	// зафиксировать версию
	err = m.CommitMigrateVersion(false, idVersion)
	if err != nil {
		return fmt.Errorf("failed commit migrate version: %w", err)
	}
	// Выполнить накат
	if err = executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Up); err != nil {
		return fmt.Errorf("migration failed for %s: %w", file, err)
	}
	// зафиксировать версию
	err = m.CommitMigrateVersion(true, idVersion)
	if err != nil {
		return fmt.Errorf("failed commit migrate version: %w", err)
	}

	return nil
}
