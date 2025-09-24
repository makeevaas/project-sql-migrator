package mng

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (m *Management) CreateTableForMigrateVersion() error {
	createTableReq := `
 CREATE TABLE IF NOT EXISTS DB_version (
  id BIGSERIAL PRIMARY KEY,
  version_id VARCHAR NOT NULL,
  is_applied BOOL NOT NULL,
  tstamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
 );`
	_, err := m.Cfg.DB.Conn.Exec(context.Background(), createTableReq)
	if err != nil {
		return fmt.Errorf("failed to create table DB_version: %w", err)
	}
	return nil
}

func (m *Management) CommitMigrateVersion(applied bool, idVersion string) error {
	insertDataReq := `INSERT INTO DB_version (version_id, is_applied) VALUES ($1, $2)`
	_, err := m.Cfg.DB.Conn.Exec(context.Background(), insertDataReq, idVersion, applied)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	return nil
}

func (m *Management) CheckMigrateVersion(task, idVersion string) (bool, error) {
	err := m.CreateTableForMigrateVersion()
	if err != nil {
		return false, fmt.Errorf("failed to create table db_version: %w", err)
	}
	var approve bool
	getMigrateVersion := `SELECT version_id,is_applied,tstamp 
	from DB_version 
	where version_id=$1 
	ORDER BY tstamp 
	DESC LIMIT 1;`

	rows, err := m.Cfg.DB.Conn.Query(context.Background(), getMigrateVersion, idVersion)
	if err != nil {
		return false, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var versionID string
	var isApplied bool
	var tstamp time.Time
	for rows.Next() {
		if err := rows.Scan(&versionID, &isApplied, &tstamp); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("failed to scan row: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("error during rows iteration: %w", err)
	}

	switch task {
	case "UP":
		if errors.Is(err, sql.ErrNoRows) || !isApplied {
			// накатить можно
			approve = true
		}
	case "DOWN":
		if isApplied {
			// откатить можно
			approve = true
		}
	}
	return approve, nil
}
