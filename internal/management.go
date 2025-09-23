package mng

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v2"
)

func executeMigration(ctx context.Context, db *pgx.Conn, fileName, dataMigrate string) error {
	// Выполнение SQL-запроса
	_, err := db.Exec(ctx, dataMigrate)
	if err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", dataMigrate, err)
	}

	fmt.Printf("migrate execute: %s\n", filepath.Base(fileName))

	return nil
}

func readFile(filePath string) (*Migration, error) {
	// Чтение файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Создаем переменную для хранения конфигурации
	var config Migration

	// Распарсиваем YAML в структуру
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
