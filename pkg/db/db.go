package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

func Conn(ctx context.Context, dbConnPath string) (*DB, error) {
	// Подключение к базе данных
	db, err := pgx.Connect(ctx, dbConnPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{Conn: db}, nil
}
