package mng

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func (m *Management) UpMigrations() error {
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
		approve, err := m.CheckMigrateVersion("UP", idVersion)
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
			return err
		}
		// Выполнить миграцию
		if err := executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Up); err != nil {
			log.Fatalf("migration failed for %s: %v", file, err)
		}
		// зафиксировать версию
		err = m.CommitMigrateVersion(true, idVersion)
		if err != nil {
			log.Fatalf("failed commit migrate version: %v", err)
		}
	}
	return nil
}
