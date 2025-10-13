package mng

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func (m *Management) GetStatusMigrations(idVersion string) (string, error) {
	rows, err := m.Cfg.DB.Conn.Query(context.Background(), GetMigrateDataReq, idVersion)
	if err != nil {
		return "", fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var versionID string
	var isApplied bool
	var tstamp time.Time
	for rows.Next() {
		if err := rows.Scan(&versionID, &isApplied, &tstamp); err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}
	}

	outputData := `
	version_id: %s
	is_applied: %v
	tstamp: %v
	=========================
	`
	if versionID == "" {
		return fmt.Sprintf("migration is not registered: %s", idVersion), nil
	}
	res := fmt.Sprintf(outputData, versionID, isApplied, tstamp)
	return res, nil
}

func (m *Management) StatusMigrations() ([]string, error) {
	resMigratesStatus := []string{}
	for _, file := range m.Cfg.MigrationFiles {
		// проверить версию миграции
		filename := filepath.Base(file)
		idVersion := strings.Split(filename, "_")[0]
		res, err := m.GetStatusMigrations(idVersion)
		if err != nil {
			return nil, fmt.Errorf("failed check migrate version: %w", err)
		}
		resMigratesStatus = append(resMigratesStatus, res)
	}
	return resMigratesStatus, nil
}
