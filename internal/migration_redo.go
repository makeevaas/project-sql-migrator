package mng

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func (m *Management) RedoMigrations() error {
	// получить последний файл миграции
	file := m.Cfg.MigrationFiles[len(m.Cfg.MigrationFiles)-1]
	// проверить версию миграции
	idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
	approve, err := m.CheckMigrateVersion("DOWN", idVersion)
	if err != nil {
		log.Fatalf("failed check migrate version: %v", err)
	}
	if !approve {
		log.Info("not migrate\n")
		return nil
	}
	// Получить данные миграции
	dataMigrate, err := readFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// Выполнить откат
	if err := executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Down); err != nil {
		log.Fatalf("migration failed for %s: %v", file, err)
	}
	// зафиксировать версию
	err = m.CommitMigrateVersion(false, idVersion)
	if err != nil {
		log.Fatalf("failed commit migrate version: %v", err)
	}
	// Выполнить накат
	if err = executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Up); err != nil {
		log.Fatalf("migration failed for %s: %v", file, err)
	}
	// зафиксировать версию
	err = m.CommitMigrateVersion(true, idVersion)
	if err != nil {
		log.Fatalf("failed commit migrate version: %v", err)
	}

	return nil
}
