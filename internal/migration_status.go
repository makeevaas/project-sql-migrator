package mng

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (m *Management) GetStatusMigrations(idVersion string) (string, error) {
	rows, err := m.Cfg.Db.Query(context.Background(), "SELECT version_id,is_applied,tstamp from db_version where version_id=$1 ORDER BY tstamp DESC LIMIT 1;", idVersion)
	if err != nil {
		log.Fatalf("failed to query data: %v", err)
	}
	defer rows.Close()

	var versionId string
	var isApplied bool
	var tstamp time.Time
	for rows.Next() {
		if err := rows.Scan(&versionId, &isApplied, &tstamp); err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Fatalf("failed to scan row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("error during rows iteration: %v", err)
	}
	res := fmt.Sprintf("\nversion_id: %s\nis_applied: %v\ntstamp: %v\n=========================\n", versionId, isApplied, tstamp)
	return res, nil
}

func (m *Management) StatusMigrations() ([]string, error) {
	var resMigratesStatus []string
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		idVersion := strings.Split(strings.Split(file, "/")[1], "_")[0]
		res, err := m.GetStatusMigrations(idVersion)
		if err != nil {
			log.Fatalf("failed check migrate version: %v", err)
		}
		resMigratesStatus = append(resMigratesStatus, res)
	}
	return resMigratesStatus, nil
}
