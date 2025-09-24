package mng

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/makeevaas/project/sql-migrator/pkg/db"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func executeMigration(ctx context.Context, db *db.DB, fileName, dataMigrate string) error {
	// Выполнение SQL-запроса
	_, err := db.Conn.Exec(ctx, dataMigrate)
	if err != nil {
		log.Fatalf("failed to execute migration %s: %v", dataMigrate, err)
		return err
	}

	log.Info("migrate execute: ", filepath.Base(strings.Split(fileName, ".")[0]))

	return nil
}

func readFile(filePath string) (*Migration, error) {
	// Чтение файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error read file: %w", err)
	}

	// Создаем переменную для хранения конфигурации
	var config Migration

	// Распарсиваем YAML в структуру
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal yaml file: %w", err)
	}
	return &config, nil
}
