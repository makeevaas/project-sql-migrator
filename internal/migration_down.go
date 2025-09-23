package mng

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (m *Management) DownMigrations() error {
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
		approve, err := m.CheckMigrateVersion("DOWN", idVersion)
		if err != nil {
			log.Fatalf("failed check migrate version: %v", err)
		}
		if !approve {
			fmt.Printf("not migrate\n")
			return nil
		}
		// Получить данные миграции
		dataMigrate, err := readFile(file)
		if err != nil {
			log.Fatal(err)
		}
		// Выполнить миграцию
		if err := executeMigration(m.Cfg.Ctx, m.Cfg.Db, file, dataMigrate.Down); err != nil {
			log.Fatalf("migration failed for %s: %v", file, err)
		}
		// зафиксировать версию
		err = m.CommitMigrateVersion(false, idVersion)
		if err != nil {
			log.Fatalf("failed commit migrate version: %v", err)
		}
	}

	return nil
}
