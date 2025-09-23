package mng

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v2"
)

func executeMigration(ctx context.Context, db *pgx.Conn, fileName, dataMigrate string) error {
	// Выполнение SQL-запроса
	_, err := db.Exec(ctx, dataMigrate)
	if err != nil {
		log.Fatalf("failed to execute migration %s: %v", dataMigrate, err)
	}

	log.Info("migrate execute: ", filepath.Base(strings.Split(fileName, ".")[0]))

	return nil
}

func readFile(filePath string) (*Migration, error) {
	// Чтение файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error read file: %v", err)
	}

	// Создаем переменную для хранения конфигурации
	var config Migration

	// Распарсиваем YAML в структуру
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshal yaml file: %v", err)
	}
	return &config, nil
}
